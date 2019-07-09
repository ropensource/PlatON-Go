package cbft

import (
	"bytes"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p"
	"math"
	"reflect"
	"sync"
)

const (
	DEFAULT_FANOUT_VALUE = 5
)

type router struct {
	cbft       *Cbft
	msgHandler handler                               // Used to send or receive logical processing of messages.
	filter     func(*peer, uint64, interface{}) bool // Used for filtering node
	routerLock sync.RWMutex
}

func NewRouter(cbft *Cbft, hd handler) *router {
	return &router{
		cbft:       cbft,
		msgHandler: hd,
		filter: func(p *peer, msgType uint64, condition interface{}) bool {
			return p.knownMessageHash.Contains(condition)
		},
	}
}

// pass the message by gossip protocol.
func (r *router) gossip(m *MsgPackage) {
	// todo: need to check
	msgType := MessageType(m.msg.Message)
	msgHash := m.msg.MsgHash()

	switch msgType {
	case ConfirmedPrepareBlockMsg, PrepareBlockHashMsg:
		// check the message is repeated
		if r.repeatedCheck(m.peerID, msgHash) {
			log.Debug("The message is repeated, not to forward again", "msgType", reflect.TypeOf(m.msg.Message), "msgHash", msgHash.TerminalString())
			return
		}
	}
	peers, err := r.selectNodesByMsgType(msgType, msgHash)
	if err != nil {
		log.Error("select nodes fail in the gossip method. gossip fail", "msgType", msgType)
		return
	}
	switch m.mode {
	case MixMode:
	case FullMode:
	case PartMode:
		transfer := peers[:int(math.Sqrt(float64(len(peers))))]
		peers = transfer
	}
	pids := formatPeers(peers)
	log.Debug("Gossip message", "msgHash", msgHash.TerminalString(), "msgType", reflect.TypeOf(m.msg.Message), "targetPeer", pids)

	r.cbft.tracing.RecordSend(r.cbft.config.NodeID.TerminalString(), msgHash.TerminalString(), fmt.Sprintf("%T", m.msg.Message), pids)
	for _, peer := range peers {
		if err := p2p.Send(peer.rw, msgType, m.msg.Message); err != nil {
			log.Error("Send message failed", "peer", peer.id, "err", err)
		} else {
			peer.MarkMessageHash(msgHash)
		}
	}
}

// formatPeers is used to print the information about peer
func formatPeers(peers []*peer) string {
	var bf bytes.Buffer
	for idx, peer := range peers {
		bf.WriteString(peer.id)
		if idx < len(peers)-1 {
			bf.WriteString(",")
		}
	}
	return bf.String()
}

func (r *router) selectNodesByMsgType(msgType uint64, condition interface{}) ([]*peer, error) {
	r.routerLock.RLock()
	defer r.routerLock.RUnlock()
	switch msgType {
	case PrepareBlockMsg, PrepareVoteMsg, ConfirmedPrepareBlockMsg,
		PrepareBlockHashMsg:
		return r.kMixingRandomNodes(msgType, condition)
	case ViewChangeMsg, GetPrepareBlockMsg, GetHighestPrepareBlockMsg, ViewChangeVoteMsg:
		return r.kConsensusRandomNodes(msgType, condition)
	}
	return nil, fmt.Errorf("no found nodes")
}

// Return consensus nodes by random.
func (r *router) kConsensusRandomNodes(msgType uint64, condition interface{}) ([]*peer, error) {
	cNodes, err := r.cbft.ConsensusNodes()
	if err != nil {
		return nil, err
	}
	existsPeers := r.msgHandler.PeerSet().Peers()
	log.Debug("kConsensusRandomNodes select node", "msgHash", condition, "cNodesLen", len(cNodes), "peerSetLen", len(existsPeers))
	consensusPeers := make([]*peer, 0)
	for _, peer := range existsPeers {
		if msgType != GetHighestPrepareBlockMsg && peer.knownMessageHash.Contains(condition) {
			continue
		}
		for _, node := range cNodes {
			if peer.id == fmt.Sprintf("%x", node.Bytes()[:8]) {
				consensusPeers = append(consensusPeers, peer)
				break
			}
		}
	}
	//kConsensusNodes := kRandomNodes(len(consensusPeers), consensusPeers, msgType, condition, r.filter)
	return consensusPeers, nil

}

// Return the nodes of consensus and non-consensus.
func (r *router) kMixingRandomNodes(msgType uint64, condition interface{}) ([]*peer, error) {
	// all consensus nodes + a number of k non-consensus nodes
	cNodes, err := r.cbft.ConsensusNodes()
	if err != nil {
		return nil, err
	}
	existsPeers := r.msgHandler.PeerSet().Peers()
	consensusPeers := make([]*peer, 0)
	nonconsensusPeers := make([]*peer, 0)
	for _, peer := range existsPeers { //
		isConsensus := false
		for _, node := range cNodes {
			if peer.id == fmt.Sprintf("%x", node.Bytes()[:8]) {
				isConsensus = true
				break
			}
		}
		if isConsensus {
			consensusPeers = append(consensusPeers, peer)
		} else {
			nonconsensusPeers = append(nonconsensusPeers, peer)
		}
	}
	// todo: need to test
	kNonconsensusNodes := kRandomNodes(len(nonconsensusPeers)/2, nonconsensusPeers, msgType, condition, r.filter)
	consensusPeers = append(consensusPeers, kNonconsensusNodes...)
	return consensusPeers, nil
}

// Return the completely random nodes.
func (r *router) kPureRandomNodes(msgType uint64, condition interface{}) ([]*peer, error) {
	existsPeers := r.msgHandler.PeerSet().Peers()
	kConsensusNodes := kRandomNodes(DEFAULT_FANOUT_VALUE, existsPeers, msgType, condition, r.filter)
	return kConsensusNodes, nil
}

// kRandomNodes is used to select up to k random nodes, excluding any nodes where
// the filter function returns true. It is possible that less than k nodes are returned.
func kRandomNodes(k int, peers []*peer, msgType uint64, condition interface{}, filterFn func(*peer, uint64, interface{}) bool) []*peer {
	n := len(peers)
	kNodes := make([]*peer, 0, k)
OUTER:
	// Probe up to 3*n times, with large n this is not necessary
	// since k << n, but with small n we want search to be
	// exhaustive.
	for i := 0; i < 3*n && len(kNodes) < k; i++ {
		// Get random node
		idx := randomOffset(n)
		node := peers[idx]

		// Give the filter a shot at it.
		if filterFn != nil && filterFn(node, msgType, condition) {
			continue OUTER
		}

		// Check if we have this node already
		for j := 0; j < len(kNodes); j++ {
			if node == kNodes[j] {
				continue OUTER
			}
		}
		// Append the node
		kNodes = append(kNodes, node)
	}
	return kNodes
}

func (r *router) repeatedCheck(peerId string, msgHash common.Hash) bool {
	peers := r.cbft.handler.PeerSet().Peers()
	for _, peer := range peers {
		if peer.id == peerId {
			continue
		}
		if peer.knownMessageHash.Contains(msgHash) {
			return true
		}
	}
	return false
}
