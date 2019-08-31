package vm

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/mock"
	cvm "github.com/PlatONnetwork/PlatON-Go/common/vm"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/crypto/bls"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/PlatONnetwork/PlatON-Go/x/gov"
	"github.com/PlatONnetwork/PlatON-Go/x/plugin"
	"github.com/PlatONnetwork/PlatON-Go/x/restricting"
	"github.com/PlatONnetwork/PlatON-Go/x/staking"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
	"github.com/PlatONnetwork/PlatON-Go/x/xutil"
)

func init() {
	//log.Root().SetHandler(log.CallerFileHandler(log.LvlFilterHandler(log.Lvl(4), log.StreamHandler(os.Stderr, log.TerminalFormat(true)))))
	bls.Init(bls.BLS12_381)
}

const initGas = 10000000

var (
	nodeIdArr = []discover.NodeID{
		discover.MustHexID("0x362003c50ed3a523cdede37a001803b8f0fed27cb402b3d6127a1a96661ec202318f68f4c76d9b0bfbabfd551a178d4335eaeaa9b7981a4df30dfc8c0bfe3384"),
		discover.MustHexID("0xced880d4769331f47af07a8d1b79de1e40c95a37ea1890bb9d3f0da8349e1a7c0ea4cadbb9c5bf185b051061eef8e5eadca251c24e1db1d9faf0fb24cbd06f9a"),
		discover.MustHexID("0xda56501a77fc1dfe0399b81f3909061d9a176cb9433fab4d3dfb1a10344c243274e38155e18878c7a0b3fcdd6182000c7784a95e2c4d9e0691ce67798624786e"),
		discover.MustHexID("0x89a4409abe1ace8b77c4497c2073a8a2046dbdabb58c8bb58fe73926bbdc572fb848d739b1d2d09dd0796abcc1ed8d9a33bb3ef0a6c2e106e408090df179b041"),
		discover.MustHexID("0x65e2ab09161e32e6d07d82adaa416ee6d41d617c52db20e3145a4d1b7d396af38d095c87508ad5bb35df741513bdc4bf12fec215e58450e255f05d194d41d089"),
		discover.MustHexID("0x9bfacd628f3adb0f94e8b3968064d5248fa18efa75c680fdffea3af2575406461f3395817dd2a1be07a79bd81ffa00f57ad82286061d4a6caceece048e352380"),
		discover.MustHexID("0x1e07d66b56bbc931ddce7cc5b9f55672d7fe4e19897a42f19d4ad7c969435cad652d720401d68f5769e245ec0f4e23362c8b1b062771d614876fdbb875ba9d44"),
		discover.MustHexID("0x11a315747ce79cdf3d6aaf87ff2b6897950a20bda281838f922ea9407736fec9029d85f6202fd059a57a9119d05895402e7570948ae759cb093a54c3da9e0a4a"),
		discover.MustHexID("0x248af08a775ff63a47a5970e4928bcccd1a8cef984fd4142ea7f89cd13015bdab9ca4a8c5e1070dc00fa81a047542f53ca596f553c4acfb7abe75a8fb5019057"),
		discover.MustHexID("0xfd790ff5dc48baccb9418ce5cfac6a10c3646f20a3fe32d9502c4edce3a77fa90bfee0361d8a72093b7994f8cbc28ee537bdda2b634c5966b1a9253d9d270145"),
		discover.MustHexID("0x56d243db84a521cb204f582ee84bca7f4af29437dd447a6e36d17f4853888e05343844bd64294b99b835ca7f72ef5b1325ef1c89b0c5c2744154cdadf7c4e9fa"),
		discover.MustHexID("0x8796a6fcefd9037d8433e3a959ff8f3c4552a482ce727b00a90bfd1ec365ce2faa33e19aa6a172b5c186b51f5a875b5acd35063171f0d9501a9c8f1c98513825"),
		discover.MustHexID("0x547b876036165d66274ce31692165c8acb6f140a65cab0e0e12f1f09d1c7d8d53decf997830919e4f5cacb2df1adfe914c53d22e3ab284730b78f5c63a273b8c"),
		discover.MustHexID("0x9fdbeb873bea2557752eabd2c96419b8a700b680716081472601ddf7498f0db9b8a40797b677f2fac541031f742c2bbd110ff264ae3400bf177c456a76a93d42"),
		discover.MustHexID("0xc553783799bfef7c34a84b2737f2c77f8f2c5cfedc3fd7af2d944da6ece90aa94cf621e6de5c4495881fbfc9beec655ffb10e39cb4ca9be7768d284409040f32"),
		discover.MustHexID("0x75ad2ee8ca77619c3ba0ddcec5dab1375fe4fa90bab9e751caef3996ce082dfed32fe4c137401ee05e501c079b2e4400397b09de14b08b09c9e7f9698e9e4f0a"),
		discover.MustHexID("0xdb18af9be2af9dff2347c3d06db4b1bada0598d099a210275251b68fa7b5a863d47fcdd382cc4b3ea01e5b55e9dd0bdbce654133b7f58928ce74629d5e68b974"),
		discover.MustHexID("0x472d19e5e9888368c02f24ebbbe0f2132096e7183d213ab65d96b8c03205f88398924af8876f3c615e08aa0f9a26c38911fda26d51c602c8d4f8f3cb866808d7"),
		discover.MustHexID("4f1f036e5e18cc812347d5073cbec2a8da7930de323063c39b0d4413a396e088bfa90e8c28174313d8d82e9a14bc0884b13a48fc28e619e44c48a49b4fd9f107"),
		discover.MustHexID("f18c596232d637409c6295abb1e720db99ffc12363a1eb8123d6f54af80423a5edd06f91115115a1dca1377e97b9031e2ddb864d34d9b3491d6fa07e8d9b951b"),
		discover.MustHexID("7a8f7a28ac1c4eaf98b2be890f372e5abc58ebe6d3aab47aedcb0076e34eb42882e926676ebab327a4ef4e2ea5c4296e9c7bc0991360cb44f52672631012db1b"),
		discover.MustHexID("9eeb448babf9e93449e831b91f98d9cbc0c2324fe8c43baac69d090717454f3f930713084713fe3a9f01e4ca59b80a0f2b41dbd6d531f414650bab0363e3691a"),
		discover.MustHexID("cc1d7314c15e30dc5587f675eb5f803b1a2d88bfe76cec591cec1ff678bc6abce98f40054325bdcb44fb83174f27d38a54fbce4846af8f027b333868bc5144a4"),
		discover.MustHexID("e4d99694be2fc8a53d8c2446f947aec1c7de3ee26f7cd43f4f6f77371f56f11156218dec32b51ddce470e97127624d330bb7a3237ba5f0d87d2d3166faf1035e"),
		discover.MustHexID("9c61f59f70296b6d494e7230888e58f19b13c5c6c85562e57e1fe02d0ff872b4957238c73559d017c8770b999891056aa6329dbf628bc19028d8f4d35ec35823"),
	}

	addrArr = []common.Address{
		common.HexToAddress("0xc9E1C2B330Cf7e759F2493c5C754b34d98B07f93"),
		common.HexToAddress("0xd87E10F8efd2C32f5e88b7C279953aEF6EE58902"),
		common.HexToAddress("0xeAEc60C738eeD9468e6AcCc1d403faCF1A670F6D"),
		common.HexToAddress("0x5c5994165265Ac31AAFE874a231f2C5d0eF29C3a"),
		common.HexToAddress("0xB9449Eb226cb93c3BF5FeCA16c85a737538e24f0"),
		common.HexToAddress("0x908bad1823BddA66cc65E788b9d0194b7975976A"),
		common.HexToAddress("0x3DfC64A87db521662675DffEa48d0c208414D4f8"),
		common.HexToAddress("0xad8adf35068Cdf572c9eFb5a069dA48D2E165Aa1"),
		common.HexToAddress("0xf33b5Da47c6ECbC61cF07C7387Afc6ef0EA2f866"),
		common.HexToAddress("0x2E5FB4F78E3FB9b1898DE7d7D8dB3d44C62040be"),
		common.HexToAddress("0x285CF84ea3E177E1fC9F396aEbc9329a08f51bb5"),
		common.HexToAddress("0x91BffdC88329AfDD97DF6fe92cfd4FcB7927Aecd"),
		common.HexToAddress("0x58b62FfF5046aF2252F1F8Ecb5c3342ada394F72"),
		common.HexToAddress("0x8ec116c11d8515e8222Cabc4BEc06A880C51D929"),
		common.HexToAddress("0x364eCBade4c35beE2F8a04F8209BaB236B48A35a"),
		common.HexToAddress("0x26896c394A1E12095e822e5b080e8EfA050c738C"),
		common.HexToAddress("0x314253824CD6b7BCF1613CAB00126D6076F7a389"),
		common.HexToAddress("0x5544F05D51E45fa6497AFEC0F1A5d64531B21be0"),
		common.HexToAddress("0x3da830FAd2A6983d948d7262B2AdF7eA53b953be"),
		common.HexToAddress("0x815A7910C035F2FB9451cDA349969788449c2288"),
		common.HexToAddress("0x4Cdd49e08587c824c7629e7d124390B70d105740"),
		common.HexToAddress("0xD041b5fAaa4B721241A55107FE8F19ce1ba3E7fD"),
		common.HexToAddress("0xcbc583DEdbbE6b51B86036C040596bB0a0299a73"),
		common.HexToAddress("0x1c0A4509Ba46deA47775Ad8B20A19f398B820642"),
		common.HexToAddress("0xEEE10Fc4A3AB339f5a788f3b82Eb57738F075EcE"),
	}

	priKeyArr = []*ecdsa.PrivateKey{
		crypto.HexMustToECDSA("0c6ccec28e36dc5581ea3d8af1303c774b51523da397f55cdc4acd9d2b988132"),
		crypto.HexMustToECDSA("07c0b2525cbff7dad6211cf901507e3814a77d864d31bdaa5785a94ee20a8da1"),
		crypto.HexMustToECDSA("564a4965c2bd98654c275c6b63713c936f2dc91bb6a91bdd47e8320d4d9ebcf4"),
		crypto.HexMustToECDSA("dfdbb19f4c18bb4964392b6e1998c62c9e2b53e9400c4dd64d6659c2191625f5"),
		crypto.HexMustToECDSA("2ceb44fe9196cd1ad2c4bb4c657098fdf8baeb85910011243f6da23a47ea9781"),
		crypto.HexMustToECDSA("904160f823ede58e83584f0e2b98f0994fc9626f250457873b16570a446c9e92"),
		crypto.HexMustToECDSA("9ed9a0d08f8354539e336f85b3011d609ded672bba12b63d62931b111471cb99"),
		crypto.HexMustToECDSA("85e733ec3f15aab14848e21c1ac69624ce8547c92d3f0ee8abc0f3ab412158a1"),
		crypto.HexMustToECDSA("343d10559147d42e1632b4e932aeae36e360d3e0083b9d8d30bb8cc9bb6923c1"),
		crypto.HexMustToECDSA("15439211a0e25c58d7985e11138ce60f675e5243e2b4387fadbd6a0c85755791"),
		crypto.HexMustToECDSA("4a931cfc05fd33b3f3b0f3d910b4358b4cfeac6e1f13b3461a56945ab0de8d96"),
		crypto.HexMustToECDSA("72c8e5bc83fd79debd0af75dab09617198c5f06656ef24009bf7e9a944750bd2"),
		crypto.HexMustToECDSA("d58b015ad107166bd648ba3fb15672e4958f8df668d85acacda7a2ed6f855683"),
		crypto.HexMustToECDSA("1fa19b3862cb9ec584da03d56a84766abdc03cbb3a5e07645531563c1fe2ede6"),
		crypto.HexMustToECDSA("a2be5c2766e9eeed2575448364313cfa91caeb1f1fd03cdbe6f9cee1ded2bffa"),
		crypto.HexMustToECDSA("7da86d7aca8b5dbec9d0bd3c0c2e91552f504df3a42a6e4493992b251bc6c438"),
		crypto.HexMustToECDSA("ed46c6521237ffba7626c67574f8e29d2941ef4bdef561e6d2b4bc877f7c4745"),
		crypto.HexMustToECDSA("b5f8a8bff108a3e674eef019121bdb1c1e0c14857888ff4052954db5700520c3"),
		crypto.HexMustToECDSA("548ceef29a39093e48ef65bc98b210320dedd79ca40acebeb573f8eb72018aac"),
		crypto.HexMustToECDSA("73a2bd8694f883ff5f11551c04303ff7180ae6ef1b89170a67ace10d04c7c3e2"),
		crypto.HexMustToECDSA("996e2bb9c1371e50125fb8b1d0e6f9c46148dfb8b01d9edd6e8b5ec1a6241316"),
		crypto.HexMustToECDSA("51c977a01d5517406fcce2bf7bbb44c67e6b876641a5dac6d2fc26b2f6a97001"),
		crypto.HexMustToECDSA("41d4ce3f8b18fc7ccb4bb0e9514e0863d0c0bd4bb26e9fba3c2a384189c2000b"),
		crypto.HexMustToECDSA("3653b25ba39e59d12a3f45f0fb324b8588db839de4bafd9b938315c356a37051"),
		crypto.HexMustToECDSA("e066f9c4daabcc354162165f8aa161c0bc1cede1b0d14a269f63f6d6bdb1ec5d"),
	}

	blockNumber = big.NewInt(1)
	blockHash   = common.HexToHash("9d4fb5346abcf593ad80a0d3d5a371b22c962418ad34189d5b1b39065668d663")

	blockNumber2 = big.NewInt(2)
	blockHash2   = common.HexToHash("c95876b92443d652d7eb7d7a9c0e2c58a95e934c0c1197978c5445180cc60980")

	blockNumber3 = big.NewInt(3)
	blockHash3   = common.HexToHash("c95876b92443d652d7eb7d7a9c0e2c58a95e934c0c1197978c5445180cc60345")

	lastBlockNumber uint64
	lastBlockHash   common.Hash
	lastHeader      types.Header

	sender            = common.HexToAddress("0xeef233120ce31b3fac20dac379db243021a5234")
	anotherSender     = common.HexToAddress("0xeef233120ce31b3fac20dac379db243021a5233")
	sender_balance, _ = new(big.Int).SetString("9999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999", 10)

	delegate_sender            = common.HexToAddress("0xc1f330b214668beac2e6418dd651b09c759a4bf5")
	delegate_sender_balance, _ = new(big.Int).SetString("9999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999", 10)

	sndb = snapshotdb.Instance()

	txHashArr = []common.Hash{
		common.HexToHash("0x00000000000000000000000000000000000000886d5ba2d3dfb2e2f6a1814f22"),
		common.HexToHash("0x000000000000000000000000000000005249b59609286f2fa91a2abc8555e887"),
		common.HexToHash("0x000000008dba388834e2515c4d9ccb02a48bae177e73959330e55067211c2456"),
		common.HexToHash("0x0000000000000000000000000000000000009a715a765a72b8a289156f9543c9"),
		common.HexToHash("0x0000e1b4a5508c11772b61f463657585c33b577019e4a23bd359c018a4e306d1"),
		common.HexToHash("0x00fd854f940e2d2af8e74c33e640ea6f75c1d9ee49b816b8a4647611d0c91863"),
		common.HexToHash("0x0000000000001038575739a53385cfe42321585a56050e18f8ea2b3e8dc21966"),
		common.HexToHash("0x0000000000000000000000000000000000000048f3b312dc8d081e1186abe8c2"),
		common.HexToHash("0x000000000000000000000000f5bd37579e7ca954eba8fbe7a65646250e92ab7d"),
		common.HexToHash("0x00000000000000000000000000000000000000001d65a5a69fed6ddb0cb58dff"),
		common.HexToHash("0x00000000000000000000000000000000000000000000000000000000000000d2"),
		common.HexToHash("0x0000000000000000000000000000000000000000000000000000f2e8b2706c9e"),
		common.HexToHash("0x00000000000000000000000000e22a393898aac376b079e0894e8e2be6024d03"),
		common.HexToHash("0x000000000000000000000000000000000000000000000000483570dd0679860a"),
		common.HexToHash("0x000000000000000000000000000000000000007fc9e1dc435b5d0064ac50fd4e"),
		common.HexToHash("0x00000000000000000000000000cbeb8f4d51969d7eb70a4f6e8505950d870df7"),
		common.HexToHash("0x00000000000000000000000000000000000000000000000000000000000000b4"),
		common.HexToHash("0x000000008fd2abdf28d87efb2c7fa2d37618c8dba97059376d6a58007bee3d8b"),
		common.HexToHash("0x0000000000000000000000003566f3a0adf49d90e610ef3d3548b5a72b1fe199"),
		common.HexToHash("0x00000000000054fa3d19eb57e98aa1dd69d216722054d8539ede4b89c5b77ee9"),
		common.HexToHash("0x00000000000000000000000000cbeb8f4d51969d7eb70a4f6e8505950d870ef3"),
		common.HexToHash("0x00000000000000000000000000000000000000000000000000000000000011b4"),
		common.HexToHash("0x000000008fd2abdf28d87efb2c7fa2d37618c8dba97059376d6a58007bee3d84"),
		common.HexToHash("0x0000000000000000000000003566f3a0adf49d90e610ef3d3548b5a72b1fe178"),
		common.HexToHash("0x00000000000054fa3d19eb57e98aa1dd69d216722054d8539ede4b89c5b77ee5"),
	}

	//initProgramVersion      = uint32(params.VersionMajor<<16 | params.VersionMinor<<8 | params.VersionPatch)
	//initProgramVersionBytes = common.Uint32ToBytes(initProgramVersion)

	initProgramVersion      = uint32(1<<16 | 0<<8 | 0) // 65536, version: 1.0.0
	initProgramVersionBytes = common.Uint32ToBytes(initProgramVersion)

	promoteVersion      = uint32(2<<16 | 0<<8 | 0) // 131072, version: 2.0.0
	promoteVersionBytes = common.Uint32ToBytes(promoteVersion)

	balanceStr = []string{
		"90000000000000000000000000",
		"600000000000000000000000000",
		"13000000000000000000000000",
		"11000000000000000000000000",
		"10000000000000000000000000",
		"48790000000000000000000000",
		"18000000000000000000000000",
		"10000000000000000000000000",
		"10000000000000000000000000",
		"700000000000000000000000000",
		"55500000000000000000000000",
		"90000000000000000000000000",
		"600000000000000000000000000",
		"13000000000000000000000000",
		"11000000000000000000000000",
		"10000000000000000000000000",
		"48790000000000000000000000",
		"18000000000000000000000000",
		"10000000000000000000000000",
		"10000000000000000000000000",
		"700000000000000000000000000",
		"55500000000000000000000000",
		"10000000000000000000000000",
		"700000000000000000000000000",
		"55500000000000000000000000",
	}

	nodeNameArr = []string{
		"PlatON",
		"Gavin",
		"Emma",
		"Kally",
		"Juzhen",
		"Baidu",
		"Alibaba",
		"Tencent",
		"ming",
		"hong",
		"gang",
		"guang",
		"hua",
		"PlatON_2",
		"Gavin_2",
		"Emma_2",
		"Kally_2",
		"Juzhen_2",
		"Baidu_2",
		"Alibaba_2",
		"Tencent_2",
		"ming_2",
		"hong_2",
		"gang_2",
		"guang_2",
	}

	chaList = []string{"A", "a", "B", "b", "C", "c", "D", "d", "E", "e", "F", "f", "G", "g", "H", "h", "J", "j", "K", "k", "M", "m",
		"N", "n", "P", "p", "Q", "q", "R", "r", "S", "s", "T", "t", "U", "u", "V", "v", "W", "w", "X", "x", "Y", "y", "Z", "z"}

	specialCharList = []string{
		"☄", "★", "☎", "☻", "♨", "✠", "❝", "♚", "♘", "✎", "♞", "✩", "✪", "❦", "❥", "❣", "웃", "卍", "Ⓞ", "▶", "◙", "⊕", "◌", "⅓", "∭",
		"∮", "╳", "㏒", "㏕", "‱", "㎏", "❶", "Ň", "🅱", "🅾", "𝖋", "𝕻", "𝕼", "𝕽", "お", "な", "ぬ", "㊎", "㊞", "㊮", "✘"}
)

