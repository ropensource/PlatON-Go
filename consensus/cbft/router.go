package cbft

import (
	"bytes"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/event"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p"
	"reflect"
	"sync"
	"time"
)

const (
	DEFAULT_FANOUT_VALUE = 5
)

type peerFilter func(peers []*peer) []*peer

type router struct {
	msgHandler *handler 							// Used to send or receive logical processing of messages.
	filter func(*peer, uint64, interface{}) bool	// Used for filtering node
	pFilter peerFilter
	routerLock sync.RWMutex

	preBlockCh         		chan NewPrepareBlockEvent
	preBlockSub        		event.Subscription

	preBlockHashCh     		chan NewPrepareBlockHashEvent
	preBlockHashSub   	 	event.Subscription

	confPreBlockCh     		chan NewConfirmedPrepareBlockEvent
	confPreBlockSub    		event.Subscription

	preVoteCh         		chan NewPrepareVoteEvent
	preVoteSub        		event.Subscription

	//viewChangeCh        	chan NewViewChangeEvent
	//viewChangeSub       	event.Subscription

	viewChangeVoteCh       	chan NewViewChangeVoteEvent
	viewChangeVoteSub       event.Subscription
}

func NewRouter(hd *handler) *router {
	r := &router{
		msgHandler: hd,
		filter: func(p *peer, msgType uint64, condition interface{}) bool {
			// todo: Need to distinguish according to message type
			return p.knownMessageHash.Contains(condition)
		},
		pFilter: func(peers []*peer) []*peer {
			return peers
		},
	}

	r.preBlockCh = make(chan NewPrepareBlockEvent, 1024)
	r.preBlockSub = r.msgHandler.cbft.SubscribeNewPrepareBlockEvent(r.preBlockCh)
	r.preVoteCh = make(chan NewPrepareVoteEvent, 1024)
	r.preVoteSub = r.msgHandler.cbft.SubscribeNewPrepareVoteEvent(r.preVoteCh)
	r.confPreBlockCh = make(chan NewConfirmedPrepareBlockEvent, 1024)
	r.confPreBlockSub = r.msgHandler.cbft.SubscribeNewConfirmedPrepareBlockEvent(r.confPreBlockCh)
	r.preBlockHashCh = make(chan NewPrepareBlockHashEvent, 2048)
	r.preBlockHashSub = r.msgHandler.cbft.SubscribeNewPrepareBlockHashEvent(r.preBlockHashCh)
	//r.viewChangeCh = make(chan NewViewChangeEvent,1024)
	//r.viewChangeSub = r.msgHandler.cbft.SubscribeNewViewChangeEvent(r.viewChangeCh)
	r.viewChangeVoteCh = make(chan NewViewChangeVoteEvent, 2048)
	r.viewChangeVoteSub = r.msgHandler.cbft.SubscribeNewViewChangeVoteEvent(r.viewChangeVoteCh)

	go r.preBlockBroadcastLoop()
	go r.prepareVoteBroadcastLoop()
	go r.confirmedPrepareBlockBroadcastLoop()
	go r.prepareBlockHashBroadcastLoop()
	//go r.viewChangeBroadcastLoop()
	go r.viewChangeVoteBroadcastLoop()

	return r
}

func (r *router) preBlockBroadcastLoop() {
	for {
		select {
		case event := <-r.preBlockCh:
			r.BroadcastPreBlock(event.PrepareBlock, r.pFilter)

		case <-r.preBlockSub.Err():
			return
		}
	}
}

func (r *router) BroadcastPreBlock(msg *prepareBlock, filter peerFilter) {
	msgHash := msg.MsgHash()
	peersSet := r.msgHandler.peers.PeersWithoutPrepareBlock(msgHash)
	peers := filter(peersSet)
	//todo: need to split the type of nodes
	for _, peer := range peers {
		peer.AsyncSendPrepareBlock(msg)
	}
	log.Trace("Propagated prepare block", "hash", msg.Block.Hash(), "number", msg.Block.Number(), "msgHash", msg.MsgHash(), "duration", common.PrettyDuration(time.Since(msg.Block.ReceivedAt)))
}

// the loop to broadcast prepare vote
func (r *router) prepareVoteBroadcastLoop() {
	for {
		select {
		case event := <-r.preVoteCh:
			r.BroadcastPrepareVote(event.PrepareVote, r.pFilter)

		case <-r.preVoteSub.Err():
			return
		}
	}
}

func (r *router) BroadcastPrepareVote(msg *prepareVote, filter peerFilter) {
	msgHash := msg.MsgHash()
	peersSet := r.msgHandler.peers.PeersWithoutPrepareVote(msgHash)
	peers := filter(peersSet)
	for _, peer := range peers {
		peer.AsyncSendPrepareVote(msg)
	}
	log.Trace("Propagated prepare vote", "hash", msg.Hash, "number", msg.Number, "msgHash", msg.MsgHash())
}

func (r *router) confirmedPrepareBlockBroadcastLoop() {
	for {
		select {
		case event := <-r.confPreBlockCh:
			r.BroadcastConfirmedPrepareBlock(event.ConfirmedPrepareBlock, r.pFilter)

		case <-r.confPreBlockSub.Err():
			return
		}
	}
}

func (r *router) BroadcastConfirmedPrepareBlock(msg *confirmedPrepareBlock, filter peerFilter) {
	msgHash := msg.MsgHash()
	peersSet := r.msgHandler.peers.PeersWithoutConfirmedPrepareBlock(msgHash)
	peers := filter(peersSet)
	for _, peer := range peers {
		peer.AsyncSendConfirmedPrepareBlock(msg)
	}
	log.Trace("Propagated confirmed prepare block", "hash", msg.Hash, "number", msg.Number, "msgHash", msg.MsgHash())
}

