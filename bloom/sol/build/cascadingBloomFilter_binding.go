// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bloom

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

// BloomMetaData contains all meta data concerning the Bloom contract.
var BloomMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"layerIdx\",\"type\":\"uint256\"}],\"name\":\"chunkCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"layerIdx\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"chunkIdx\",\"type\":\"uint256\"}],\"name\":\"getChunk\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"layerIdx\",\"type\":\"uint256\"}],\"name\":\"getLayerMetadata\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"chunkSizeBytes_\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"filterSizeBits_\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"k_\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"layerCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"token\",\"type\":\"bytes\"}],\"name\":\"testToken\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"layerIdx\",\"type\":\"uint256\"}],\"name\":\"totalBytesOfLayer\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[][]\",\"name\":\"newChunksByLayer\",\"type\":\"bytes[][]\"},{\"internalType\":\"uint256[]\",\"name\":\"ks\",\"type\":\"uint256[]\"}],\"name\":\"updateCascade\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234602257600e60a5565b60146026565b6125286100b1823961252890f35b602c565b60405190565b600080fd5b60001b90565b90604660018060a01b03916031565b9181191691161790565b60018060a01b031690565b90565b606d60696071926050565b605b565b6050565b90565b607b90605e565b90565b6085906074565b90565b90565b90609b609760a192607e565b6088565b82546037565b9055565b60ae336000608b565b56fe60806040526004361015610013575b610618565b61001e6000356100ad565b806323b3bf32146100a85780634d75704a146100a357806351decc7d1461009e57806356e7f6c7146100995780638da5cb5b14610094578063a50e2b431461008f578063ad88b6e41461008a578063d423db2a146100855763f2fde38b0361000e576105e5565b61056c565b610498565b6103ac565b610344565b6102d5565b61028e565b610166565b610131565b60e01c90565b60405190565b600080fd5b600080fd5b600080fd5b90565b6100d4816100c8565b036100db57565b600080fd5b905035906100ed826100cb565b565b9060208282031261010957610106916000016100e0565b90565b6100be565b610117906100c8565b9052565b919061012f9060006020850194019061010e565b565b346101615761015d61014c6101473660046100ef565b6107ba565b6101546100b3565b9182918261011b565b0390f35b6100b9565b346101965761019261018161017c3660046100ef565b610819565b6101896100b3565b9182918261011b565b0390f35b6100b9565b600080fd5b600080fd5b600080fd5b909182601f830112156101e45781359167ffffffffffffffff83116101df5760200192602083028401116101da57565b6101a5565b6101a0565b61019b565b909182601f830112156102235781359167ffffffffffffffff831161021e57602001926020830284011161021957565b6101a5565b6101a0565b61019b565b909160408284031261028357600082013567ffffffffffffffff811161027e57836102549184016101aa565b929093602082013567ffffffffffffffff81116102795761027592016101e9565b9091565b6100c3565b6100c3565b6100be565b60000190565b346102c0576102aa6102a1366004610228565b92919091611ab0565b6102b26100b3565b806102bc81610288565b0390f35b6100b9565b60009103126102d057565b6100be565b34610305576102e53660046102c5565b6103016102f0611abe565b6102f86100b3565b9182918261011b565b0390f35b6100b9565b60018060a01b031690565b61031e9061030a565b90565b61032a90610315565b9052565b919061034290600060208501940190610321565b565b34610374576103543660046102c5565b61037061035f611ad9565b6103676100b3565b9182918261032e565b0390f35b6100b9565b6040906103a36103aa94969593966103996060840198600085019061010e565b602083019061010e565b019061010e565b565b346103df576103db6103c76103c23660046100ef565b611aef565b6103d29391936100b3565b93849384610379565b0390f35b6100b9565b919060408382031261040d578061040161040a92600086016100e0565b936020016100e0565b90565b6100be565b5190565b60209181520190565b60005b838110610433575050906000910152565b806020918301518185015201610422565b601f801991011690565b61046d61047660209361047b9361046481610412565b93848093610416565b9586910161041f565b610444565b0190565b610495916020820191600081840391015261044e565b90565b346104c9576104c56104b46104ae3660046103e4565b90611ccb565b6104bc6100b3565b9182918261047f565b0390f35b6100b9565b909182601f830112156105085781359167ffffffffffffffff83116105035760200192600183028401116104fe57565b6101a5565b6101a0565b61019b565b9060208282031261053f57600082013567ffffffffffffffff811161053a5761053692016104ce565b9091565b6100c3565b6100be565b151590565b61055290610544565b9052565b919061056a90600060208501940190610549565b565b3461059d5761059961058861058236600461050d565b90612106565b6105906100b3565b91829182610556565b0390f35b6100b9565b6105ab81610315565b036105b257565b600080fd5b905035906105c4826105a2565b565b906020828203126105e0576105dd916000016105b7565b90565b6100be565b34610613576105fd6105f83660046105c6565b6124e7565b6106056100b3565b8061060f81610288565b0390f35b6100b9565b600080fd5b600090565b5490565b60209181520190565b60007f496e76616c6964206c6179657220696e64657800000000000000000000000000910152565b6106646013602092610626565b61066d8161062f565b0190565b6106879060208101906000818303910152610657565b90565b1561069157565b6106996100b3565b62461bcd60e51b8152806106af60048201610671565b0390fd5b634e487b7160e01b600052603260045260246000fd5b600052602060002090565b600052602060002090565b9060206106f1818306601f03936106d4565b91040191565b61070081610622565b82101561071b576107126004916106c9565b91020190600090565b6106b3565b60001c90565b90565b61073561073a91610720565b610726565b90565b6107479054610729565b90565b90565b90565b61076461075f6107699261074a565b61074d565b6100c8565b90565b634e487b7160e01b600052601260045260246000fd5b634e487b7160e01b600052601160045260246000fd5b6107a46107aa916100c8565b916100c8565b9081156107b5570490565b61076c565b61080260026107fb610812936107ce61061d565b506107f4816107ee6107e86107e36001610622565b6100c8565b916100c8565b1061068a565b60016106f7565b500161073d565b61080c6008610750565b90610798565b90565b5490565b600061085761085e9261082a61061d565b506108508161084a61084461083f6001610622565b6100c8565b916100c8565b1061068a565b60016106f7565b5001610815565b90565b60018060a01b031690565b61087861087d91610720565b610861565b90565b61088a905461086c565b90565b60007f43616c6c6572206973206e6f74206f776e657200000000000000000000000000910152565b6108c26013602092610626565b6108cb8161088d565b0190565b6108e590602081019060008183039101526108b5565b90565b156108ef57565b6108f76100b3565b62461bcd60e51b81528061090d600482016108cf565b0390fd5b9061094293929161093d3361093761093161092c6000610880565b610315565b91610315565b146108e8565b611762565b565b5090565b5090565b60007f4d75737420737570706c79206f6e65206b20706572206c617965720000000000910152565b610981601b602092610626565b61098a8161094c565b0190565b6109a49060208101906000818303910152610974565b90565b156109ae57565b6109b66100b3565b62461bcd60e51b8152806109cc6004820161098e565b0390fd5b90565b6109e76109e26109ec926109d0565b61074d565b6100c8565b90565b60007f4e656564206174206c65617374206f6e65206c61796572000000000000000000910152565b610a246017602092610626565b610a2d816109ef565b0190565b610a479060208101906000818303910152610a17565b90565b15610a5157565b610a596100b3565b62461bcd60e51b815280610a6f60048201610a31565b0390fd5b634e487b7160e01b600052604160045260246000fd5b610a98610a9e919392936100c8565b926100c8565b91610aaa8382026100c8565b928184041490151715610ab957565b610782565b610ac9906004610a89565b90565b600190818003010490565b600052602060002090565b634e487b7160e01b600052602260045260246000fd5b9060016002830492168015610b18575b6020831014610b1357565b610ae2565b91607f1691610b08565b600052602060002090565b1c90565b90610b459060001990602003600802610b2d565b8154169055565b1b90565b91906008610b6c910291610b6660001984610b4c565b92610b4c565b9181191691161790565b610b8a610b85610b8f926100c8565b61074d565b6100c8565b90565b90565b9190610bab610ba6610bb393610b76565b610b92565b908354610b50565b9055565b610bc991610bc361061d565b91610b95565b565b5b818110610bd7575050565b80610be56000600193610bb7565b01610bcc565b90610bfc9060001990600802610b2d565b191690565b81610c0b91610beb565b906002021790565b90600091610c2b610c2382610b22565b928354610c01565b905555565b601f602091010490565b91929060208210600014610c9457601f8411600114610c6457610c5e929350610c01565b90555b5b565b5090610c8a610c8f936001610c81610c7b85610b22565b92610c30565b82019101610bcb565b610c13565b610c61565b50610ccb8293610ca5600194610b22565b610cc4610cb185610c30565b820192601f861680610cd6575b50610c30565b0190610bcb565b600202179055610c62565b610ce290888603610b31565b38610cbe565b929091680100000000000000008211610d4a57602011600014610d3b5760208110600014610d1f57610d1991610c01565b90555b5b565b60019160ff1916610d2f84610b22565b55600202019055610d1c565b60019150600202019055610d1d565b610a73565b908154610d5b81610af8565b90818311610d84575b818310610d72575b50505050565b610d7b93610c3a565b38808080610d6c565b610d9083838387610ce8565b610d64565b6000610da091610d4f565b565b634e487b7160e01b600052600060045260246000fd5b90600003610dcb57610dc990610d95565b565b610da2565b5b818110610ddc575050565b80610dea6000600193610db8565b01610dd1565b9091828110610dff575b505050565b610e1d610e17610e11610e2895610acc565b92610acc565b92610ad7565b918201910190610dd0565b388080610dfa565b90680100000000000000008111610e595781610e4e610e5793610815565b90828155610df0565b565b610a73565b6000610e6991610e30565b565b90600003610e7e57610e7c90610e5e565b565b610da2565b60006003610eb792610e9783808301610e6b565b610ea48360018301610bb7565b610eb18360028301610bb7565b01610bb7565b565b90600003610ecc57610eca90610e83565b565b610da2565b5b818110610edd575050565b80610eeb6000600493610eb9565b01610ed2565b9091828110610f00575b505050565b610f1e610f18610f12610f2995610abe565b92610abe565b926106c9565b918201910190610ed1565b388080610efb565b90680100000000000000008111610f5a5781610f4f610f5893610622565b90828155610ef1565b565b610a73565b6000610f6a91610f31565b565b90600003610f7f57610f7d90610f5f565b565b610da2565b6001610f9091016100c8565b90565b600080fd5b600080fd5b600080fd5b903590600160200381360303821215610fe4570180359067ffffffffffffffff8211610fdf57602001916020820236038313610fda57565b610f9d565b610f98565b610f93565b908210156110045760206110009202810190610fa2565b9091565b6106b3565b9190811015611019576020020190565b6106b3565b35611028816100cb565b90565b60007f6b206d757374206265203e203000000000000000000000000000000000000000910152565b611060600d602092610626565b6110698161102b565b0190565b6110839060208101906000818303910152611053565b90565b1561108d57565b6110956100b3565b62461bcd60e51b8152806110ab6004820161106d565b0390fd5b5090565b60207f6e6b000000000000000000000000000000000000000000000000000000000000917f45616368206c61796572206e65656473206d6f7265207468616e20312063687560008201520152565b61110e6022604092610626565b611117816110b3565b0190565b6111319060208101906000818303910152611101565b90565b1561113b57565b6111436100b3565b62461bcd60e51b8152806111596004820161111b565b0390fd5b90359060016020038136030382121561119f570180359067ffffffffffffffff821161119a5760200191600182023603831361119557565b610f9d565b610f98565b610f93565b908210156111bf5760206111bb920281019061115d565b9091565b6106b3565b5090565b60007f4368756e6b732063616e6e6f74206265207a65726f2d6c656e67746800000000910152565b6111fd601c602092610626565b611206816111c8565b0190565b61122090602081019060008183039101526111f0565b90565b1561122a57565b6112326100b3565b62461bcd60e51b8152806112486004820161120a565b0390fd5b90565b61126361125e6112689261124c565b61074d565b6100c8565b90565b61127a611280919392936100c8565b926100c8565b820391821161128b57565b610782565b60007f4368756e6b73206d69736d617463682073697a65000000000000000000000000910152565b6112c56014602092610626565b6112ce81611290565b0190565b6112e890602081019060008183039101526112b8565b90565b156112f257565b6112fa6100b3565b62461bcd60e51b815280611310600482016112d2565b0390fd5b611323611329919392936100c8565b926100c8565b820180921161133457565b610782565b90565b5490565b600052602060002090565b6113548161133c565b82101561136f57611366600491611340565b91020190600090565b6106b3565b61137d8161133c565b680100000000000000008110156113a15761139d9160018201815561134b565b9091565b610a73565b90565b60001b90565b906113bc600019916113a9565b9181191691161790565b906113db6113d66113e292610b76565b610b92565b82546113af565b9055565b906113f090610444565b810190811067ffffffffffffffff82111761140a57604052565b610a73565b9061142261141b6100b3565b92836113e6565b565b67ffffffffffffffff811161143c5760208091020190565b610a73565b9061145361144e83611424565b61140f565b918252565b606090565b60005b82811061146c57505050565b602090611477611458565b8184015201611460565b906114a661148e83611441565b9260208061149c8693611424565b920191039061145d565b565b60200190565b5190565b5190565b9190601f81116114c6575b505050565b6114d26114f793610b22565b9060206114de84610c30565b830193106114ff575b6114f090610c30565b0190610bcb565b3880806114c1565b91506114f0819290506114e7565b9061151781610412565b9067ffffffffffffffff82116115d95761153b826115358554610af8565b856114b6565b602090601f83116001146115705791809161155f93600092611564575b5050610c01565b90555b565b90915001513880611558565b601f1983169161157f85610b22565b9260005b8181106115c1575091600293918560019694106115a7575b50505002019055611562565b6115b7910151601f841690610beb565b905538808061159b565b91936020600181928787015181550195019201611583565b610a73565b906115e89161150d565b565b61160f6116096115f9846114ae565b936116048585610e30565b6114a8565b91610ad7565b6000915b8383106116205750505050565b600160208261163861163284956114b2565b866115de565b01920192019190611613565b9061164e916115ea565b565b61165981610815565b8210156116745761166b600191610ad7565b91020190600090565b6106b3565b9161168490826111c4565b9067ffffffffffffffff8211611746576116a8826116a28554610af8565b856114b6565b600090601f83116001146116dd579180916116cc936000926116d1575b5050610c01565b90555b565b909150013538806116c5565b601f198316916116ec85610b22565b9260005b81811061172e57509160029391856001969410611714575b505050020190556116cf565b611724910135601f841690610beb565b9055388080611708565b919360206001819287870135815501950192016116f0565b610a73565b92919061175d5761175b92611679565b565b610da2565b91909392611796611774848790610944565b61179061178a611785868690610948565b6100c8565b916100c8565b146109a7565b6117bd6117a4848790610944565b6117b76117b160006109d3565b916100c8565b11610a4a565b6117c960006001610f6c565b6117d360006109d3565b5b806117f16117eb6117e6878a90610944565b6100c8565b916100c8565b1015611aa85761180384878391610fe9565b9390939461181b61181683868691611009565b61101e565b966118398861183361182d60006109d3565b916100c8565b11611086565b6118606118478789906110af565b61185a61185460006109d3565b916100c8565b11611134565b61187e611878878961187260006109d3565b916111a4565b906111c4565b9461189c8661189661189060006109d3565b916100c8565b11611223565b6118a660006109d3565b936118b160006109d3565b945b856118cf6118c96118c48d8d6110af565b6100c8565b916100c8565b101561197457611930611936918a8c8961190e6119086119036118f38686906110af565b6118fd600161124f565b9061126b565b6100c8565b916100c8565b1061193c575b6119249061192a92908b916111a4565b906111c4565b90611314565b95610f84565b946118b3565b61194f6119559161196d93908c916111a4565b906111c4565b6119676119618d6100c8565b916100c8565b146112eb565b8a8c611914565b9093969a929599919450600161198990611339565b61199290611374565b50506001806119a090610622565b60016119ab9061124f565b6119b49161126b565b6119bd916106f7565b506119c7906113a6565b9a8b600101906119d6916113c6565b60086119e190610750565b6119ea91610a89565b8a600201906119f8916113c6565b8960030190611a06916113c6565b8587611a11916110af565b611a1a90611481565b8960000190611a2891611644565b6000611a33906109d3565b5b80611a51611a4b611a468a8c906110af565b6100c8565b916100c8565b1015611a8c57611a8790611a8260008c611a7c611a718d8d9087916111a4565b939092018590611650565b9061174b565b610f84565b611a34565b509450945094611a9d919650610f84565b9493909291946117d4565b505050509050565b90611abc939291610911565b565b611ac661061d565b50611ad16001610622565b90565b600090565b611ae1611ad4565b50611aec6000610880565b90565b611b3d611b4391611afe61061d565b50611b0761061d565b50611b1061061d565b50611b3681611b30611b2a611b256001610622565b6100c8565b916100c8565b1061068a565b60016106f7565b506113a6565b90611b506001830161073d565b90611b696003611b626002860161073d565b940161073d565b91929190565b606090565b60007f496e76616c6964206368756e6b20696e64657800000000000000000000000000910152565b611ba96013602092610626565b611bb281611b74565b0190565b611bcc9060208101906000818303910152611b9c565b90565b15611bd657565b611bde6100b3565b62461bcd60e51b815280611bf460048201611bb6565b0390fd5b60209181520190565b9060009291805490611c1c611c1583610af8565b8094611bf8565b91600181169081600014611c755750600114611c38575b505050565b611c459192939450610b22565b916000925b818410611c5d5750500190388080611c33565b60018160209295939554848601520191019290611c4a565b92949550505060ff1916825215156020020190388080611c33565b90611c9a91611c01565b90565b90611cbd611cb692611cad6100b3565b93848092611c90565b03836113e6565b565b611cc890611c9d565b90565b611d48916000611d16611d10611d4294611ce3611b6f565b50611d0981611d03611cfd611cf86001610622565b6100c8565b916100c8565b1061068a565b60016106f7565b506113a6565b611d3c83611d36611d30611d2b868601610815565b6100c8565b916100c8565b10611bcf565b01611650565b50611cbf565b90565b600090565b60007f4e6f206c617965727320696e697469616c697a65640000000000000000000000910152565b611d856015602092610626565b611d8e81611d50565b0190565b611da89060208101906000818303910152611d78565b90565b15611db257565b611dba6100b3565b62461bcd60e51b815280611dd060048201611d92565b0390fd5b600080fd5b67ffffffffffffffff8111611df757611df3602091610444565b0190565b610a73565b90826000939282370152565b90929192611e1d611e1882611dd9565b61140f565b93818552602085019082840111611e3957611e3792611dfc565b565b611dd4565b611e49913691611e08565b90565b60200190565b611e5e611e6391610720565b610b76565b90565b90565b60ff1690565b611e83611e7e611e8892611e66565b61074d565b611e69565b90565b611eaa90611ea4611e9e611eaf94611e69565b916100c8565b90610b2d565b6100c8565b90565b67ffffffffffffffff1690565b611ed3611ece611ed8926100c8565b61074d565b611eb2565b90565b90565b611ef2611eed611ef792611edb565b61074d565b611e69565b90565b90565b611f11611f0c611f1692611efa565b61074d565b6100c8565b90565b90565b611f30611f2b611f3592611eb2565b61074d565b6100c8565b90565b611f44611f4a916100c8565b916100c8565b908115611f55570690565b61076c565b90565b611f71611f6c611f7692611f5a565b61074d565b611e69565b90565b5490565b600052602060002090565b611f9181611f79565b821015611fac57611fa3600191611f7d565b91020190600090565b6106b3565b90565b611fbe9054610af8565b90565b90611fcb82611fb4565b80821015611ff957602011600014611fe95760209006601f0390915b565b611ff2916106df565b9091611fe7565b6106b3565b60f81b90565b61200d90611ffe565b90565b6120209060086120259302610b2d565b612004565b90565b906120339154612010565b90565b60f81c90565b61205061204b61205592611e69565b61074d565b611e69565b90565b61206461206991612036565b61203c565b90565b90565b61208361207e6120889261206c565b61074d565b6100c8565b90565b61209f61209a6120a4926100c8565b61074d565b611e69565b90565b6120c6906120c06120ba6120cb94611e69565b91611e69565b90610b2d565b611e69565b90565b6120e26120dd6120e79261124c565b61074d565b611e69565b90565b6120fe6120f9612103926109d0565b61074d565b611e69565b90565b9061214461215c91612116611d4b565b506121216001610622565b9361213f8561213961213360006109d3565b916100c8565b11611dab565b611e3e565b61215661215082610412565b91611e4c565b20611e52565b6121a661218f61217e6121798461217360c0611e6f565b90611e8b565b611ebf565b926121896080611ede565b90611e8b565b6121a067ffffffffffffffff611efd565b16611ebf565b916121b160006109d3565b925b836121c66121c0846100c8565b916100c8565b101561235d576121e460026121dd600187906106f7565b500161073d565b906121fc60016121f58188906106f7565b500161073d565b93612215600361220e600189906106f7565b500161073d565b9561222e6000612227600184906106f7565b5001611f19565b9461223960006109d3565b5b8061224d6122478b6100c8565b916100c8565b101561234457612305886122ff6122ee6122e96122e38c6122de6122d861229d8f8f8f612292612298928f9261228561228c91611f1c565b9391611f1c565b90610a89565b90611314565b611f38565b976122d26122b58a6122af6003611f5d565b90611e8b565b916122cc6122c4848390610798565b939184610a89565b9061126b565b93611f88565b50611fb1565b611fc1565b90612028565b612058565b916122f9600761206f565b1661208b565b906120a7565b61230f60016120ce565b1661232361231d60006120ea565b91611e69565b146123365761233190610f84565b61223a565b505050505050505050600090565b50945091945094506123569150610f84565b92916121b3565b50505050600190565b6123949061238f3361238961238361237e6000610880565b610315565b91610315565b146108e8565b6124b4565b565b6123aa6123a56123af926109d0565b61074d565b61030a565b90565b6123bb90612396565b90565b60007f4e6577206f776e6572206973207a65726f206164647265737300000000000000910152565b6123f36019602092610626565b6123fc816123be565b0190565b61241690602081019060008183039101526123e6565b90565b1561242057565b6124286100b3565b62461bcd60e51b81528061243e60048201612400565b0390fd5b9061245360018060a01b03916113a9565b9181191691161790565b61247161246c6124769261030a565b61074d565b61030a565b90565b6124829061245d565b90565b61248e90612479565b90565b90565b906124a96124a46124b092612485565b612491565b8254612442565b9055565b6124e5906124de816124d76124d16124cc60006123b2565b610315565b91610315565b1415612419565b6000612494565b565b6124f090612366565b56fea264697066735822122044d00c9edebb86655d3022b8a188cde3eedf0970e9ea44aecadec409b108fdb764736f6c634300081e0033",
}