func newPlugins() {
	plugin.GovPluginInstance()
	plugin.StakingInstance()
	plugin.SlashInstance()
	plugin.RestrictingInstance()
	plugin.RewardMgrInstance()

	snapshotdb.Instance()
}

func newChainState() (*mock.MockStateDB, *types.Block, error) {

	testGenesis := new(types.Block)
	chain := mock.NewChain(testGenesis)
	//	var state *state.StateDB

	chain.StateDB.AddBalance(sender, sender_balance)
	chain.StateDB.AddBalance(delegate_sender, delegate_sender_balance)
	for i, addr := range addrArr {
		amount, _ := new(big.Int).SetString(balanceStr[len(addrArr)-1-i], 10)
		amount = new(big.Int).Mul(common.Big257, amount)
		chain.StateDB.AddBalance(addr, amount)
	}
	return chain.StateDB, chain.Genesis, nil
}

func newEvm(blockNumber *big.Int, blockHash common.Hash, state *mock.MockStateDB) *EVM {

	if nil == state {
		state, _, _ = newChainState()
	}
	evm := &EVM{
		StateDB: state,
	}
	context := Context{
		BlockNumber: blockNumber,
		BlockHash:   blockHash,
	}
	evm.Context = context

	//set a default active version

	gov.AddActiveVersion(initProgramVersion, 0, state)

	return evm
}

