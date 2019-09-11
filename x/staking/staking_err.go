package staking

import "github.com/PlatONnetwork/PlatON-Go/common"

var (
	ErrWrongBlsPubKey            = common.NewBizError(302000, "The bls public key is wrong")
	ErrDescriptionLen            = common.NewBizError(302001, "The Description length is wrong")
	ErrWrongProgramVersionSign   = common.NewBizError(302002, "The program version sign is wrong")
	ErrProgramVersionTooLow      = common.NewBizError(302003, "The program version of the relates node's is too low")
	ErrDeclVsFialedCreateCan     = common.NewBizError(302004, "DeclareVersion is failed on create staking")
	ErrNoSameStakingAddr         = common.NewBizError(302005, "The address must be the same as initiated staking")
	ErrStakeVonTooLow            = common.NewBizError(302100, "Staking deposit too low")
	ErrCanAlreadyExist           = common.NewBizError(302101, "This candidate is already exist")
	ErrCanNoExist                = common.NewBizError(302102, "This candidate is not exist")
	ErrCanStatusInvalid          = common.NewBizError(302103, "This candidate status was invalided")
	ErrIncreaseStakeVonTooLow    = common.NewBizError(302104, "IncreaseStake von is too low")
	ErrDelegateVonTooLow         = common.NewBizError(302105, "Delegate deposit too low")
	ErrAccountNoAllowToDelegate  = common.NewBizError(302106, "The account is not allowed to be used for delegating")
	ErrCanNoAllowDelegate        = common.NewBizError(302107, "This candidate is not allow to delegate")
	ErrWithdrewDelegateVonTooLow = common.NewBizError(302108, "Withdrew delegate von is too low")
	ErrDelegateNoExist           = common.NewBizError(302109, "This is delegate is not exist")
	ErrWrongVonOptType           = common.NewBizError(302110, "The von operationType is wrong")
	ErrAccountVonNoEnough        = common.NewBizError(302111, "The von of account is not enough")
	ErrBlockNumberDisordered     = common.NewBizError(302112, "The blockNumber is disordered")
	ErrDelegateVonNoEnough       = common.NewBizError(302113, "The von of delegate is not enough")
	ErrWrongWithdrewDelVonCalc   = common.NewBizError(302114, "Withdrew delegate von calculate is wrong")
	ErrValidatorNoExist          = common.NewBizError(302115, "The validator is not exist")
	ErrWrongFuncParams           = common.NewBizError(302116, "The fn params is wrong")
	ErrWrongSlashType            = common.NewBizError(302117, "The slash type is wrong")
	ErrSlashVonTooLarge          = common.NewBizError(302118, "Slash amount is too large")
	ErrWrongSlashVonCalc         = common.NewBizError(302119, "Slash candidate von calculate is wrong")
	ErrGetVerifierList           = common.NewBizError(302200, "Getting verifierList is failed")
	ErrGetValidatorList          = common.NewBizError(302201, "Getting validatorList is failed")
	ErrGetCandidateList          = common.NewBizError(302202, "Getting candidateList is failed")
	ErrGetDelegateRelated        = common.NewBizError(302203, "Getting related of delegate is failed")
	ErrQueryCandidateInfo        = common.NewBizError(302204, "Query candidate info failed")
	ErrQueryDelegateInfo         = common.NewBizError(302205, "Query delegate info failed")
)