// BloomABI is the input ABI used to generate the binding from.
// Deprecated: Use BloomMetaData.ABI instead.
var BloomABI = BloomMetaData.ABI

// BloomBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use BloomMetaData.Bin instead.
var BloomBin = BloomMetaData.Bin

// DeployBloom deploys a new Ethereum contract, binding an instance of Bloom to it.
func DeployBloom(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Bloom, error) {
	parsed, err := BloomMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(BloomBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Bloom{BloomCaller: BloomCaller{contract: contract}, BloomTransactor: BloomTransactor{contract: contract}, BloomFilterer: BloomFilterer{contract: contract}}, nil
}

// Bloom is an auto generated Go binding around an Ethereum contract.
type Bloom struct {
	BloomCaller     // Read-only binding to the contract
	BloomTransactor // Write-only binding to the contract
	BloomFilterer   // Log filterer for contract events
}

// BloomCaller is an auto generated read-only Go binding around an Ethereum contract.
type BloomCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BloomTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BloomTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BloomFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BloomFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BloomSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BloomSession struct {
	Contract     *Bloom            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BloomCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BloomCallerSession struct {
	Contract *BloomCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// BloomTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BloomTransactorSession struct {
	Contract     *BloomTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BloomRaw is an auto generated low-level Go binding around an Ethereum contract.
type BloomRaw struct {
	Contract *Bloom // Generic contract binding to access the raw methods on
}

// BloomCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BloomCallerRaw struct {
	Contract *BloomCaller // Generic read-only contract binding to access the raw methods on
}

// BloomTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BloomTransactorRaw struct {
	Contract *BloomTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBloom creates a new instance of Bloom, bound to a specific deployed contract.
func NewBloom(address common.Address, backend bind.ContractBackend) (*Bloom, error) {
	contract, err := bindBloom(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Bloom{BloomCaller: BloomCaller{contract: contract}, BloomTransactor: BloomTransactor{contract: contract}, BloomFilterer: BloomFilterer{contract: contract}}, nil
}

// NewBloomCaller creates a new read-only instance of Bloom, bound to a specific deployed contract.
func NewBloomCaller(address common.Address, caller bind.ContractCaller) (*BloomCaller, error) {
	contract, err := bindBloom(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BloomCaller{contract: contract}, nil
}

// NewBloomTransactor creates a new write-only instance of Bloom, bound to a specific deployed contract.
func NewBloomTransactor(address common.Address, transactor bind.ContractTransactor) (*BloomTransactor, error) {
	contract, err := bindBloom(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BloomTransactor{contract: contract}, nil
}

// NewBloomFilterer creates a new log filterer instance of Bloom, bound to a specific deployed contract.
func NewBloomFilterer(address common.Address, filterer bind.ContractFilterer) (*BloomFilterer, error) {
	contract, err := bindBloom(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BloomFilterer{contract: contract}, nil
}

// bindBloom binds a generic wrapper to an already deployed contract.
func bindBloom(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BloomMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Bloom *BloomRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Bloom.Contract.BloomCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Bloom *BloomRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Bloom.Contract.BloomTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Bloom *BloomRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Bloom.Contract.BloomTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Bloom *BloomCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Bloom.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Bloom *BloomTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Bloom.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Bloom *BloomTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Bloom.Contract.contract.Transact(opts, method, params...)
}

// ChunkCount is a free data retrieval call binding the contract method 0x4d75704a.
//
// Solidity: function chunkCount(uint256 layerIdx) view returns(uint256)
func (_Bloom *BloomCaller) ChunkCount(opts *bind.CallOpts, layerIdx *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Bloom.contract.Call(opts, &out, "chunkCount", layerIdx)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ChunkCount is a free data retrieval call binding the contract method 0x4d75704a.
//
// Solidity: function chunkCount(uint256 layerIdx) view returns(uint256)
func (_Bloom *BloomSession) ChunkCount(layerIdx *big.Int) (*big.Int, error) {
	return _Bloom.Contract.ChunkCount(&_Bloom.CallOpts, layerIdx)
}

// ChunkCount is a free data retrieval call binding the contract method 0x4d75704a.
//
// Solidity: function chunkCount(uint256 layerIdx) view returns(uint256)
func (_Bloom *BloomCallerSession) ChunkCount(layerIdx *big.Int) (*big.Int, error) {
	return _Bloom.Contract.ChunkCount(&_Bloom.CallOpts, layerIdx)
}

// GetChunk is a free data retrieval call binding the contract method 0xad88b6e4.
//
// Solidity: function getChunk(uint256 layerIdx, uint256 chunkIdx) view returns(bytes)
func (_Bloom *BloomCaller) GetChunk(opts *bind.CallOpts, layerIdx *big.Int, chunkIdx *big.Int) ([]byte, error) {
	var out []interface{}
	err := _Bloom.contract.Call(opts, &out, "getChunk", layerIdx, chunkIdx)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetChunk is a free data retrieval call binding the contract method 0xad88b6e4.
//
// Solidity: function getChunk(uint256 layerIdx, uint256 chunkIdx) view returns(bytes)
func (_Bloom *BloomSession) GetChunk(layerIdx *big.Int, chunkIdx *big.Int) ([]byte, error) {
	return _Bloom.Contract.GetChunk(&_Bloom.CallOpts, layerIdx, chunkIdx)
}

// GetChunk is a free data retrieval call binding the contract method 0xad88b6e4.
//
// Solidity: function getChunk(uint256 layerIdx, uint256 chunkIdx) view returns(bytes)
func (_Bloom *BloomCallerSession) GetChunk(layerIdx *big.Int, chunkIdx *big.Int) ([]byte, error) {
	return _Bloom.Contract.GetChunk(&_Bloom.CallOpts, layerIdx, chunkIdx)
}

// GetLayerMetadata is a free data retrieval call binding the contract method 0xa50e2b43.
//
// Solidity: function getLayerMetadata(uint256 layerIdx) view returns(uint256 chunkSizeBytes_, uint256 filterSizeBits_, uint256 k_)
func (_Bloom *BloomCaller) GetLayerMetadata(opts *bind.CallOpts, layerIdx *big.Int) (struct {
	ChunkSizeBytes *big.Int
	FilterSizeBits *big.Int
	K              *big.Int
}, error) {
	var out []interface{}
	err := _Bloom.contract.Call(opts, &out, "getLayerMetadata", layerIdx)

	outstruct := new(struct {
		ChunkSizeBytes *big.Int
		FilterSizeBits *big.Int
		K              *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ChunkSizeBytes = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.FilterSizeBits = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.K = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetLayerMetadata is a free data retrieval call binding the contract method 0xa50e2b43.
//
// Solidity: function getLayerMetadata(uint256 layerIdx) view returns(uint256 chunkSizeBytes_, uint256 filterSizeBits_, uint256 k_)
func (_Bloom *BloomSession) GetLayerMetadata(layerIdx *big.Int) (struct {
	ChunkSizeBytes *big.Int
	FilterSizeBits *big.Int
	K              *big.Int
}, error) {
	return _Bloom.Contract.GetLayerMetadata(&_Bloom.CallOpts, layerIdx)
}

// GetLayerMetadata is a free data retrieval call binding the contract method 0xa50e2b43.
//
// Solidity: function getLayerMetadata(uint256 layerIdx) view returns(uint256 chunkSizeBytes_, uint256 filterSizeBits_, uint256 k_)
func (_Bloom *BloomCallerSession) GetLayerMetadata(layerIdx *big.Int) (struct {
	ChunkSizeBytes *big.Int
	FilterSizeBits *big.Int
	K              *big.Int
}, error) {
	return _Bloom.Contract.GetLayerMetadata(&_Bloom.CallOpts, layerIdx)
}

// LayerCount is a free data retrieval call binding the contract method 0x56e7f6c7.
//
// Solidity: function layerCount() view returns(uint256)
func (_Bloom *BloomCaller) LayerCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Bloom.contract.Call(opts, &out, "layerCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LayerCount is a free data retrieval call binding the contract method 0x56e7f6c7.
//
// Solidity: function layerCount() view returns(uint256)
func (_Bloom *BloomSession) LayerCount() (*big.Int, error) {
	return _Bloom.Contract.LayerCount(&_Bloom.CallOpts)
}

// LayerCount is a free data retrieval call binding the contract method 0x56e7f6c7.
//
// Solidity: function layerCount() view returns(uint256)
func (_Bloom *BloomCallerSession) LayerCount() (*big.Int, error) {
	return _Bloom.Contract.LayerCount(&_Bloom.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Bloom *BloomCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Bloom.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Bloom *BloomSession) Owner() (common.Address, error) {
	return _Bloom.Contract.Owner(&_Bloom.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Bloom *BloomCallerSession) Owner() (common.Address, error) {
	return _Bloom.Contract.Owner(&_Bloom.CallOpts)
}

// TestToken is a free data retrieval call binding the contract method 0xd423db2a.
//
// Solidity: function testToken(bytes token) view returns(bool)
func (_Bloom *BloomCaller) TestToken(opts *bind.CallOpts, token []byte) (bool, error) {
	var out []interface{}
	err := _Bloom.contract.Call(opts, &out, "testToken", token)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// TestToken is a free data retrieval call binding the contract method 0xd423db2a.
//
// Solidity: function testToken(bytes token) view returns(bool)
func (_Bloom *BloomSession) TestToken(token []byte) (bool, error) {
	return _Bloom.Contract.TestToken(&_Bloom.CallOpts, token)
}

// TestToken is a free data retrieval call binding the contract method 0xd423db2a.
//
// Solidity: function testToken(bytes token) view returns(bool)
func (_Bloom *BloomCallerSession) TestToken(token []byte) (bool, error) {
	return _Bloom.Contract.TestToken(&_Bloom.CallOpts, token)
}

// TotalBytesOfLayer is a free data retrieval call binding the contract method 0x23b3bf32.
//
// Solidity: function totalBytesOfLayer(uint256 layerIdx) view returns(uint256)
func (_Bloom *BloomCaller) TotalBytesOfLayer(opts *bind.CallOpts, layerIdx *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Bloom.contract.Call(opts, &out, "totalBytesOfLayer", layerIdx)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalBytesOfLayer is a free data retrieval call binding the contract method 0x23b3bf32.
//
// Solidity: function totalBytesOfLayer(uint256 layerIdx) view returns(uint256)
func (_Bloom *BloomSession) TotalBytesOfLayer(layerIdx *big.Int) (*big.Int, error) {
	return _Bloom.Contract.TotalBytesOfLayer(&_Bloom.CallOpts, layerIdx)
}

// TotalBytesOfLayer is a free data retrieval call binding the contract method 0x23b3bf32.
//
// Solidity: function totalBytesOfLayer(uint256 layerIdx) view returns(uint256)
func (_Bloom *BloomCallerSession) TotalBytesOfLayer(layerIdx *big.Int) (*big.Int, error) {
	return _Bloom.Contract.TotalBytesOfLayer(&_Bloom.CallOpts, layerIdx)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Bloom *BloomTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Bloom.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Bloom *BloomSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Bloom.Contract.TransferOwnership(&_Bloom.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Bloom *BloomTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Bloom.Contract.TransferOwnership(&_Bloom.TransactOpts, newOwner)
}

// UpdateCascade is a paid mutator transaction binding the contract method 0x51decc7d.
//
// Solidity: function updateCascade(bytes[][] newChunksByLayer, uint256[] ks) returns()
func (_Bloom *BloomTransactor) UpdateCascade(opts *bind.TransactOpts, newChunksByLayer [][][]byte, ks []*big.Int) (*types.Transaction, error) {
	return _Bloom.contract.Transact(opts, "updateCascade", newChunksByLayer, ks)
}

// UpdateCascade is a paid mutator transaction binding the contract method 0x51decc7d.
//
// Solidity: function updateCascade(bytes[][] newChunksByLayer, uint256[] ks) returns()
func (_Bloom *BloomSession) UpdateCascade(newChunksByLayer [][][]byte, ks []*big.Int) (*types.Transaction, error) {
	return _Bloom.Contract.UpdateCascade(&_Bloom.TransactOpts, newChunksByLayer, ks)
}

// UpdateCascade is a paid mutator transaction binding the contract method 0x51decc7d.
//
// Solidity: function updateCascade(bytes[][] newChunksByLayer, uint256[] ks) returns()
func (_Bloom *BloomTransactorSession) UpdateCascade(newChunksByLayer [][][]byte, ks []*big.Int) (*types.Transaction, error) {
	return _Bloom.Contract.UpdateCascade(&_Bloom.TransactOpts, newChunksByLayer, ks)
}
