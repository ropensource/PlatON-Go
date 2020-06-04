package main

import (
	"encoding/json"
	"github.com/PlatONnetwork/PlatON-Go/cmd/utils"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"io/ioutil"

	"gopkg.in/urfave/cli.v1"
)

var HexAccountFileFlag = cli.StringFlag{
	Name:  "hexAddressFile",
	Usage: "file hex accounts want to update to bech32 address,file like  [hex,hex...]",
}

type addressPair struct {
	Address       common.AddressOutput
	OriginAddress string
}

type KeyPair struct {
	PrivateKey string `json:"private_key"`
	Address string `json:"address"`
}

var commandAddressHexToBech32 = cli.Command{
	Name:      "updateaddress",
	Usage:     "update hex address to bech32 address",
	ArgsUsage: "[<hex_address> <hex_address>...]",
	Description: `
update hex address to bech32 address.
`,
	Flags: []cli.Flag{
		jsonFlag,
		HexAccountFileFlag,
	},
	Action: func(ctx *cli.Context) error {
		var accounts []KeyPair
		if ctx.IsSet(HexAccountFileFlag.Name) {
			accountPath := ctx.String(HexAccountFileFlag.Name)
			accountjson, err := ioutil.ReadFile(accountPath)
			if err != nil {
				utils.Fatalf("Failed to read the keyfile at '%s': %v", accountPath, err)
			}
			if err := json.Unmarshal(accountjson, &accounts); err != nil {
				utils.Fatalf("Failed to json decode '%s': %v", accountPath, err)
			}
		} else {
			/*for _, add := range ctx.Args() {
				if add == "" {
					utils.Fatalf("the account can't be nil")
				}
				accounts = append(accounts, add)
			}*/
		}
		var outAddress []KeyPair
		genesis := make(map[string]map[string]string)
		for _, account := range accounts {
			address := common.HexToAddress(account.Address)
			out := KeyPair{
				PrivateKey: account.PrivateKey,
				Address:       address.Bech32WithPrefix(common.MainNetAddressPrefix),
			}
			outAddress = append(outAddress, out)
			genesis[out.Address] = map[string]string{
				"balance": "0x200000000000000000000000000000000000000000000000000000000000",
			}
		}

		if ctx.Bool(jsonFlag.Name) {
			//mustPrintJSON(outAddress)
			data, _ := json.Marshal(outAddress)
			ioutil.WriteFile("C:\\Users\\DELL\\Desktop\\0.13.0\\bech32_all_addr_and_private_keys.json",
				data, 0600)

			// 产出创始区块文件
			genesisData, _ := json.Marshal(genesis)
			ioutil.WriteFile("C:\\Users\\DELL\\Desktop\\0.13.0\\bech_genesis_alloc_10000.json", genesisData, 0600)

		} else {
			/*for i, address := range outAddress {
				fmt.Println("originAddress: ", address.OriginAddress)
				address.Address.Print()
				if i != len(outAddress)-1 {
					fmt.Println("---")
				}
			}*/
		}
		return nil
	},
}