func newContract(value *big.Int, sender common.Address) *Contract {
	callerAddress := AccountRef(sender)
	fmt.Println("newContract sender :", callerAddress.Address().Hex())
	contract := NewContract(callerAddress, callerAddress, value, uint64(initGas))
	return contract
}

func build_staking_data(genesisHash common.Hash) {

	stakingDB := staking.NewStakingDB()
	sndb.NewBlock(big.NewInt(1), genesisHash, blockHash)
	// MOCK

	nodeId_A := nodeIdArr[0]
	addr_A, _ := xutil.NodeId2Addr(nodeId_A)

	nodeId_B := nodeIdArr[1]
	addr_B, _ := xutil.NodeId2Addr(nodeId_B)

	nodeId_C := nodeIdArr[2]
	addr_C, _ := xutil.NodeId2Addr(nodeId_C)

	//canArr := make(staking.CandidateQueue, 0)

	var blsKey1 bls.SecretKey
	blsKey1.SetByCSPRNG()
	c1 := &staking.Candidate{
		NodeId:             nodeId_A,
		BlsPubKey:          *blsKey1.GetPublicKey(),
		StakingAddress:     sender,
		BenefitAddress:     addrArr[0],
		StakingTxIndex:     uint32(2),
		ProgramVersion:     uint32(1),
		Status:             staking.Valided,
		StakingEpoch:       uint32(1),
		StakingBlockNum:    uint64(1),
		Shares:             common.Big256,
		Released:           common.Big2,
		ReleasedHes:        common.Big32,
		RestrictingPlan:    common.Big1,
		RestrictingPlanHes: common.Big257,
		Description: staking.Description{
			ExternalId: "xxccccdddddddd",
			NodeName:   "I Am " + fmt.Sprint(1),
			Website:    "www.baidu.com",
			Details:    "this is  baidu ~~",
		},
	}

	var blsKey2 bls.SecretKey
	blsKey2.SetByCSPRNG()
	c2 := &staking.Candidate{
		NodeId:             nodeId_B,
		BlsPubKey:          *blsKey2.GetPublicKey(),
		StakingAddress:     sender,
		BenefitAddress:     addrArr[1],
		StakingTxIndex:     uint32(3),
		ProgramVersion:     uint32(1),
		Status:             staking.Valided,
		StakingEpoch:       uint32(1),
		StakingBlockNum:    uint64(1),
		Shares:             common.Big256,
		Released:           common.Big2,
		ReleasedHes:        common.Big32,
		RestrictingPlan:    common.Big1,
		RestrictingPlanHes: common.Big257,
		Description: staking.Description{
			ExternalId: "SFSFSFSFSFSFSSFS",
			NodeName:   "I Am " + fmt.Sprint(2),
			Website:    "www.JD.com",
			Details:    "this is  JD ~~",
		},
	}

	var blsKey3 bls.SecretKey
	blsKey3.SetByCSPRNG()
	c3 := &staking.Candidate{
		NodeId:             nodeId_C,
		BlsPubKey:          *blsKey3.GetPublicKey(),
		StakingAddress:     sender,
		BenefitAddress:     addrArr[2],
		StakingTxIndex:     uint32(4),
		ProgramVersion:     uint32(1),
		Status:             staking.Valided,
		StakingEpoch:       uint32(1),
		StakingBlockNum:    uint64(1),
		Shares:             common.Big256,
		Released:           common.Big2,
		ReleasedHes:        common.Big32,
		RestrictingPlan:    common.Big1,
		RestrictingPlanHes: common.Big257,
		Description: staking.Description{
			ExternalId: "FWAGGDGDGG",
			NodeName:   "I Am " + fmt.Sprint(3),
			Website:    "www.alibaba.com",
			Details:    "this is  alibaba ~~",
		},
	}

	//canArr = append(canArr, c1)
	//canArr = append(canArr, c2)
	//canArr = append(canArr, c3)

	stakingDB.SetCanPowerStore(blockHash, addr_A, c1)
	stakingDB.SetCanPowerStore(blockHash, addr_B, c2)
	stakingDB.SetCanPowerStore(blockHash, addr_C, c3)

	stakingDB.SetCandidateStore(blockHash, addr_A, c1)
	stakingDB.SetCandidateStore(blockHash, addr_B, c2)
	stakingDB.SetCandidateStore(blockHash, addr_C, c3)

	queue := make(staking.ValidatorQueue, 0)

	v1 := &staking.Validator{
		NodeAddress:   addr_A,
		NodeId:        c1.NodeId,
		BlsPubKey:     c1.BlsPubKey,
		StakingWeight: [staking.SWeightItem]string{"1", common.Big256.String(), fmt.Sprint(c1.StakingBlockNum), fmt.Sprint(c1.StakingTxIndex)},
		ValidatorTerm: 0,
	}

	v2 := &staking.Validator{
		NodeAddress:   addr_B,
		NodeId:        c2.NodeId,
		BlsPubKey:     c2.BlsPubKey,
		StakingWeight: [staking.SWeightItem]string{"1", common.Big256.String(), fmt.Sprint(c2.StakingBlockNum), fmt.Sprint(c2.StakingTxIndex)},
		ValidatorTerm: 0,
	}

	v3 := &staking.Validator{
		NodeAddress:   addr_C,
		NodeId:        c3.NodeId,
		BlsPubKey:     c3.BlsPubKey,
		StakingWeight: [staking.SWeightItem]string{"1", common.Big256.String(), fmt.Sprint(c3.StakingBlockNum), fmt.Sprint(c3.StakingTxIndex)},
		ValidatorTerm: 0,
	}

	queue = append(queue, v1)
	queue = append(queue, v2)
	queue = append(queue, v3)

	epoch_Arr := &staking.Validator_array{
		Start: 1,
		End:   uint64(xutil.CalcBlocksEachEpoch()),
		Arr:   queue,
	}

	pre_Arr := &staking.Validator_array{
		Start: 0,
		End:   0,
		Arr:   queue,
	}

	curr_Arr := &staking.Validator_array{
		Start: 1,
		End:   uint64(xutil.ConsensusSize()),
		Arr:   queue,
	}

	setVerifierList(blockHash, epoch_Arr)
	setRoundValList(blockHash, pre_Arr)
	setRoundValList(blockHash, curr_Arr)

	lastBlockHash = blockHash
	lastBlockNumber = blockNumber.Uint64()
	lastHeader = types.Header{
		Number: blockNumber,
	}
}

