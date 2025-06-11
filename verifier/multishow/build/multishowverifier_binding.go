// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package verifier

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// VerifierMetaData contains all meta data concerning the Verifier contract.
var VerifierMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_bloom\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_zkpVerifier\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_x\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_y\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"bloom\",\"outputs\":[{\"internalType\":\"contractCascadingBloomFilter\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[8]\",\"name\":\"proof\",\"type\":\"uint256[8]\"},{\"internalType\":\"uint256\",\"name\":\"token\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"}],\"name\":\"checkCredential\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"valid\",\"type\":\"bool\"},{\"internalType\":\"uint8\",\"name\":\"errorCode\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"issuer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"issuerPubKeyX\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"issuerPubKeyY\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[8]\",\"name\":\"proof\",\"type\":\"uint256[8]\"},{\"internalType\":\"uint256\",\"name\":\"token\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"}],\"name\":\"measureCheckCredentialGas\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"valid\",\"type\":\"bool\"},{\"internalType\":\"uint8\",\"name\":\"errorCode\",\"type\":\"uint8\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"newFilters\",\"type\":\"bytes[]\"},{\"internalType\":\"uint256[]\",\"name\":\"ks\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"bitLens\",\"type\":\"uint256[]\"}],\"name\":\"update\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"verifier\",\"outputs\":[{\"internalType\":\"contractVerifier\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x6080604052346100335761001d61001461014d565b929190916102d0565b610025610038565b610f4761031a8239610f4790f35b61003e565b60405190565b600080fd5b601f801991011690565b634e487b7160e01b600052604160045260246000fd5b9061006d90610043565b810190811060018060401b0382111761008557604052565b61004d565b9061009d610096610038565b9283610063565b565b600080fd5b60018060a01b031690565b6100b8906100a4565b90565b6100c4816100af565b036100cb57565b600080fd5b905051906100dd826100bb565b565b90565b6100eb816100df565b036100f257565b600080fd5b90505190610104826100e2565b565b6080818303126101485761011d82600083016100d0565b9261014561012e84602085016100d0565b9361013c81604086016100f7565b936060016100f7565b90565b61009f565b61016b611261803803806101608161008a565b928339810190610106565b90919293565b60001b90565b9061018860018060a01b0391610171565b9181191691161790565b90565b6101a96101a46101ae926100a4565b610192565b6100a4565b90565b6101ba90610195565b90565b6101c6906101b1565b90565b90565b906101e16101dc6101e8926101bd565b6101c9565b8254610177565b9055565b6101f590610195565b90565b610201906101ec565b90565b61020d906101ec565b90565b90565b9061022861022361022f92610204565b610210565b8254610177565b9055565b61023c90610195565b90565b61024890610233565b90565b61025490610233565b90565b90565b9061026f61026a6102769261024b565b610257565b8254610177565b9055565b9061028760001991610171565b9181191691161790565b6102a56102a06102aa926100df565b610192565b6100df565b90565b90565b906102c56102c06102cc92610291565b6102ad565b825461027a565b9055565b91610309610302610310936102fd6102f661031798976102f13360026101cc565b6101f8565b6000610213565b61023f565b600161025a565b60036102b0565b60046102b0565b56fe60806040526004361015610013575b61064a565b61001e60003561009d565b80631d143848146100985780632b7ac3f314610093578063351c47f61461008e5780634d757b5914610089578063d26da14f14610084578063e099677d1461007f578063e9b2cd3f1461007a5763ffde64fc0361000e57610615565b6105cb565b61047f565b61044a565b6103a5565b6102a6565b610212565b610142565b60e01c90565b60405190565b600080fd5b600080fd5b60009103126100be57565b6100ae565b1c90565b60018060a01b031690565b6100e29060086100e793026100c3565b6100c7565b90565b906100f591546100d2565b90565b61010560026000906100ea565b90565b60018060a01b031690565b61011c90610108565b90565b61012890610113565b9052565b91906101409060006020850194019061011f565b565b34610172576101523660046100b3565b61016e61015d6100f8565b6101656100a3565b9182918261012c565b0390f35b6100a9565b60018060a01b031690565b61019290600861019793026100c3565b610177565b90565b906101a59154610182565b90565b6101b5600160009061019a565b90565b90565b6101cf6101ca6101d492610108565b6101b8565b610108565b90565b6101e0906101bb565b90565b6101ec906101d7565b90565b6101f8906101e3565b9052565b9190610210906000602085019401906101ef565b565b34610242576102223660046100b3565b61023e61022d6101a8565b6102356100a3565b918291826101fc565b0390f35b6100a9565b90565b61025a90600861025f93026100c3565b610247565b90565b9061026d915461024a565b90565b61027d6003600090610262565b90565b90565b61028c90610280565b9052565b91906102a490600060208501940190610283565b565b346102d6576102b63660046100b3565b6102d26102c1610270565b6102c96100a3565b91829182610290565b0390f35b6100a9565b600080fd5b600080fd5b919060206008028301116102f557565b6102e0565b61030381610280565b0361030a57565b600080fd5b9050359061031c826102fa565b565b9091610140828403126103575761035461033b84600085016102e5565b9361034a81610100860161030f565b936101200161030f565b90565b6100ae565b151590565b61036a9061035c565b9052565b60ff1690565b61037d9061036e565b9052565b9160206103a392949361039c60408201966000830190610361565b0190610374565b565b346103d7576103be6103b836600461031e565b916109c6565b906103d36103ca6100a3565b92839283610381565b0390f35b6100a9565b60018060a01b031690565b6103f79060086103fc93026100c3565b6103dc565b90565b9061040a91546103e7565b90565b6104186000806103ff565b90565b610424906101d7565b90565b6104309061041b565b9052565b919061044890600060208501940190610427565b565b3461047a5761045a3660046100b3565b61047661046561040d565b61046d6100a3565b91829182610434565b0390f35b6100a9565b346104b15761049861049236600461031e565b91610bbd565b906104ad6104a46100a3565b92839283610381565b0390f35b6100a9565b600080fd5b600080fd5b909182601f830112156104fa5781359167ffffffffffffffff83116104f55760200192602083028401116104f057565b6102e0565b6104bb565b6104b6565b909182601f830112156105395781359167ffffffffffffffff831161053457602001926020830284011161052f57565b6102e0565b6104bb565b6104b6565b906060828203126105c057600082013567ffffffffffffffff81116105bb57816105699184016104c0565b929093602082013567ffffffffffffffff81116105b6578361058c9184016104ff565b929093604082013567ffffffffffffffff81116105b1576105ad92016104ff565b9091565b6102db565b6102db565b6102db565b6100ae565b60000190565b34610600576105ea6105de36600461053e565b94939093929192610f01565b6105f26100a3565b806105fc816105c5565b0390f35b6100a9565b6106126004600090610262565b90565b34610645576106253660046100b3565b610641610630610605565b6106386100a3565b91829182610290565b0390f35b6100a9565b600080fd5b600090565b600090565b601f801991011690565b634e487b7160e01b600052604160045260246000fd5b9061068390610659565b810190811067ffffffffffffffff82111761069d57604052565b610663565b906106b56106ae6100a3565b9283610679565b565b67ffffffffffffffff81116106cc5760200290565b610663565b6106dd6106e2916106b7565b6106a2565b90565b60001c90565b6106f76106fc916106e5565b610247565b90565b61070990546106eb565b90565b9061071690610280565b9052565b61072661072b916106e5565b610177565b90565b610738905461071a565b90565b600080fd5b60e01b90565b600091031261075157565b6100ae565b9037565b6107679161010091610756565b565b50600490565b905090565b90565b61078090610280565b9052565b9061079181602093610777565b0190565b60200190565b6107b76107b16107aa83610769565b809461076f565b91610774565b6000915b8383106107c85750505050565b6107de6107d86001928451610784565b92610795565b920191906107bb565b9161010061080b929493610804610180820196600083019061075a565b019061079b565b565b6108156100a3565b3d6000823e3d90fd5b90565b61083561083061083a9261081e565b6101b8565b61036e565b90565b61084961084e916106e5565b6103dc565b90565b61085b905461083d565b90565b90565b60001b90565b61087b61087661088092610280565b610861565b61085e565b90565b90565b6108926108979161085e565b610883565b9052565b6108a781602093610886565b0190565b6108b48161035c565b036108bb57565b600080fd5b905051906108cd826108ab565b565b905051906108dc826102fa565b565b919060408382031261090757806108fb61090492600086016108c0565b936020016108cf565b90565b6100ae565b5190565b60209181520190565b60005b83811061092d575050906000910152565b80602091830151818501520161091c565b61095d61096660209361096b936109548161090c565b93848093610910565b95869101610919565b610659565b0190565b610985916020820191600081840391015261093e565b90565b90565b61099f61099a6109a492610988565b6101b8565b61036e565b90565b90565b6109be6109b96109c3926109a7565b6101b8565b61036e565b90565b90916109d061064f565b506109d9610654565b50610a2a6109e760046106d1565b916109fe6109f560036106ff565b6000850161070c565b610a14610a0b60046106ff565b6020850161070c565b610a21856040850161070c565b6060830161070c565b90610a3d610a38600161072e565b6101e3565b916323572511919092803b15610bb857610a6a600093610a75610a5e6100a3565b96879586948594610740565b8452600484016107e7565b03915afa9081610b8b575b5015600014610b7d576001610b6c576040610ad4610b02925b610af7610aae610aa96000610851565b61041b565b91610ae3610ac063d423db2a92610867565b610ac86100a3565b9586916020830161089b565b60208201810382520385610679565b610aeb6100a3565b95869485938493610740565b83526004830161096f565b03915afa908115610b6757600091610b3a575b50610b2a57600190610b2760006109aa565b90565b600090610b37600261098b565b90565b610b5b915060403d8111610b60575b610b538183610679565b8101906108de565b610b15565b503d610b49565b61080d565b50600090610b7a6001610821565b90565b6040610ad4610b0292610a99565b610bab9060003d8111610bb1575b610ba38183610679565b810190610746565b38610a80565b503d610b99565b61073b565b91610bdc92610bca61064f565b50610bd3610654565b509190916109c6565b91909190565b610bee610bf3916106e5565b6100c7565b90565b610c009054610be2565b90565b60209181520190565b60007f4e6f742069737375657200000000000000000000000000000000000000000000910152565b610c41600a602092610c03565b610c4a81610c0c565b0190565b610c649060208101906000818303910152610c34565b90565b15610c6e57565b610c766100a3565b62461bcd60e51b815280610c8c60048201610c4e565b0390fd5b90610cc39594939291610cbe33610cb8610cb2610cad6002610bf6565b610113565b91610113565b14610c67565b610e61565b565b60209181520190565b90565b60209181520190565b90826000939282370152565b9190610d0081610cf981610d0595610cd1565b8095610cda565b610659565b0190565b90610d149291610ce6565b90565b600080fd5b600080fd5b600080fd5b9035600160200382360303811215610d6757016020813591019167ffffffffffffffff8211610d62576001820236038313610d5d57565b610d1c565b610d17565b610d21565b60200190565b9181610d7d91610cc5565b9081610d8e60208302840194610cce565b92836000925b848410610da45750505050505090565b9091929394956020610dd0610dca8385600195038852610dc48b88610d26565b90610d09565b98610d6c565b940194019294939190610d94565b60209181520190565b600080fd5b909182610df891610dde565b9160018060fb1b038111610e1b5782916020610e179202938491610756565b0190565b610de7565b94929093610e42610e5e9795610e5094606089019189830360008b0152610d72565b918683036020880152610dec565b926040818503910152610dec565b90565b9194909293610e78610e736000610851565b61041b565b9263b163337d90949695919295843b15610efc57600096610ead948894610eb893610ea16100a3565b9b8c9a8b998a98610740565b885260048801610e20565b03925af18015610ef757610eca575b50565b610eea9060003d8111610ef0575b610ee28183610679565b810190610746565b38610ec7565b503d610ed8565b61080d565b61073b565b90610f0f9594939291610c90565b56fea26469706673582212203ce38d7d19987d0ae2c19b18e8127ebe3184e904feb24a1680ef2cee67fa83d264736f6c634300081e0033",
}

