package staking

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"math/big"
)

const (
	/**
	######   ######   ######   ######
	#	  THE CANDIDATE  STATUS     #
	######   ######   ######   ######
	*/
	Invalided  = 1 << iota // 0001: The current candidate withdraws from the staking qualification (Active OR Passive)
	LowRatio               // 0010: The candidate was low package ratio
	NotEnough              // 0100: The current candidate's von does not meet the minimum staking threshold
	DoubleSign             // 1000: The Double package or Double sign
	Valided    = 0         // 0000: The current candidate is in force
	NotExist   = 1 << 31   // 1000,xxxx,... : The candidate is not exist
)

const SWeightItem = 4

func Is_Valid(status uint32) bool {
	return status&Valided == Valided
}

func Is_Invalid(status uint32) bool {
	return status&Invalided == Invalided
}

func Is_PureInvalid(status uint32) bool {
	return status&Invalided == status|Invalided
}

func Is_LowRatio(status uint32) bool {
	return status&LowRatio == LowRatio
}

func Is_PureLowRatio(status uint32) bool {
	return status&LowRatio == status|LowRatio
}

func Is_NotEnough(status uint32) bool {
	return status&NotEnough == NotEnough
}

func Is_PureNotEnough(status uint32) bool {
	return status&NotEnough == status|NotEnough
}

func Is_Invalid_LowRatio(status uint32) bool {
	return status&(Invalided|LowRatio) == (Invalided | LowRatio)
}

func Is_Invalid_NotEnough(status uint32) bool {
	return status&(Invalided|NotEnough) == (Invalided | NotEnough)
}

func Is_Invalid_LowRatio_NotEnough(status uint32) bool {
	return status&(Invalided|LowRatio|NotEnough) == (Invalided | LowRatio | NotEnough)
}

func Is_LowRatio_NotEnough(status uint32) bool {
	return status&(LowRatio|NotEnough) == (LowRatio | NotEnough)
}

func Is_DoubleSign(status uint32) bool {
	return status&DoubleSign == DoubleSign
}

func Is_DoubleSign_Invalid(status uint32) bool {
	return status&(DoubleSign|Invalided) == (DoubleSign | Invalided)
}

// The Candidate info
type Candidate struct {
	NodeId discover.NodeID
	// The account used to initiate the staking
	StakingAddress common.Address
	// The account receive the block rewards and the staking rewards
	BenifitAddress common.Address
	// The tx index at the time of staking
	StakingTxIndex uint32
	// The version of the node process
	ProcessVersion uint32
	// The candidate status
	// Reference `THE CANDIDATE  STATUS`
	Status uint32
	// The epoch number at staking or edit
	StakingEpoch uint32
	// Block height at the time of staking
	StakingBlockNum uint64
	// All vons of staking and delegated
	Shares *big.Int
	// The staking von  is circulating for effective epoch (in effect)
	Released *big.Int
	// The staking von  is circulating for hesitant epoch (in hesitation)
	ReleasedHes *big.Int
	// The staking von  is RestrictingPlan for effective epoch (in effect)
	RestrictingPlan *big.Int
	// The staking von  is RestrictingPlan for hesitant epoch (in hesitation)
	RestrictingPlanHes *big.Int

	// Node desc
	Description
}

//// EncodeRLP implements rlp.Encoder
//func (c *Candidate) EncodeRLP(w io.Writer) error {
//	return rlp.Encode(w, &c)
//}
//
//
//// DecodeRLP implements rlp.Decoder
//func (c *Candidate) DecodeRLP(s *rlp.Stream) error {
//	if err := s.Decode(&c); err != nil {
//		return err
//	}
//	return nil
//}



type Description struct {
	// External Id for the third party to pull the node description (with length limit)
	ExternalId string
	// The Candidate Node's Name  (with a length limit)
	NodeName string
	// The third-party home page of the node (with a length limit)
	Website string
	// Description of the node (with a length limit)
	Details string
}

