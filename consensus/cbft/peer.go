package cbft

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"math/big"
	"sync"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/p2p"
	"github.com/deckarep/golang-set"
)

var (
	errClosed            = errors.New("peer set is closed")
	errAlreadyRegistered = errors.New("peer is already registered")
	errNotRegistered     = errors.New("peer is not registered")
)

const (
	maxKnownMessageHash = 60000

	handshakeTimeout = 5 * time.Second

	maxKnownPreBlocks       = 10000 // Maximum prepare blocks to keep in the known list (prevent DOS)
	maxKnownPreVotes        = 1024  // Maximum prepare votes to keep in the known list (prevent DOS)
	maxKnownConfPreBlocks   = 1024
	maxKnownPreBlockHashes  = 2048
	maxKnownViewChanges     = 2048
	maxKnownViewChangeVotes = 2048

	maxQueuedPreBlocks       = 4
	maxQueuedPreVotes        = 4
	maxQueuedConfPreBlocks   = 4
	maxQueuedPreBlockHashes  = 4
	maxQueuedViewChanges     = 4
	maxQueuedViewChangeVotes = 10
)

type peer struct {
	id   string
	term chan struct{} // Termination channel to stop the broadcaster
	*p2p.Peer
	rw p2p.MsgReadWriter

	knownMessageHash mapset.Set

	// Messages that have been received
	knownPreBlock       mapset.Set
	knownPreVote        mapset.Set
	knownConfPreBlock   mapset.Set
	knownPreBlockHash   mapset.Set
	knownViewChange     mapset.Set
	knownViewChangeVote mapset.Set

	// pending message queue
	queuedPreBlocks       chan *prepareBlock
	queuedPreVotes        chan *prepareVote
	queuedConfPreBlocks   chan *confirmedPrepareBlock
	queuedPreBlockHashes  chan *prepareBlockHash
	queuedViewChanges     chan *viewChange
	queuedViewChangeVotes chan *viewChangeVote
}

func newPeer(p *p2p.Peer, rw p2p.MsgReadWriter) *peer {
	return &peer{
		Peer:                p,
		rw:                  rw,
		id:                  fmt.Sprintf("%x", p.ID().Bytes()[:8]),
		term:                make(chan struct{}),
		knownMessageHash:    mapset.NewSet(),
		knownPreBlock:       mapset.NewSet(),
		knownPreVote:        mapset.NewSet(),
		knownConfPreBlock:   mapset.NewSet(),
		knownPreBlockHash:   mapset.NewSet(),
		knownViewChange:     mapset.NewSet(),
		knownViewChangeVote: mapset.NewSet(),

		queuedPreBlocks:       make(chan *prepareBlock, maxQueuedPreBlocks),
		queuedPreVotes:        make(chan *prepareVote, maxQueuedPreVotes),
		queuedConfPreBlocks:   make(chan *confirmedPrepareBlock, maxQueuedConfPreBlocks),
		queuedPreBlockHashes:  make(chan *prepareBlockHash, maxQueuedPreBlockHashes),
		queuedViewChanges:     make(chan *viewChange, maxQueuedViewChanges),
		queuedViewChangeVotes: make(chan *viewChangeVote, maxQueuedViewChangeVotes),
	}
}

func (p *peer) close() {
	close(p.term)
}

func (p *peer) MarkMessageHash(hash common.Hash) {
	for p.knownMessageHash.Cardinality() >= maxKnownMessageHash {
		p.knownMessageHash.Pop()
	}
	p.knownMessageHash.Add(hash)
}

func (p *peer) MarkPrepareBlock(msgHash common.Hash) {
	for p.knownPreBlock.Cardinality() >= maxKnownPreBlocks {
		p.knownPreBlock.Pop()
	}
	p.knownPreBlock.Add(msgHash)
}

func (p *peer) MarkPrepareVote(msgHash common.Hash) {
	for p.knownPreVote.Cardinality() >= maxKnownPreVotes {
		p.knownPreVote.Pop()
	}
	p.knownPreVote.Add(msgHash)
}

func (p *peer) MarkConfirmedPrepareBlock(msgHash common.Hash) {
	for p.knownConfPreBlock.Cardinality() >= maxKnownConfPreBlocks {
		p.knownConfPreBlock.Pop()
	}
	p.knownConfPreBlock.Add(msgHash)
}