// VerifierABI is the input ABI used to generate the binding from.
// Deprecated: Use VerifierMetaData.ABI instead.
var VerifierABI = VerifierMetaData.ABI

// VerifierBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use VerifierMetaData.Bin instead.
var VerifierBin = VerifierMetaData.Bin

// DeployVerifier deploys a new Ethereum contract, binding an instance of Verifier to it.
func DeployVerifier(auth *bind.TransactOpts, backend bind.ContractBackend, _bloom common.Address, _zkpVerifier common.Address, _x *big.Int, _y *big.Int) (common.Address, *types.Transaction, *Verifier, error) {
	parsed, err := VerifierMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VerifierBin), backend, _bloom, _zkpVerifier, _x, _y)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Verifier{VerifierCaller: VerifierCaller{contract: contract}, VerifierTransactor: VerifierTransactor{contract: contract}, VerifierFilterer: VerifierFilterer{contract: contract}}, nil
}

// Verifier is an auto generated Go binding around an Ethereum contract.
type Verifier struct {
	VerifierCaller     // Read-only binding to the contract
	VerifierTransactor // Write-only binding to the contract
	VerifierFilterer   // Log filterer for contract events
}

// VerifierCaller is an auto generated read-only Go binding around an Ethereum contract.
type VerifierCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VerifierTransactor is an auto generated write-only Go binding around an Ethereum contract.
type VerifierTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VerifierFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type VerifierFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VerifierSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type VerifierSession struct {
	Contract     *Verifier         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// VerifierCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type VerifierCallerSession struct {
	Contract *VerifierCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// VerifierTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type VerifierTransactorSession struct {
	Contract     *VerifierTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// VerifierRaw is an auto generated low-level Go binding around an Ethereum contract.
type VerifierRaw struct {
	Contract *Verifier // Generic contract binding to access the raw methods on
}

// VerifierCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type VerifierCallerRaw struct {
	Contract *VerifierCaller // Generic read-only contract binding to access the raw methods on
}

// VerifierTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type VerifierTransactorRaw struct {
	Contract *VerifierTransactor // Generic write-only contract binding to access the raw methods on
}

// NewVerifier creates a new instance of Verifier, bound to a specific deployed contract.
func NewVerifier(address common.Address, backend bind.ContractBackend) (*Verifier, error) {
	contract, err := bindVerifier(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Verifier{VerifierCaller: VerifierCaller{contract: contract}, VerifierTransactor: VerifierTransactor{contract: contract}, VerifierFilterer: VerifierFilterer{contract: contract}}, nil
}

// NewVerifierCaller creates a new read-only instance of Verifier, bound to a specific deployed contract.
func NewVerifierCaller(address common.Address, caller bind.ContractCaller) (*VerifierCaller, error) {
	contract, err := bindVerifier(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VerifierCaller{contract: contract}, nil
}

// NewVerifierTransactor creates a new write-only instance of Verifier, bound to a specific deployed contract.
func NewVerifierTransactor(address common.Address, transactor bind.ContractTransactor) (*VerifierTransactor, error) {
	contract, err := bindVerifier(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VerifierTransactor{contract: contract}, nil
}

// NewVerifierFilterer creates a new log filterer instance of Verifier, bound to a specific deployed contract.
func NewVerifierFilterer(address common.Address, filterer bind.ContractFilterer) (*VerifierFilterer, error) {
	contract, err := bindVerifier(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VerifierFilterer{contract: contract}, nil
}

// bindVerifier binds a generic wrapper to an already deployed contract.
func bindVerifier(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VerifierMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Verifier *VerifierRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Verifier.Contract.VerifierCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Verifier *VerifierRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Verifier.Contract.VerifierTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Verifier *VerifierRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Verifier.Contract.VerifierTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Verifier *VerifierCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Verifier.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Verifier *VerifierTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Verifier.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Verifier *VerifierTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Verifier.Contract.contract.Transact(opts, method, params...)
}

// Bloom is a free data retrieval call binding the contract method 0xd26da14f.
//
// Solidity: function bloom() view returns(address)
func (_Verifier *VerifierCaller) Bloom(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Verifier.contract.Call(opts, &out, "bloom")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Bloom is a free data retrieval call binding the contract method 0xd26da14f.
//
// Solidity: function bloom() view returns(address)
func (_Verifier *VerifierSession) Bloom() (common.Address, error) {
	return _Verifier.Contract.Bloom(&_Verifier.CallOpts)
}

// Bloom is a free data retrieval call binding the contract method 0xd26da14f.
//
// Solidity: function bloom() view returns(address)
func (_Verifier *VerifierCallerSession) Bloom() (common.Address, error) {
	return _Verifier.Contract.Bloom(&_Verifier.CallOpts)
}

// CheckCredential is a free data retrieval call binding the contract method 0x4d757b59.
//
// Solidity: function checkCredential(uint256[8] proof, uint256 token, uint256 epoch) view returns(bool valid, uint8 errorCode)
func (_Verifier *VerifierCaller) CheckCredential(opts *bind.CallOpts, proof [8]*big.Int, token *big.Int, epoch *big.Int) (struct {
	Valid     bool
	ErrorCode uint8
}, error) {
	var out []interface{}
	err := _Verifier.contract.Call(opts, &out, "checkCredential", proof, token, epoch)

	outstruct := new(struct {
		Valid     bool
		ErrorCode uint8
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Valid = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.ErrorCode = *abi.ConvertType(out[1], new(uint8)).(*uint8)

	return *outstruct, err

}

// CheckCredential is a free data retrieval call binding the contract method 0x4d757b59.
//
// Solidity: function checkCredential(uint256[8] proof, uint256 token, uint256 epoch) view returns(bool valid, uint8 errorCode)
func (_Verifier *VerifierSession) CheckCredential(proof [8]*big.Int, token *big.Int, epoch *big.Int) (struct {
	Valid     bool
	ErrorCode uint8
}, error) {
	return _Verifier.Contract.CheckCredential(&_Verifier.CallOpts, proof, token, epoch)
}

// CheckCredential is a free data retrieval call binding the contract method 0x4d757b59.
//
// Solidity: function checkCredential(uint256[8] proof, uint256 token, uint256 epoch) view returns(bool valid, uint8 errorCode)
func (_Verifier *VerifierCallerSession) CheckCredential(proof [8]*big.Int, token *big.Int, epoch *big.Int) (struct {
	Valid     bool
	ErrorCode uint8
}, error) {
	return _Verifier.Contract.CheckCredential(&_Verifier.CallOpts, proof, token, epoch)
}

// Issuer is a free data retrieval call binding the contract method 0x1d143848.
//
// Solidity: function issuer() view returns(address)
func (_Verifier *VerifierCaller) Issuer(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Verifier.contract.Call(opts, &out, "issuer")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Issuer is a free data retrieval call binding the contract method 0x1d143848.
//
// Solidity: function issuer() view returns(address)
func (_Verifier *VerifierSession) Issuer() (common.Address, error) {
	return _Verifier.Contract.Issuer(&_Verifier.CallOpts)
}

// Issuer is a free data retrieval call binding the contract method 0x1d143848.
//
// Solidity: function issuer() view returns(address)
func (_Verifier *VerifierCallerSession) Issuer() (common.Address, error) {
	return _Verifier.Contract.Issuer(&_Verifier.CallOpts)
}

// IssuerPubKeyX is a free data retrieval call binding the contract method 0x351c47f6.
//
// Solidity: function issuerPubKeyX() view returns(uint256)
func (_Verifier *VerifierCaller) IssuerPubKeyX(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Verifier.contract.Call(opts, &out, "issuerPubKeyX")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// IssuerPubKeyX is a free data retrieval call binding the contract method 0x351c47f6.
//
// Solidity: function issuerPubKeyX() view returns(uint256)
func (_Verifier *VerifierSession) IssuerPubKeyX() (*big.Int, error) {
	return _Verifier.Contract.IssuerPubKeyX(&_Verifier.CallOpts)
}

// IssuerPubKeyX is a free data retrieval call binding the contract method 0x351c47f6.
//
// Solidity: function issuerPubKeyX() view returns(uint256)
func (_Verifier *VerifierCallerSession) IssuerPubKeyX() (*big.Int, error) {
	return _Verifier.Contract.IssuerPubKeyX(&_Verifier.CallOpts)
}

// IssuerPubKeyY is a free data retrieval call binding the contract method 0xffde64fc.
//
// Solidity: function issuerPubKeyY() view returns(uint256)
func (_Verifier *VerifierCaller) IssuerPubKeyY(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Verifier.contract.Call(opts, &out, "issuerPubKeyY")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// IssuerPubKeyY is a free data retrieval call binding the contract method 0xffde64fc.
//
// Solidity: function issuerPubKeyY() view returns(uint256)
func (_Verifier *VerifierSession) IssuerPubKeyY() (*big.Int, error) {
	return _Verifier.Contract.IssuerPubKeyY(&_Verifier.CallOpts)
}

// IssuerPubKeyY is a free data retrieval call binding the contract method 0xffde64fc.
//
// Solidity: function issuerPubKeyY() view returns(uint256)
func (_Verifier *VerifierCallerSession) IssuerPubKeyY() (*big.Int, error) {
	return _Verifier.Contract.IssuerPubKeyY(&_Verifier.CallOpts)
}

// Verifier is a free data retrieval call binding the contract method 0x2b7ac3f3.
//
// Solidity: function verifier() view returns(address)
func (_Verifier *VerifierCaller) Verifier(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Verifier.contract.Call(opts, &out, "verifier")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Verifier is a free data retrieval call binding the contract method 0x2b7ac3f3.
//
// Solidity: function verifier() view returns(address)
func (_Verifier *VerifierSession) Verifier() (common.Address, error) {
	return _Verifier.Contract.Verifier(&_Verifier.CallOpts)
}

// Verifier is a free data retrieval call binding the contract method 0x2b7ac3f3.
//
// Solidity: function verifier() view returns(address)
func (_Verifier *VerifierCallerSession) Verifier() (common.Address, error) {
	return _Verifier.Contract.Verifier(&_Verifier.CallOpts)
}

// MeasureCheckCredentialGas is a paid mutator transaction binding the contract method 0xe099677d.
//
// Solidity: function measureCheckCredentialGas(uint256[8] proof, uint256 token, uint256 epoch) returns(bool valid, uint8 errorCode)
func (_Verifier *VerifierTransactor) MeasureCheckCredentialGas(opts *bind.TransactOpts, proof [8]*big.Int, token *big.Int, epoch *big.Int) (*types.Transaction, error) {
	return _Verifier.contract.Transact(opts, "measureCheckCredentialGas", proof, token, epoch)
}

// MeasureCheckCredentialGas is a paid mutator transaction binding the contract method 0xe099677d.
//
// Solidity: function measureCheckCredentialGas(uint256[8] proof, uint256 token, uint256 epoch) returns(bool valid, uint8 errorCode)
func (_Verifier *VerifierSession) MeasureCheckCredentialGas(proof [8]*big.Int, token *big.Int, epoch *big.Int) (*types.Transaction, error) {
	return _Verifier.Contract.MeasureCheckCredentialGas(&_Verifier.TransactOpts, proof, token, epoch)
}

// MeasureCheckCredentialGas is a paid mutator transaction binding the contract method 0xe099677d.
//
// Solidity: function measureCheckCredentialGas(uint256[8] proof, uint256 token, uint256 epoch) returns(bool valid, uint8 errorCode)
func (_Verifier *VerifierTransactorSession) MeasureCheckCredentialGas(proof [8]*big.Int, token *big.Int, epoch *big.Int) (*types.Transaction, error) {
	return _Verifier.Contract.MeasureCheckCredentialGas(&_Verifier.TransactOpts, proof, token, epoch)
}

// Update is a paid mutator transaction binding the contract method 0xe9b2cd3f.
//
// Solidity: function update(bytes[] newFilters, uint256[] ks, uint256[] bitLens) returns()
func (_Verifier *VerifierTransactor) Update(opts *bind.TransactOpts, newFilters [][]byte, ks []*big.Int, bitLens []*big.Int) (*types.Transaction, error) {
	return _Verifier.contract.Transact(opts, "update", newFilters, ks, bitLens)
}

// Update is a paid mutator transaction binding the contract method 0xe9b2cd3f.
//
// Solidity: function update(bytes[] newFilters, uint256[] ks, uint256[] bitLens) returns()
func (_Verifier *VerifierSession) Update(newFilters [][]byte, ks []*big.Int, bitLens []*big.Int) (*types.Transaction, error) {
	return _Verifier.Contract.Update(&_Verifier.TransactOpts, newFilters, ks, bitLens)
}

// Update is a paid mutator transaction binding the contract method 0xe9b2cd3f.
//
// Solidity: function update(bytes[] newFilters, uint256[] ks, uint256[] bitLens) returns()
func (_Verifier *VerifierTransactorSession) Update(newFilters [][]byte, ks []*big.Int, bitLens []*big.Int) (*types.Transaction, error) {
	return _Verifier.Contract.Update(&_Verifier.TransactOpts, newFilters, ks, bitLens)
}