type CandidateQueue []*Candidate

// the Validator info
// They are Simplified Candidate
// They are consensus nodes and Epoch nodes snapshot
type Validator struct {
	NodeAddress common.Address
	NodeId      discover.NodeID
	// The weight
	// NOTE:
	// converted from the weight of Candidate is: (Int.Max - candidate.shares) + blocknum + txindex
	StakingWeight [SWeightItem]string
	// Validator's term in the consensus round
	ValidatorTerm uint32
}

type ValidatorQueue []*Validator

//type SlashMark map[discover.NodeID]struct{}
type SlashCandidate map[common.Address]*Candidate

func (arr ValidatorQueue) ValidatorSort(slashs SlashCandidate) {
	if len(arr) <= 1 {
		return
	}
	arr.quickSort(slashs, 0, len(arr)-1)
}
func (arr ValidatorQueue) quickSort(slashs SlashCandidate, left, right int) {
	if left < right {
		pivot := arr.partition(slashs, left, right)
		arr.quickSort(slashs, left, pivot-1)
		arr.quickSort(slashs, pivot+1, right)
	}
}
func (arr ValidatorQueue) partition(slashs SlashCandidate, left, right int) int {
	for left < right {
		for left < right && compare(slashs, arr[left], arr[right]) >= 0 {
			right--
		}
		if left < right {
			arr[left], arr[right] = arr[right], arr[left]
			left++
		}
		for left < right && compare(slashs, arr[left], arr[right]) >= 0 {
			left++
		}
		if left < right {
			arr[left], arr[right] = arr[right], arr[left]
			right--
		}
	}
	return left
}

func compare(slashs SlashCandidate, c, can *Validator) int {
	// TODO
	return -1
}

// some consensus round validators or current epoch validators
type Validator_array struct {
	// the round start blockNumber or epoch start blockNumber
	Start uint64
	// the round end blockNumber or epoch blockNumber
	End uint64
	// the round validators or epoch validators
	Arr ValidatorQueue
}

type ValidatorEx struct {
	NodeId discover.NodeID
	// The account used to initiate the staking
	StakingAddress common.Address
	// The account receive the block rewards and the staking rewards
	BenifitAddress common.Address
	// The tx index at the time of staking
	StakingTxIndex uint32
	// The version of the node process
	ProcessVersion uint32
	// Block height at the time of staking
	StakingBlockNum uint64
	// All vons of staking and delegated
	Shares *big.Int
	// Node desc
	Description
	// this is the term of validator in consensus round
	// [0, N]
	ValidatorTerm uint32
}

type ValidatorExQueue = []*ValidatorEx

// the Delegate information
type Delegation struct {
	// The epoch number at delegate or edit
	DelegateEpoch uint32
	// The delegate von  is circulating for effective epoch (in effect)
	Released *big.Int
	// The delegate von  is circulating for hesitant epoch (in hesitation)
	ReleasedHes *big.Int
	// The delegate von  is RestrictingPlan for effective epoch (in effect)
	RestrictingPlan *big.Int
	// The delegate von  is RestrictingPlan for hesitant epoch (in hesitation)
	RestrictingPlanHes *big.Int
	// Total amount in all cancellation plans
	Reduction *big.Int
}

type DelegationEx struct {
	Addr            common.Address
	NodeId          discover.NodeID
	StakingBlockNum uint64
	Delegation
}

type DelegateRelated struct {
	Addr            common.Address
	NodeId          discover.NodeID
	StakingBlockNum uint64
}

type DelRelatedQueue = []*DelegateRelated

/*type UnStakeItem struct {
	// this is the nodeAddress
	KeySuffix  	[]byte
	Amount 		*big.Int
}*/

type UnDelegateItem struct {
	// this is the `delegateAddress` + `nodeAddress` + `stakeBlockNumber`
	KeySuffix []byte
	Amount    *big.Int
}