func (p *peer) MarkPrepareBlockHash(msgHash common.Hash) {
	for p.knownPreBlockHash.Cardinality() >= maxKnownPreBlockHashes {
		p.knownPreBlockHash.Pop()
	}
	p.knownPreBlockHash.Add(msgHash)
}

func (p *peer) MarkViewChange(msgHash common.Hash) {
	for p.knownViewChange.Cardinality() >= maxKnownViewChanges {
		p.knownViewChange.Pop()
	}
	p.knownViewChange.Add(msgHash)
}

func (p *peer) MarkViewChangeVote(msgHash common.Hash) {
	for p.knownViewChangeVote.Cardinality() >= maxKnownViewChangeVotes {
		p.knownViewChangeVote.Pop()
	}
	p.knownViewChangeVote.Add(msgHash)
}


// exchange node information with each other.
func (p *peer) Handshake(bn *big.Int, head common.Hash) error {
	errc := make(chan error, 2)
	var status cbftStatusData

	go func() {
		errc <- p2p.Send(p.rw, CBFTStatusMsg, &cbftStatusData{
			BN:           bn,
			CurrentBlock: head,
		})
	}()
	go func() {
		errc <- p.readStatus(&status)
		if status.BN != nil {
			p.Log().Debug("[Method:Handshake] Receive the cbftStatusData message", "blockHash", status.CurrentBlock.TerminalString(), "blockNumber", status.BN.Int64())
		}
	}()
	timeout := time.NewTicker(handshakeTimeout)
	defer timeout.Stop()
	for i := 0; i < 2; i++ {
		select {

		case err := <-errc:
			if err != nil {
				return err
			}
		case <-timeout.C:
			return p2p.DiscReadTimeout
		}
	}
	// todo: Maybe there is something to be done.
	return nil
}

func (p *peer) readStatus(status *cbftStatusData) error {
	msg, err := p.rw.ReadMsg()

	if err != nil {
		return err
	}
	if msg.Code != CBFTStatusMsg {
		return errResp(ErrNoStatusMsg, "first msg has code %x (!= %x)", msg.Code, CBFTStatusMsg)
	}
	if msg.Size > CbftProtocolMaxMsgSize {
		return errResp(ErrMsgTooLarge, "%v > %v", msg.Size, CbftProtocolMaxMsgSize)
	}
	if err := msg.Decode(&status); err != nil {

		return errResp(ErrDecode, "msg %v: %v", msg, err)
	}
	// todo: additional judgment.
	return nil
}

// main loop to send message
func (p *peer) broadcast() {
	for {
		select {
		case preBlock := <-p.queuedPreBlocks:
			if err := p.SendPrepareBlock(preBlock); err != nil {
				return
			}
			p.Log().Trace("Broadcast prepare block", "number", preBlock.Block.Number(), "hash", preBlock.Block.Hash(), "msgHash", preBlock.MsgHash().TerminalString())

		case preVote := <-p.queuedPreVotes:
			if err := p.SendPrepareVote(preVote); err != nil {
				return
			}
			p.Log().Trace("Broadcast prepare vote", "number", preVote.Number, "hash", preVote.Hash, "msgHash", preVote.MsgHash().TerminalString())

		case confPreBlock := <-p.queuedConfPreBlocks:
			if err := p.SendConfirmedPrepareBlock(confPreBlock); err != nil {
				return
			}
			p.Log().Trace("Broadcast confirm prepare block", "number", confPreBlock.Number, "hash", confPreBlock.Hash, "msgHash", confPreBlock.MsgHash().TerminalString())

		case preBlockHash := <-p.queuedPreBlockHashes:
			if err := p.SendPrepareBlockHash(preBlockHash); err != nil {
				return
			}
			p.Log().Trace("Broadcast prepare block hash", "number", preBlockHash.Number, "hash", preBlockHash.Hash, "msgHash", preBlockHash.MsgHash().TerminalString())

		case viewChange := <-p.queuedViewChanges:
			if err := p.SendViewChange(viewChange); err != nil {
				return
			}
			p.Log().Trace("Broadcast view change", "number", viewChange.BaseBlockNum, "hash", viewChange.BaseBlockHash, "msgHash", viewChange.MsgHash().TerminalString())

		case viewChangeVote := <-p.queuedViewChangeVotes:
			if err := p.SendViewChangeVote(viewChangeVote); err != nil {
				return
			}
			p.Log().Trace("Broadcast view change vote", "number", viewChangeVote.BlockNum, "hash", viewChangeVote.BlockHash, "msgHash", viewChangeVote.MsgHash().TerminalString())

		case <-p.term:
			return
		}
	}
}

