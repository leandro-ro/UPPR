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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_bloom\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_zkpVerifier\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_x\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_y\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"bloom\",\"outputs\":[{\"internalType\":\"contractCascadingBloomFilter\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[8]\",\"name\":\"proof\",\"type\":\"uint256[8]\"},{\"internalType\":\"uint256\",\"name\":\"pubKeyX\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"pubKeyY\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"token\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"}],\"name\":\"checkCredential\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"valid\",\"type\":\"bool\"},{\"internalType\":\"uint8\",\"name\":\"errorCode\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"pubKeyXx\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"pubKeyYy\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"usedToken\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"usedEpoch\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"issuer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"issuerPubKeyX\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"issuerPubKeyY\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"newFilters\",\"type\":\"bytes[]\"},{\"internalType\":\"uint256[]\",\"name\":\"ks\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"bitLens\",\"type\":\"uint256[]\"}],\"name\":\"update\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"verifier\",\"outputs\":[{\"internalType\":\"contractVerifier\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x6080604052346100335761001d61001461014d565b929190916102d0565b610025610038565b610d7061031a8239610d7090f35b61003e565b60405190565b600080fd5b601f801991011690565b634e487b7160e01b600052604160045260246000fd5b9061006d90610043565b810190811060018060401b0382111761008557604052565b61004d565b9061009d610096610038565b9283610063565b565b600080fd5b60018060a01b031690565b6100b8906100a4565b90565b6100c4816100af565b036100cb57565b600080fd5b905051906100dd826100bb565b565b90565b6100eb816100df565b036100f257565b600080fd5b90505190610104826100e2565b565b6080818303126101485761011d82600083016100d0565b9261014561012e84602085016100d0565b9361013c81604086016100f7565b936060016100f7565b90565b61009f565b61016b61108a803803806101608161008a565b928339810190610106565b90919293565b60001b90565b9061018860018060a01b0391610171565b9181191691161790565b90565b6101a96101a46101ae926100a4565b610192565b6100a4565b90565b6101ba90610195565b90565b6101c6906101b1565b90565b90565b906101e16101dc6101e8926101bd565b6101c9565b8254610177565b9055565b6101f590610195565b90565b610201906101ec565b90565b61020d906101ec565b90565b90565b9061022861022361022f92610204565b610210565b8254610177565b9055565b61023c90610195565b90565b61024890610233565b90565b61025490610233565b90565b90565b9061026f61026a6102769261024b565b610257565b8254610177565b9055565b9061028760001991610171565b9181191691161790565b6102a56102a06102aa926100df565b610192565b6100df565b90565b90565b906102c56102c06102cc92610291565b6102ad565b825461027a565b9055565b91610309610302610310936102fd6102f661031798976102f13360026101cc565b6101f8565b6000610213565b61023f565b600161025a565b60036102b0565b60046102b0565b56fe60806040526004361015610013575b610662565b61001e60003561008d565b80631d143848146100885780632b7ac3f314610083578063351c47f61461007e5780639d2cd56914610079578063d26da14f14610074578063e9b2cd3f1461006f5763ffde64fc0361000e5761062d565b6105e3565b610499565b6103eb565b610296565b610202565b610132565b60e01c90565b60405190565b600080fd5b600080fd5b60009103126100ae57565b61009e565b1c90565b60018060a01b031690565b6100d29060086100d793026100b3565b6100b7565b90565b906100e591546100c2565b90565b6100f560026000906100da565b90565b60018060a01b031690565b61010c906100f8565b90565b61011890610103565b9052565b91906101309060006020850194019061010f565b565b34610162576101423660046100a3565b61015e61014d6100e8565b610155610093565b9182918261011c565b0390f35b610099565b60018060a01b031690565b61018290600861018793026100b3565b610167565b90565b906101959154610172565b90565b6101a5600160009061018a565b90565b90565b6101bf6101ba6101c4926100f8565b6101a8565b6100f8565b90565b6101d0906101ab565b90565b6101dc906101c7565b90565b6101e8906101d3565b9052565b9190610200906000602085019401906101df565b565b34610232576102123660046100a3565b61022e61021d610198565b610225610093565b918291826101ec565b0390f35b610099565b90565b61024a90600861024f93026100b3565b610237565b90565b9061025d915461023a565b90565b61026d6003600090610252565b90565b90565b61027c90610270565b9052565b919061029490600060208501940190610273565b565b346102c6576102a63660046100a3565b6102c26102b1610260565b6102b9610093565b91829182610280565b0390f35b610099565b600080fd5b600080fd5b919060206008028301116102e557565b6102d0565b6102f381610270565b036102fa57565b600080fd5b9050359061030c826102ea565b565b9190610180838203126103655761032881600085016102d5565b926103378261010083016102ff565b926103626103498461012085016102ff565b936103588161014086016102ff565b93610160016102ff565b90565b61009e565b151590565b6103789061036a565b9052565b60ff1690565b61038b9061037c565b9052565b91946103d86103e2929897956103ce60a0966103c46103e99a6103ba60c08a019e60008b019061036f565b6020890190610382565b6040870190610273565b6060850190610273565b6080830190610273565b0190610273565b565b346104265761042261040a61040136600461030e565b93929092610879565b92610419969496929192610093565b9687968761038f565b0390f35b610099565b60018060a01b031690565b61044690600861044b93026100b3565b61042b565b90565b906104599154610436565b90565b61046760008061044e565b90565b610473906101c7565b90565b61047f9061046a565b9052565b919061049790600060208501940190610476565b565b346104c9576104a93660046100a3565b6104c56104b461045c565b6104bc610093565b91829182610483565b0390f35b610099565b600080fd5b600080fd5b909182601f830112156105125781359167ffffffffffffffff831161050d57602001926020830284011161050857565b6102d0565b6104d3565b6104ce565b909182601f830112156105515781359167ffffffffffffffff831161054c57602001926020830284011161054757565b6102d0565b6104d3565b6104ce565b906060828203126105d857600082013567ffffffffffffffff81116105d357816105819184016104d8565b929093602082013567ffffffffffffffff81116105ce57836105a4918401610517565b929093604082013567ffffffffffffffff81116105c9576105c59201610517565b9091565b6102cb565b6102cb565b6102cb565b61009e565b60000190565b34610618576106026105f6366004610556565b94939093929192610d2a565b61060a610093565b80610614816105dd565b0390f35b610099565b61062a6004600090610252565b90565b3461065d5761063d3660046100a3565b61065961064861061d565b610650610093565b91829182610280565b0390f35b610099565b600080fd5b600090565b600090565b600090565b601f801991011690565b634e487b7160e01b600052604160045260246000fd5b906106a090610676565b810190811067ffffffffffffffff8211176106ba57604052565b610680565b906106d26106cb610093565b9283610696565b565b67ffffffffffffffff81116106e95760200290565b610680565b6106fa6106ff916106d4565b6106bf565b90565b9061070c90610270565b9052565b60001c90565b61072261072791610710565b610167565b90565b6107349054610716565b90565b600080fd5b60e01b90565b600091031261074d57565b61009e565b9037565b6107639161010091610752565b565b50600490565b905090565b90565b61077c90610270565b9052565b9061078d81602093610773565b0190565b60200190565b6107b36107ad6107a683610765565b809461076b565b91610770565b6000915b8383106107c45750505050565b6107da6107d46001928451610780565b92610791565b920191906107b7565b916101006108079294936108006101808201966000830190610756565b0190610797565b565b610811610093565b3d6000823e3d90fd5b61082661082b91610710565b610237565b90565b610838905461081a565b90565b90565b61085261084d6108579261083b565b6101a8565b61037c565b90565b90565b61087161086c6108769261085a565b6101a8565b61037c565b90565b916108d790959495610889610667565b5061089261066c565b5061089b610671565b506108a4610671565b506108ad610671565b506108b6610671565b506108ce6108c460046106ee565b9360008501610702565b60208301610702565b6108e48360408301610702565b6108f18560608301610702565b906109046108ff600161072a565b6101d3565b916323572511919092803b156109e55761093160009361093c610925610093565b9687958694859461073c565b8452600484016107e3565b03915afa90816109b8575b50156000146109b3576001610985575b600192600093610967600361082e565b9361097e610975600461082e565b9493929661085d565b9493929190565b600092600193610995600361082e565b936109ac6109a3600461082e565b9493929661083e565b9493929190565b610957565b6109d89060003d81116109de575b6109d08183610696565b810190610742565b38610947565b503d6109c6565b610737565b6109f66109fb91610710565b6100b7565b90565b610a0890546109ea565b90565b60209181520190565b60007f4e6f742069737375657200000000000000000000000000000000000000000000910152565b610a49600a602092610a0b565b610a5281610a14565b0190565b610a6c9060208101906000818303910152610a3c565b90565b15610a7657565b610a7e610093565b62461bcd60e51b815280610a9460048201610a56565b0390fd5b90610acb9594939291610ac633610ac0610aba610ab560026109fe565b610103565b91610103565b14610a6f565b610c8a565b565b610ad9610ade91610710565b61042b565b90565b610aeb9054610acd565b90565b60209181520190565b90565b60209181520190565b90826000939282370152565b9190610b2981610b2281610b2e95610afa565b8095610b03565b610676565b0190565b90610b3d9291610b0f565b90565b600080fd5b600080fd5b600080fd5b9035600160200382360303811215610b9057016020813591019167ffffffffffffffff8211610b8b576001820236038313610b8657565b610b45565b610b40565b610b4a565b60200190565b9181610ba691610aee565b9081610bb760208302840194610af7565b92836000925b848410610bcd5750505050505090565b9091929394956020610bf9610bf38385600195038852610bed8b88610b4f565b90610b32565b98610b95565b940194019294939190610bbd565b60209181520190565b600080fd5b909182610c2191610c07565b9160018060fb1b038111610c445782916020610c409202938491610752565b0190565b610c10565b94929093610c6b610c879795610c7994606089019189830360008b0152610b9b565b918683036020880152610c15565b926040818503910152610c15565b90565b9194909293610ca1610c9c6000610ae1565b61046a565b9263b163337d90949695919295843b15610d2557600096610cd6948894610ce193610cca610093565b9b8c9a8b998a9861073c565b885260048801610c49565b03925af18015610d2057610cf3575b50565b610d139060003d8111610d19575b610d0b8183610696565b810190610742565b38610cf0565b503d610d01565b610809565b610737565b90610d389594939291610a98565b56fea2646970667358221220ce61f88d6bf1042aeeae8101d4f68a819f7d986e6f2f5971a5c9a72208571bb664736f6c634300081e0033",
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