func (r *router) prepareBlockHashBroadcastLoop() {
	for {
		select {
		case event := <-r.preBlockHashCh:
			r.BroadcastPrepareBlockHash(event.PrepareBlockHash, r.pFilter)

		case <-r.preBlockHashSub.Err():
			return
		}
	}
}

func (r *router) BroadcastPrepareBlockHash(msg *prepareBlockHash, filter peerFilter) {
	msgHash := msg.MsgHash()
	peersSet := r.msgHandler.peers.PeersWithoutPrepareBlockHash(msgHash)
	peers := filter(peersSet)
	for _, peer := range peers {
		peer.AsyncSendPrepareBlockHash(msg)
	}
	log.Trace("Propagated prepare block hash", "hash", msg.Hash, "number", msg.Number, "msgHash", msg.MsgHash())
}

/*func (r *router) viewChangeBroadcastLoop() {
	for {
		select {
		case event := <-r.viewChangeCh:
			r.BroadcastViewChange(event.ViewChange, r.pFilter)

		case <-r.viewChangeSub.Err():
			return
		}
	}
}*/

/*func (r *router) BroadcastViewChange(msg *viewChange, filter peerFilter) {
	msgHash := msg.MsgHash()
	peersSet := r.msgHandler.peers.PeersWithoutViewChange(msgHash)
	peers := filter(peersSet)
	for _, peer := range peers {
		peer.AsyncSendViewChnage(msg)
	}
	log.Trace("Propagated view change", "hash", msg.BaseBlockHash, "number", msg.BaseBlockNum, "msgHash", msg.MsgHash())
}*/

func (r *router) viewChangeVoteBroadcastLoop() {
	for {
		select {
		case event := <-r.viewChangeVoteCh:
			r.BroadcastViewChangeVote(event.ViewChangeVote, r.pFilter)

		case <-r.viewChangeVoteSub.Err():
			return
		}
	}
}

func (r *router) BroadcastViewChangeVote(msg *viewChangeVote, filter peerFilter) {
	msgHash := msg.MsgHash()
	peersSet := r.msgHandler.peers.PeersWithoutViewChangeVote(msgHash)
	peers := filter(peersSet)
	for _, peer := range peers {
		peer.AsyncSendViewChnageVote(msg)
	}
	log.Trace("Propagated view change vote", "hash", msg.BlockHash, "number", msg.BlockNum, "msgHash", msg.MsgHash())
}

// pass the message by gossip protocol.
func (r *router) gossip(m *MsgPackage) {
	msgType := MessageType(m.msg)
	msgHash := m.msg.MsgHash()
	peers, err := r.selectNodesByMsgType(msgType, msgHash)
	if err != nil {
		log.Error("select nodes fail in the gossip method. gossip fail", "msgType", msgType)
		return
	}
	log.Debug("[Method:gossip] gossip message", "msgHash", msgHash.TerminalString(), "msgType", reflect.TypeOf(m.msg), "targetPeer", formatPeers(peers))
	filter := func(peers []*peer) []*peer {
		return peers
	}
	//
	switch msgType {
	case PrepareBlockMsg:
		r.BroadcastPreBlock(m.msg.(*prepareBlock), filter)
	case PrepareVoteMsg:
		r.BroadcastPrepareVote(m.msg.(*prepareVote), filter)
	case ViewChangeVoteMsg:
		r.BroadcastViewChangeVote(m.msg.(*viewChangeVote), filter)
	case ConfirmedPrepareBlockMsg:
		r.BroadcastConfirmedPrepareBlock(m.msg.(*confirmedPrepareBlock), filter)
	case PrepareBlockHashMsg:
		r.BroadcastPrepareBlockHash(m.msg.(*prepareBlockHash), filter)
	case ViewChangeMsg:
		fallthrough
	case GetPrepareVoteMsg:
		fallthrough
	case PrepareVotesMsg:
		fallthrough
	case GetPrepareBlockMsg:
		fallthrough
	case GetHighestPrepareBlockMsg:
		fallthrough
	case HighestPrepareBlockMsg:
		fallthrough
	default:
		for _, peer := range peers {
			if err := p2p.Send(peer.rw, msgType, m.msg); err != nil {
				log.Error("Send message failed", "peer", peer.id, "err", err)
			} else {
				peer.MarkMessageHash(msgHash)
			}
		}
	}
}

// formatPeers is used to print the information about peer
func formatPeers(peers []*peer) string {
	var bf bytes.Buffer
	for idx, peer := range peers {
		bf.WriteString(peer.id)
		if idx < len(peers) - 1 {
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
	cNodes, err := r.msgHandler.cbft.ConsensusNodes()
	if err != nil {
		return nil, err
	}
	existsPeers := r.msgHandler.peers.Peers()
	consensusPeers := make([]*peer, 0)
	for _, peer := range existsPeers {
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
	cNodes, err := r.msgHandler.cbft.ConsensusNodes()
	if err != nil {
		return nil, err
	}
	existsPeers := r.msgHandler.peers.Peers()
	consensusPeers := make([]*peer, 0)
	nonconsensusPeers := make([]*peer, 0)
	for _, peer := range existsPeers {	//
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
	kNonconsensusNodes := kRandomNodes(len(nonconsensusPeers) / 2, nonconsensusPeers, msgType, condition, r.filter)
	consensusPeers = append(consensusPeers, kNonconsensusNodes...)
	return consensusPeers, nil
}

// Return the completely random nodes.
func (r *router) kPureRandomNodes(msgType uint64, condition interface{}) ([]*peer, error) {
	existsPeers := r.msgHandler.peers.Peers()
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
	for i := 0; i < 3 * n && len(kNodes) < k; i++ {
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