func build_staking_data_more(block uint64) {

	no := int64(block)
	header := types.Header{
		Number: big.NewInt(no),
	}
	hash := header.Hash()

	stakingDB := staking.NewStakingDB()
	sndb.NewBlock(big.NewInt(int64(block)), lastBlockHash, hash)
	// MOCK

	validatorArr := make(staking.ValidatorQueue, 0)

	// build  more data
	for i := 0; i < 1000; i++ {

		var index int
		if i >= len(balanceStr) {
			index = i % (len(balanceStr) - 1)
		}

		balance, _ := new(big.Int).SetString(balanceStr[index], 10)

		rand.Seed(time.Now().UnixNano())

		weight := rand.Intn(1000000000)

		ii := rand.Intn(len(chaList))

		balance = new(big.Int).Add(balance, big.NewInt(int64(weight)))

		randBuildFunc := func() (discover.NodeID, common.Address, error) {
			privateKey, err := crypto.GenerateKey()
			if nil != err {
				fmt.Printf("Failed to generate random NodeId private key: %v", err)
				return discover.NodeID{}, common.ZeroAddr, err
			}

			nodeId := discover.PubkeyID(&privateKey.PublicKey)

			privateKey, err = crypto.GenerateKey()
			if nil != err {
				fmt.Printf("Failed to generate random Address private key: %v", err)
				return discover.NodeID{}, common.ZeroAddr, err
			}

			addr := crypto.PubkeyToAddress(privateKey.PublicKey)

			return nodeId, addr, nil
		}

		var nodeId discover.NodeID
		var addr common.Address

		if i < 25 {
			nodeId = nodeIdArr[i]
			ar, _ := xutil.NodeId2Addr(nodeId)
			addr = ar
		} else {
			id, ar, err := randBuildFunc()
			if nil != err {
				return
			}
			nodeId = id
			addr = ar
		}

		var blsKey bls.SecretKey
		blsKey.SetByCSPRNG()
		canTmp := &staking.Candidate{
			NodeId:          nodeId,
			BlsPubKey:       *blsKey.GetPublicKey(),
			StakingAddress:  sender,
			BenefitAddress:  addr,
			StakingBlockNum: uint64(1),
			StakingTxIndex:  uint32(i + 1),
			Shares:          balance,
			ProgramVersion:  xutil.CalcVersion(initProgramVersion),
			// Prevent null pointer initialization
			Released:           common.Big0,
			ReleasedHes:        common.Big0,
			RestrictingPlan:    common.Big0,
			RestrictingPlanHes: common.Big0,

			Description: staking.Description{
				NodeName:   nodeNameArr[index] + "_" + fmt.Sprint(i),
				ExternalId: nodeNameArr[index] + chaList[(len(chaList)-1)%(index+ii+1)] + "balabalala" + chaList[index],
				Website:    "www." + nodeNameArr[index] + "_" + fmt.Sprint(i) + ".org",
				Details:    "This is " + nodeNameArr[index] + "_" + fmt.Sprint(i) + " Super Node",
			},
		}

		canAddr, _ := xutil.NodeId2Addr(canTmp.NodeId)

		stakingDB.SetCanPowerStore(hash, canAddr, canTmp)
		stakingDB.SetCandidateStore(hash, canAddr, canTmp)

		v := &staking.Validator{
			NodeAddress: canAddr,
			NodeId:      canTmp.NodeId,
			BlsPubKey:   canTmp.BlsPubKey,
			StakingWeight: [staking.SWeightItem]string{fmt.Sprint(xutil.CalcVersion(initProgramVersion)), canTmp.Shares.String(),
				fmt.Sprint(canTmp.StakingBlockNum), fmt.Sprint(canTmp.StakingTxIndex)},
			ValidatorTerm: 0,
		}
		validatorArr = append(validatorArr, v)
	}

	queue := validatorArr[:25]

	epoch_Arr := &staking.Validator_array{
		//Start: ((block-1)/22000)*22000 + 1,
		//End:   ((block-1)/22000)*22000 + 22000,
		Start: ((block-1)/uint64(xutil.CalcBlocksEachEpoch()))*uint64(xutil.CalcBlocksEachEpoch()) + 1,
		End:   ((block-1)/uint64(xutil.CalcBlocksEachEpoch()))*uint64(xutil.CalcBlocksEachEpoch()) + uint64(xutil.CalcBlocksEachEpoch()),
		Arr:   queue,
	}

	pre_Arr := &staking.Validator_array{
		Start: 0,
		End:   0,
		Arr:   queue,
	}

	curr_Arr := &staking.Validator_array{
		//Start: ((block-1)/250)*250 + 1,
		//End:   ((block-1)/250)*250 + 250,
		Start: ((block-1)/uint64(xutil.ConsensusSize()))*uint64(xutil.ConsensusSize()) + 1,
		End:   ((block-1)/uint64(xutil.ConsensusSize()))*uint64(xutil.ConsensusSize()) + uint64(xutil.ConsensusSize()),
		Arr:   queue,
	}

	setVerifierList(hash, epoch_Arr)
	setRoundValList(hash, pre_Arr)
	setRoundValList(hash, curr_Arr)

	lastBlockHash = hash
	lastBlockNumber = block
	lastHeader = header
}