func (p *peer) SendPrepareBlock(preBlock *prepareBlock) error {
	p.knownPreBlock.Add(preBlock.MsgHash())
	return p2p.Send(p.rw, PrepareBlockMsg, preBlock)
}

func (p *peer) AsyncSendPrepareBlock(preBlock *prepareBlock) {
	select {
	case p.queuedPreBlocks <- preBlock:
		p.knownPreBlock.Add(preBlock.MsgHash())
	default:
		p.Log().Debug("Dropping prepare block propagation", "number", preBlock.Block.Number(), "hash", preBlock.Block.Hash(), "msgHash", preBlock.MsgHash().TerminalString())
	}
}

func (p *peer) SendPrepareVote(preVote *prepareVote) error {
	p.knownPreVote.Add(preVote.MsgHash())
	return p2p.Send(p.rw, PrepareVoteMsg, preVote)
}

func (p *peer) AsyncSendPrepareVote(preVote *prepareVote) {
	select {
	case p.queuedPreVotes <- preVote:
		p.knownPreVote.Add(preVote.MsgHash())
	default:
		p.Log().Debug("Dropping prepare vote propagation", "number", preVote.Number, "hash", preVote.Hash, "msgHash", preVote.MsgHash())
	}
}

func (p *peer) SendConfirmedPrepareBlock(confPreBlock *confirmedPrepareBlock) error {
	p.knownConfPreBlock.Add(confPreBlock.MsgHash())
	return p2p.Send(p.rw, ConfirmedPrepareBlockMsg, confPreBlock)
}

func (p *peer) AsyncSendConfirmedPrepareBlock(confPreBlock *confirmedPrepareBlock) {
	select {
	case p.queuedConfPreBlocks <- confPreBlock:
		p.knownConfPreBlock.Add(confPreBlock.MsgHash())
	default:
		p.Log().Debug("Dropping confirmed prepare block propagation", "number", confPreBlock.Number, "hash", confPreBlock.Hash, "msgHash", confPreBlock.MsgHash().TerminalString())
	}
}

func (p *peer) SendPrepareBlockHash(preBlockHash *prepareBlockHash) error {
	p.knownPreBlockHash.Add(preBlockHash.MsgHash())
	return p2p.Send(p.rw, PrepareBlockHashMsg, preBlockHash)
}

func (p *peer) AsyncSendPrepareBlockHash(preBlockHash *prepareBlockHash) {
	select {
	case p.queuedPreBlockHashes <- preBlockHash:
		p.knownPreBlockHash.Add(preBlockHash.MsgHash())
	default:
		p.Log().Debug("Dropping prepare block hash propagation", "number", preBlockHash.Number, "hash", preBlockHash.Hash, "msgHash", preBlockHash.MsgHash().TerminalString())
	}
}

func (p *peer) SendViewChange(viewChange *viewChange) error {
	p.knownViewChange.Add(viewChange.MsgHash())
	return p2p.Send(p.rw, ViewChangeMsg, viewChange)
}

func (p *peer) AsyncSendViewChnage(viewChange *viewChange) {
	select {
	case p.queuedViewChanges <- viewChange:
		p.knownViewChange.Add(viewChange.MsgHash())
	default:
		p.Log().Debug("Dropping view change msg propagation", "number", viewChange.BaseBlockNum, "hash", viewChange.BaseBlockHash, "msgHash", viewChange.MsgHash().TerminalString())
	}
}

func (p *peer) SendViewChangeVote(viewChangeVote *viewChangeVote) error {
	p.knownViewChangeVote.Add(viewChangeVote.MsgHash())
	return p2p.Send(p.rw, ViewChangeVoteMsg, viewChangeVote)
}