// CheckCredential is a free data retrieval call binding the contract method 0x9d2cd569.
//
// Solidity: function checkCredential(uint256[8] proof, uint256 pubKeyX, uint256 pubKeyY, uint256 token, uint256 epoch) view returns(bool valid, uint8 errorCode, uint256 pubKeyXx, uint256 pubKeyYy, uint256 usedToken, uint256 usedEpoch)
func (_Verifier *VerifierCaller) CheckCredential(opts *bind.CallOpts, proof [8]*big.Int, pubKeyX *big.Int, pubKeyY *big.Int, token *big.Int, epoch *big.Int) (struct {
	Valid     bool
	ErrorCode uint8
	PubKeyXx  *big.Int
	PubKeyYy  *big.Int
	UsedToken *big.Int
	UsedEpoch *big.Int
}, error) {
	var out []interface{}
	err := _Verifier.contract.Call(opts, &out, "checkCredential", proof, pubKeyX, pubKeyY, token, epoch)

	outstruct := new(struct {
		Valid     bool
		ErrorCode uint8
		PubKeyXx  *big.Int
		PubKeyYy  *big.Int
		UsedToken *big.Int
		UsedEpoch *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Valid = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.ErrorCode = *abi.ConvertType(out[1], new(uint8)).(*uint8)
	outstruct.PubKeyXx = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.PubKeyYy = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.UsedToken = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.UsedEpoch = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// CheckCredential is a free data retrieval call binding the contract method 0x9d2cd569.
//
// Solidity: function checkCredential(uint256[8] proof, uint256 pubKeyX, uint256 pubKeyY, uint256 token, uint256 epoch) view returns(bool valid, uint8 errorCode, uint256 pubKeyXx, uint256 pubKeyYy, uint256 usedToken, uint256 usedEpoch)
func (_Verifier *VerifierSession) CheckCredential(proof [8]*big.Int, pubKeyX *big.Int, pubKeyY *big.Int, token *big.Int, epoch *big.Int) (struct {
	Valid     bool
	ErrorCode uint8
	PubKeyXx  *big.Int
	PubKeyYy  *big.Int
	UsedToken *big.Int
	UsedEpoch *big.Int
}, error) {
	return _Verifier.Contract.CheckCredential(&_Verifier.CallOpts, proof, pubKeyX, pubKeyY, token, epoch)
}

// CheckCredential is a free data retrieval call binding the contract method 0x9d2cd569.
//
// Solidity: function checkCredential(uint256[8] proof, uint256 pubKeyX, uint256 pubKeyY, uint256 token, uint256 epoch) view returns(bool valid, uint8 errorCode, uint256 pubKeyXx, uint256 pubKeyYy, uint256 usedToken, uint256 usedEpoch)
func (_Verifier *VerifierCallerSession) CheckCredential(proof [8]*big.Int, pubKeyX *big.Int, pubKeyY *big.Int, token *big.Int, epoch *big.Int) (struct {
	Valid     bool
	ErrorCode uint8
	PubKeyXx  *big.Int
	PubKeyYy  *big.Int
	UsedToken *big.Int
	UsedEpoch *big.Int
}, error) {
	return _Verifier.Contract.CheckCredential(&_Verifier.CallOpts, proof, pubKeyX, pubKeyY, token, epoch)
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
