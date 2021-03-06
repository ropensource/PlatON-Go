package plugin

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/byteutil"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
	gerr "github.com/go-errors/errors"
	"reflect"
)

type BasePlugin interface {
	BeginBlock(blockHash common.Hash, header *types.Header, state xcom.StateDB) (bool, error)
	EndBlock(blockHash common.Hash, header *types.Header, state xcom.StateDB) (bool, error)
	Confirmed(block *types.Block) error
}



var (
	DecodeTxDataErr = errors.New("decode tx data is err")
	FuncNotExistErr = errors.New("the func is not exist")
	FnParamsLenErr 	= errors.New("the params len and func params len is not equal")
)

func  Verify_tx_data(input []byte, command map[uint16]interface{} ) (fn interface{}, FnParams []reflect.Value, err error)  {

	defer func() {
		if er := recover(); nil != er {
			fn, FnParams, err = nil, nil,fmt.Errorf("parse tx data is panic: %s", gerr.Wrap(er, 2).ErrorStack())
			log.Error("Failed to Verify Staking tx data", "error", err)
		}
	}()

	var args [][]byte
	if err := rlp.Decode(bytes.NewReader(input), &args); nil != err {
		return nil, nil, DecodeTxDataErr
	}

	if fn, ok := command[byteutil.BytesToUint16(args[0])]; !ok {
			return nil, nil, FuncNotExistErr
	}else {

		// the func params type list
		paramList := reflect.TypeOf(fn)
		// the func params len
		paramNum := paramList.NumIn()

		if paramNum != len(args)-1 {
			return nil, nil, FnParamsLenErr
		}

		params := make([]reflect.Value, paramNum)

		for i := 0; i < paramNum; i++ {
			targetType := paramList.In(i).String()
			inputByte := []reflect.Value{reflect.ValueOf(args[i+1])}
			params[i] = reflect.ValueOf(byteutil.Bytes2X_CMD[targetType]).Call(inputByte)[0]
		}

		return fn, params, nil
	}
}