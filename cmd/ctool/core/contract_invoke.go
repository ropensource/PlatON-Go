package core

import (
	"Platon-go/cmd/ctool/rlp"
	"Platon-go/common/hexutil"
	"encoding/json"
	"fmt"
)

/**
  合约调用入口
*/
func ContractInvoke(contractAddr string, abiPath string, funcParams string, configPath string) {
	config := Config{}
	parseConfigJson(configPath, &config)

	//判断该合约是否存在
	if !getContractByAddress(contractAddr, config.Url) {
		fmt.Printf("the contract address is not exist ...")
		return
	}

	//解析调用的方法 参数
	funcName, inputParams := GetFuncNameAndParams(funcParams)

	//判断该方法是否存在
	abiFunc := parseFuncFromAbi(abiPath, funcName)
	if abiFunc.Method == "" {
		fmt.Printf("the function not exist ,func= %s\n", funcName)
		return
	}

	//判断参数是否正确
	if len(abiFunc.Args) != len(inputParams) {
		fmt.Printf("incorrect number of parameters ,request=%d,get=%d\n", len(abiFunc.Args), len(inputParams))
		return
	}

	//todo 参数类型校验

	paramArr := [][]byte{
		Int32ToBytes(invoke),
		[]byte(funcName),
	}

	for i, v := range inputParams {
		f := abiFunc.Args[i]
		p, e := StringConverter(v, f.RealTypeName)
		if e != nil {
			fmt.Printf("incorrect param type: %s,index:%d", v, i)
		}
		paramArr = append(paramArr, p)
	}

	paramBytes, e := rlp.EncodeToBytes(paramArr)
	if e != nil {
		fmt.Printf("rpl encode error...\n %s", e.Error())
		panic(e)
	}

	params := TxParams{
		From:     config.From,
		To:       contractAddr,
		GasPrice: config.GasPrice,
		Gas:      config.Gas,
		Data:     hexutil.Encode(paramBytes),
	}

	var r string
	var err error
	//是否走call
	if abiFunc.FuncType == "const" {
		paramList := make(List, 2)
		paramList[0] = params
		paramList[1] = "latest"
		r, err = Send(paramList, "eth_call", config.Url)
	} else {
		paramList := make(List, 1)
		paramList[0] = params
		r, err = Send(paramList, "eth_sendTransaction", config.Url)
	}

	if err != nil {
		fmt.Printf("send http post to invoke contract error ")
		return
	}

	var resp = Response{}
	err = json.Unmarshal([]byte(r), &resp)
	if err != nil {
		fmt.Printf("parse eth_sendTransaction result error ! \n %s", err.Error())
		return
	}

	if resp.Error.Code != 0 {
		fmt.Printf("eth_sendTransaction error ,error:%v", resp.Error.Message)
		return
	}

	//根据abi 返回类型判断解析什么类型
	if abiFunc.FuncType == "const" {
		bytes, _ := hexutil.Decode(resp.Result)
		result := BytesConverter(bytes, abiFunc.Return)
		fmt.Printf("\nresult: %v\n", result)
	} else {
		fmt.Printf("\ntrasaction hash: %s\n", resp.Result)
	}

}

/**
  通过eth_getCode判断合约是否存在
*/
func getContractByAddress(addr, url string) bool {

	params := []string{addr, "latest"}
	r, err := Send(params, "eth_getCode", url)
	if err != nil {
		fmt.Printf("send http post to get contract address error ")
		return false
	}

	var resp = Response{}
	err = json.Unmarshal([]byte(r), &resp)
	if err != nil {
		fmt.Printf("parse eth_getCode result error ! \n %s", err.Error())
		return false
	}

	if resp.Error.Code != 0 {
		fmt.Printf("eth_getCode error ,error:%v", resp.Error.Message)
		return false
	}
	//fmt.Printf("trasaction hash: %s\n", resp.Result)

	if resp.Result != "" && len(resp.Result) > 2 {
		return true
	} else {
		return false
	}
}