func buildDbRestrictingPlan(t *testing.T, account common.Address, balance *big.Int, epochs int, stateDB xcom.StateDB) {
	//account := addrArr[0]
	fmt.Println("buildDbRestrictingPlan, store addr", account.Hex(), "balance", balance, "epochs", epochs)
	//const Epochs = 5
	var list = make([]uint64, 0)

	for epoch := 1; epoch <= epochs; epoch++ {
		// build release account record
		releaseAccountKey := restricting.GetReleaseAccountKey(uint64(epoch), 1)
		stateDB.SetState(cvm.RestrictingContractAddr, releaseAccountKey, account.Bytes())

		// build release amount record 1eth
		releaseAmount := balance
		releaseAmountKey := restricting.GetReleaseAmountKey(uint64(epoch), account)
		stateDB.SetState(cvm.RestrictingContractAddr, releaseAmountKey, releaseAmount.Bytes())

		// build release epoch list record
		releaseEpochKey := restricting.GetReleaseEpochKey(uint64(epoch))
		stateDB.SetState(cvm.RestrictingContractAddr, releaseEpochKey, common.Uint32ToBytes(1))

		list = append(list, uint64(epoch))
	}

	lockAmount := new(big.Int).Mul(balance, big.NewInt(int64(epochs)))

	fmt.Println("")

	// build restricting user info
	var user restricting.RestrictingInfo
	user.ReleaseList = list
	user.CachePlanAmount = lockAmount
	user.StakingAmount = big.NewInt(0)
	user.NeedRelease = big.NewInt(0)

	bUser, err := rlp.EncodeToBytes(user)
	if err != nil {
		t.Fatalf("failed to rlp encode restricting info: %s", err.Error())
	}

	// build restricting account info record
	restrictingKey := restricting.GetRestrictingKey(account)
	stateDB.SetState(cvm.RestrictingContractAddr, restrictingKey, bUser)

	//stateDB.AddBalance(sender, sender_balance)

	stateDB.AddBalance(cvm.RestrictingContractAddr, lockAmount)
}