func (p *peer) AsyncSendViewChnageVote(viewChangeVote *viewChangeVote) {
	select {
	case p.queuedViewChangeVotes <- viewChangeVote:
		p.knownViewChangeVote.Add(viewChangeVote.MsgHash())
	default:
		p.Log().Debug("Dropping view change vote msg propagation", "number", viewChangeVote.BlockNum, "hash", viewChangeVote.BlockHash, "msgHash", viewChangeVote.MsgHash().TerminalString())
	}
}

type peerSet struct {
	peers  map[string]*peer
	lock   sync.RWMutex
	closed bool
}

func newPeerSet() *peerSet {
	// Monitor output node list
	ps := &peerSet{
		peers: make(map[string]*peer),
	}
	go ps.printPeers()
	return ps
}

func (ps *peerSet) Register(p *peer) {
	ps.lock.Lock()
	defer ps.lock.Unlock()
	ps.peers[p.id] = p
	go p.broadcast()
}

func (ps *peerSet) Unregister(id string) error {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	p, ok := ps.peers[id]
	if !ok {
		return errNotRegistered
	}
	delete(ps.peers, id)
	p.close()

	return nil
}

func (ps *peerSet) Get(id string) (*peer, error) {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	p, ok := ps.peers[id]
	if !ok {
		return nil, errNotRegistered
	}

	return p, nil
}

func (ps *peerSet) AllConsensusPeer() []*peer {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	list := make([]*peer, 0, len(ps.peers))
	for _, p := range ps.peers {
		list = append(list, p)
	}
	return list
}

func (ps *peerSet) Close() {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	for _, p := range ps.peers {
		p.Disconnect(p2p.DiscQuitting)
	}
	ps.closed = true
}

// Return all peer.
func (ps *peerSet) Peers() []*peer {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	list := make([]*peer, 0, len(ps.peers))
	for _, p := range ps.peers {
		list = append(list, p)
	}
	return list
}

func (ps *peerSet) PeersWithoutPrepareBlock(hash common.Hash) []*peer {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	list := make([]*peer, 0, len(ps.peers))
	for _, p := range ps.peers {
		if !p.knownPreBlock.Contains(hash) {
			list = append(list, p)
		}
	}
	return list
}

func (ps *peerSet) PeersWithoutPrepareVote(hash common.Hash) []*peer {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	list := make([]*peer, 0, len(ps.peers))
	for _, p := range ps.peers {
		if !p.knownPreVote.Contains(hash) {
			list = append(list, p)
		}
	}
	return list
}

func (ps *peerSet) PeersWithoutPrepareBlockHash(hash common.Hash) []*peer {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	list := make([]*peer, 0, len(ps.peers))
	for _, p := range ps.peers {
		if !p.knownPreBlockHash.Contains(hash) {
			list = append(list, p)
		}
	}
	return list
}

func (ps *peerSet) PeersWithoutConfirmedPrepareBlock(hash common.Hash) []*peer {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	list := make([]*peer, 0, len(ps.peers))
	for _, p := range ps.peers {
		if !p.knownConfPreBlock.Contains(hash) {
			list = append(list, p)
		}
	}
	return list
}

func (ps *peerSet) PeersWithoutViewChange(hash common.Hash) []*peer {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	list := make([]*peer, 0, len(ps.peers))
	for _, p := range ps.peers {
		if !p.knownViewChange.Contains(hash) {
			list = append(list, p)
		}
	}
	return list
}

func (ps *peerSet) PeersWithoutViewChangeVote(hash common.Hash) []*peer {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	list := make([]*peer, 0, len(ps.peers))
	for _, p := range ps.peers {
		if !p.knownViewChangeVote.Contains(hash) {
			list = append(list, p)
		}
	}
	return list
}

func (ps *peerSet) printPeers() {
	// Output in 2 seconds
	outTimer := time.NewTicker(time.Second * 2)
	for {
		if ps.closed {
			break
		}
		select {
		case <-outTimer.C:
			peers := ps.Peers()
			var bf bytes.Buffer
			for idx, peer := range peers {
				bf.WriteString(peer.id)

				if idx < len(peers)-1 {
					bf.WriteString(",")
				}
			}
			pInfo := bf.String()
			log.Debug(fmt.Sprintf("[Method:printPeers] The neighbor node owned by the current peer is : {%v}, size: {%d}", pInfo, len(peers)))
		}
	}
}
