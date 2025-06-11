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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_bloom\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_zkpVerifier\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_x\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_y\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"bloom\",\"outputs\":[{\"internalType\":\"contractCascadingBloomFilter\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[8]\",\"name\":\"proof\",\"type\":\"uint256[8]\"},{\"internalType\":\"uint256\",\"name\":\"token\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"}],\"name\":\"checkCredential\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"valid\",\"type\":\"bool\"},{\"internalType\":\"uint8\",\"name\":\"errorCode\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"issuer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"issuerPubKeyX\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"issuerPubKeyY\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"newFilters\",\"type\":\"bytes[]\"},{\"internalType\":\"uint256[]\",\"name\":\"ks\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"bitLens\",\"type\":\"uint256[]\"}],\"name\":\"update\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"verifier\",\"outputs\":[{\"internalType\":\"contractVerifier\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x6080604052346100335761001d61001461014d565b929190916102d0565b610025610038565b610edb61031a8239610edb90f35b61003e565b60405190565b600080fd5b601f801991011690565b634e487b7160e01b600052604160045260246000fd5b9061006d90610043565b810190811060018060401b0382111761008557604052565b61004d565b9061009d610096610038565b9283610063565b565b600080fd5b60018060a01b031690565b6100b8906100a4565b90565b6100c4816100af565b036100cb57565b600080fd5b905051906100dd826100bb565b565b90565b6100eb816100df565b036100f257565b600080fd5b90505190610104826100e2565b565b6080818303126101485761011d82600083016100d0565b9261014561012e84602085016100d0565b9361013c81604086016100f7565b936060016100f7565b90565b61009f565b61016b6111f5803803806101608161008a565b928339810190610106565b90919293565b60001b90565b9061018860018060a01b0391610171565b9181191691161790565b90565b6101a96101a46101ae926100a4565b610192565b6100a4565b90565b6101ba90610195565b90565b6101c6906101b1565b90565b90565b906101e16101dc6101e8926101bd565b6101c9565b8254610177565b9055565b6101f590610195565b90565b610201906101ec565b90565b61020d906101ec565b90565b90565b9061022861022361022f92610204565b610210565b8254610177565b9055565b61023c90610195565b90565b61024890610233565b90565b61025490610233565b90565b90565b9061026f61026a6102769261024b565b610257565b8254610177565b9055565b9061028760001991610171565b9181191691161790565b6102a56102a06102aa926100df565b610192565b6100df565b90565b90565b906102c56102c06102cc92610291565b6102ad565b825461027a565b9055565b91610309610302610310936102fd6102f661031798976102f13360026101cc565b6101f8565b6000610213565b61023f565b600161025a565b60036102b0565b60046102b0565b56fe60806040526004361015610013575b610603565b61001e60003561008d565b80631d143848146100885780632b7ac3f314610083578063351c47f61461007e5780634d757b5914610079578063d26da14f14610074578063e9b2cd3f1461006f5763ffde64fc0361000e576105ce565b610584565b61043a565b610395565b610296565b610202565b610132565b60e01c90565b60405190565b600080fd5b600080fd5b60009103126100ae57565b61009e565b1c90565b60018060a01b031690565b6100d29060086100d793026100b3565b6100b7565b90565b906100e591546100c2565b90565b6100f560026000906100da565b90565b60018060a01b031690565b61010c906100f8565b90565b61011890610103565b9052565b91906101309060006020850194019061010f565b565b34610162576101423660046100a3565b61015e61014d6100e8565b610155610093565b9182918261011c565b0390f35b610099565b60018060a01b031690565b61018290600861018793026100b3565b610167565b90565b906101959154610172565b90565b6101a5600160009061018a565b90565b90565b6101bf6101ba6101c4926100f8565b6101a8565b6100f8565b90565b6101d0906101ab565b90565b6101dc906101c7565b90565b6101e8906101d3565b9052565b9190610200906000602085019401906101df565b565b34610232576102123660046100a3565b61022e61021d610198565b610225610093565b918291826101ec565b0390f35b610099565b90565b61024a90600861024f93026100b3565b610237565b90565b9061025d915461023a565b90565b61026d6003600090610252565b90565b90565b61027c90610270565b9052565b919061029490600060208501940190610273565b565b346102c6576102a63660046100a3565b6102c26102b1610260565b6102b9610093565b91829182610280565b0390f35b610099565b600080fd5b600080fd5b919060206008028301116102e557565b6102d0565b6102f381610270565b036102fa57565b600080fd5b9050359061030c826102ea565b565b9091610140828403126103475761034461032b84600085016102d5565b9361033a8161010086016102ff565b93610120016102ff565b90565b61009e565b151590565b61035a9061034c565b9052565b60ff1690565b61036d9061035e565b9052565b91602061039392949361038c60408201966000830190610351565b0190610364565b565b346103c7576103ae6103a836600461030e565b9161097f565b906103c36103ba610093565b92839283610371565b0390f35b610099565b60018060a01b031690565b6103e79060086103ec93026100b3565b6103cc565b90565b906103fa91546103d7565b90565b6104086000806103ef565b90565b610414906101c7565b90565b6104209061040b565b9052565b919061043890600060208501940190610417565b565b3461046a5761044a3660046100a3565b6104666104556103fd565b61045d610093565b91829182610424565b0390f35b610099565b600080fd5b600080fd5b909182601f830112156104b35781359167ffffffffffffffff83116104ae5760200192602083028401116104a957565b6102d0565b610474565b61046f565b909182601f830112156104f25781359167ffffffffffffffff83116104ed5760200192602083028401116104e857565b6102d0565b610474565b61046f565b9060608282031261057957600082013567ffffffffffffffff81116105745781610522918401610479565b929093602082013567ffffffffffffffff811161056f57836105459184016104b8565b929093604082013567ffffffffffffffff811161056a5761056692016104b8565b9091565b6102cb565b6102cb565b6102cb565b61009e565b60000190565b346105b9576105a36105973660046104f7565b94939093929192610e95565b6105ab610093565b806105b58161057e565b0390f35b610099565b6105cb6004600090610252565b90565b346105fe576105de3660046100a3565b6105fa6105e96105be565b6105f1610093565b91829182610280565b0390f35b610099565b600080fd5b600090565b600090565b601f801991011690565b634e487b7160e01b600052604160045260246000fd5b9061063c90610612565b810190811067ffffffffffffffff82111761065657604052565b61061c565b9061066e610667610093565b9283610632565b565b67ffffffffffffffff81116106855760200290565b61061c565b61069661069b91610670565b61065b565b90565b60001c90565b6106b06106b59161069e565b610237565b90565b6106c290546106a4565b90565b906106cf90610270565b9052565b6106df6106e49161069e565b610167565b90565b6106f190546106d3565b90565b600080fd5b60e01b90565b600091031261070a57565b61009e565b9037565b610720916101009161070f565b565b50600490565b905090565b90565b61073990610270565b9052565b9061074a81602093610730565b0190565b60200190565b61077061076a61076383610722565b8094610728565b9161072d565b6000915b8383106107815750505050565b610797610791600192845161073d565b9261074e565b92019190610774565b916101006107c49294936107bd6101808201966000830190610713565b0190610754565b565b6107ce610093565b3d6000823e3d90fd5b90565b6107ee6107e96107f3926107d7565b6101a8565b61035e565b90565b6108026108079161069e565b6103cc565b90565b61081490546107f6565b90565b90565b60001b90565b61083461082f61083992610270565b61081a565b610817565b90565b90565b61084b61085091610817565b61083c565b9052565b6108608160209361083f565b0190565b61086d8161034c565b0361087457565b600080fd5b9050519061088682610864565b565b90505190610895826102ea565b565b91906040838203126108c057806108b46108bd9260008601610879565b93602001610888565b90565b61009e565b5190565b60209181520190565b60005b8381106108e6575050906000910152565b8060209183015181850152016108d5565b61091661091f6020936109249361090d816108c5565b938480936108c9565b958691016108d2565b610612565b0190565b61093e91602082019160008184039101526108f7565b90565b90565b61095861095361095d92610941565b6101a8565b61035e565b90565b90565b61097761097261097c92610960565b6101a8565b61035e565b90565b9091610989610608565b5061099261060d565b506109e36109a0600461068a565b916109b76109ae60036106b8565b600085016106c5565b6109cd6109c460046106b8565b602085016106c5565b6109da85604085016106c5565b606083016106c5565b906109f66109f160016106e7565b6101d3565b916323572511919092803b15610b7157610a23600093610a2e610a17610093565b968795869485946106f9565b8452600484016107a0565b03915afa9081610b44575b5015600014610b36576001610b25576040610a8d610abb925b610ab0610a67610a62600061080a565b61040b565b91610a9c610a7963d423db2a92610820565b610a81610093565b95869160208301610854565b60208201810382520385610632565b610aa4610093565b958694859384936106f9565b835260048301610928565b03915afa908115610b2057600091610af3575b50610ae357600190610ae06000610963565b90565b600090610af06002610944565b90565b610b14915060403d8111610b19575b610b0c8183610632565b810190610897565b610ace565b503d610b02565b6107c6565b50600090610b3360016107da565b90565b6040610a8d610abb92610a52565b610b649060003d8111610b6a575b610b5c8183610632565b8101906106ff565b38610a39565b503d610b52565b6106f4565b610b82610b879161069e565b6100b7565b90565b610b949054610b76565b90565b60209181520190565b60007f4e6f742069737375657200000000000000000000000000000000000000000000910152565b610bd5600a602092610b97565b610bde81610ba0565b0190565b610bf89060208101906000818303910152610bc8565b90565b15610c0257565b610c0a610093565b62461bcd60e51b815280610c2060048201610be2565b0390fd5b90610c579594939291610c5233610c4c610c46610c416002610b8a565b610103565b91610103565b14610bfb565b610df5565b565b60209181520190565b90565b60209181520190565b90826000939282370152565b9190610c9481610c8d81610c9995610c65565b8095610c6e565b610612565b0190565b90610ca89291610c7a565b90565b600080fd5b600080fd5b600080fd5b9035600160200382360303811215610cfb57016020813591019167ffffffffffffffff8211610cf6576001820236038313610cf157565b610cb0565b610cab565b610cb5565b60200190565b9181610d1191610c59565b9081610d2260208302840194610c62565b92836000925b848410610d385750505050505090565b9091929394956020610d64610d5e8385600195038852610d588b88610cba565b90610c9d565b98610d00565b940194019294939190610d28565b60209181520190565b600080fd5b909182610d8c91610d72565b9160018060fb1b038111610daf5782916020610dab920293849161070f565b0190565b610d7b565b94929093610dd6610df29795610de494606089019189830360008b0152610d06565b918683036020880152610d80565b926040818503910152610d80565b90565b9194909293610e0c610e07600061080a565b61040b565b9263b163337d90949695919295843b15610e9057600096610e41948894610e4c93610e35610093565b9b8c9a8b998a986106f9565b885260048801610db4565b03925af18015610e8b57610e5e575b50565b610e7e9060003d8111610e84575b610e768183610632565b8101906106ff565b38610e5b565b503d610e6c565b6107c6565b6106f4565b90610ea39594939291610c24565b56fea26469706673582212206afbb1481a5a51d40db0d410762fbc65b2c1e7a4fb70d7fa7779747f55eaafa264736f6c634300081e0033",
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