func setRoundValList(blockHash common.Hash, val_Arr *staking.Validator_array) error {

	stakeDB := staking.NewStakingDB()

	queue, err := stakeDB.GetRoundValIndexByBlockHash(blockHash)
	if nil != err && err != snapshotdb.ErrNotFound {
		log.Error("Failed to setRoundValList: Query round valIndex is failed", "blockHash",
			blockHash.Hex(), "Start", val_Arr.Start, "End", val_Arr.End, "err", err)
		return err
	}

	var indexQueue staking.ValArrIndexQueue

	index := &staking.ValArrIndex{
		Start: val_Arr.Start,
		End:   val_Arr.End,
	}

	if len(queue) == 0 {
		indexQueue = make(staking.ValArrIndexQueue, 0)
		_, indexQueue = indexQueue.ConstantAppend(index, plugin.RoundValIndexSize)
	} else {

		has := false
		for _, indexInfo := range queue {
			if indexInfo.Start == val_Arr.Start && indexInfo.End == val_Arr.End {
				has = true
				break
			}
		}
		indexQueue = queue
		if !has {

			shabby, queue := queue.ConstantAppend(index, plugin.RoundValIndexSize)
			indexQueue = queue
			// delete the shabby validators
			if nil != shabby {
				if err := stakeDB.DelRoundValListByBlockHash(blockHash, shabby.Start, shabby.End); nil != err {
					log.Error("Failed to setRoundValList: delete shabby validators is failed",
						"shabby start", shabby.Start, "shabby end", shabby.End, "blockHash", blockHash.Hex())
					return err
				}
			}
		}
	}

	// Store new index Arr
	if err := stakeDB.SetRoundValIndex(blockHash, indexQueue); nil != err {
		log.Error("Failed to setRoundValList: store round validators new indexArr is failed", "blockHash", blockHash.Hex())
		return err
	}

	// Store new round validator Item
	if err := stakeDB.SetRoundValList(blockHash, index.Start, index.End, val_Arr.Arr); nil != err {
		log.Error("Failed to setRoundValList: store new round validators is failed", "blockHash", blockHash.Hex())
		return err
	}

	return nil
}

