package vm

import (
	"bytes"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/PlatONnetwork/PlatON-Go/x/xutil"
	"math/big"
	"testing"
)

func TestRLP_encode (t *testing.T) {

	var params [][]byte
	params = make([][]byte, 0)

	fnType, err := rlp.EncodeToBytes(uint16(1102))
	if nil != err {
		fmt.Println("fnType err", err)
	}else {
		var num uint16
		rlp.DecodeBytes(fnType, &num)
		fmt.Println("num is ", num)
	}
	params = append(params, fnType)

	buf := new(bytes.Buffer)
	err = rlp.Encode(buf, params)
	if err != nil {
		fmt.Println(err)
		t.Errorf("CandidateDeposit encode rlp data fail")
	} else {
		fmt.Println("CandidateDeposit data rlp: ", hexutil.Encode(buf.Bytes()))
	}
}


func TestRLP_2 (t *testing.T) {

	var params [][]byte
	params = make([][]byte, 0)
	params = append(params, common.Uint64ToBytes(1102))
	params = append(params, []byte("GetVerifiersList"))
	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, params)
	if err != nil {
		fmt.Println(err)
		t.Errorf("GetVerifiersList encode rlp data fail")
	} else {
		fmt.Println("GetVerifiersList data rlp: ", hexutil.Encode(buf.Bytes()))
	}

}

func TestIsElection (t *testing.T) {

	num230 := big.NewInt(230)
	fmt.Println(xutil.IsElection(num230.Uint64()))

	num1 := big.NewInt(1)
	fmt.Println(xutil.IsElection(num1.Uint64()))


	num480 := big.NewInt(480)
	fmt.Println(xutil.IsElection(num480.Uint64()))

	num231 := big.NewInt(231)
	fmt.Println(xutil.IsElection(num231.Uint64()))
}