func setVerifierList(blockHash common.Hash, val_Arr *staking.Validator_array) error {

	stakeDB := staking.NewStakingDB()

	queue, err := stakeDB.GetEpochValIndexByBlockHash(blockHash)
	if nil != err && err != snapshotdb.ErrNotFound {
		log.Error("Failed to setVerifierList: Query epoch valIndex is failed", "blockHash",
			blockHash.Hex(), "Start", val_Arr.Start, "End", val_Arr.End, "err", err)
		return err
	}

	var indexQueue staking.ValArrIndexQueue

	index := &staking.ValArrIndex{
		Start: val_Arr.Start,
		End:   val_Arr.End,
	}

	if len(queue) == 0 {
		indexQueue = make(staking.ValArrIndexQueue, 0)
		_, indexQueue = indexQueue.ConstantAppend(index, plugin.EpochValIndexSize)
	} else {

		has := false
		for _, indexInfo := range queue {
			if indexInfo.Start == val_Arr.Start && indexInfo.End == val_Arr.End {
				has = true
				break
			}
		}
		indexQueue = queue
		if !has {

			shabby, queue := queue.ConstantAppend(index, plugin.EpochValIndexSize)
			indexQueue = queue
			// delete the shabby validators
			if nil != shabby {
				if err := stakeDB.DelEpochValListByBlockHash(blockHash, shabby.Start, shabby.End); nil != err {
					log.Error("Failed to setVerifierList: delete shabby validators is failed",
						"shabby start", shabby.Start, "shabby end", shabby.End, "blockHash", blockHash.Hex())
					return err
				}
			}
		}
	}

	// Store new index Arr
	if err := stakeDB.SetEpochValIndex(blockHash, indexQueue); nil != err {
		log.Error("Failed to setVerifierList: store epoch validators new indexArr is failed", "blockHash", blockHash.Hex())
		return err
	}

	// Store new epoch validator Item
	if err := stakeDB.SetEpochValList(blockHash, index.Start, index.End, val_Arr.Arr); nil != err {
		log.Error("Failed to setVerifierList: store new epoch validators is failed", "blockHash", blockHash.Hex())
		return err
	}

	return nil
}

func buildBlockNoCommit(blockNum int) {

	no := int64(blockNum)
	header := types.Header{
		Number: big.NewInt(no),
	}
	hash := header.Hash()

	staking.NewStakingDB()
	sndb.NewBlock(big.NewInt(int64(blockNum)), lastBlockHash, hash)

	lastBlockHash = hash
	lastBlockNumber = uint64(blockNum)
	lastHeader = header
}
