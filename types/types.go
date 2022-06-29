package types

import (
	"os"
	"path/filepath"

	sdk "github.com/cosmos/cosmos-sdk/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

var (
	NodeHome        string
	DefaultNodeHome string
)

func init() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	DefaultNodeHome = filepath.Join(userHomeDir, ".misestm")
	NodeHome = DefaultNodeHome
}

const (
	DefaultEndpoint       string = "http://localhost:1317/"
	DefaultChainID        string = "mainnet"
	DefaultPassPhrase     string = "mises.site"
	DefaultRpcURI         string = "tcp://127.0.0.1:26657"
	AddressPrefix         string = "mises"
	MisesIDPrefix                = "did:mises:"
	MisesAppIDPrefix             = "did:misesapp:"
	ErrorInvalidLeaseTime string = "Invalid lease time"
	ErrorKeyIsRequired    string = "Key is required"
	ErrorValueIsRequired  string = "Value is required"
	ErrorKeyFormat        string = "Key format error"
)

type MSdkOption struct {
	ChainID    string //'mainnet' for the mainnet
	PassPhrase string //8 chars needed, default is 'mises.site'
	Debug      bool
	RpcURI     string
}

type MSdk interface {
	UserMgr() MUserMgr
	SetEndpoint(endpoint string) error
	TestConnection() error
	SetLogLevel(level int) error
	Login(site string, permissions []string) (string, error)
	VerifyLogin(auth string) (string, string, error)
}

type MUserInfo interface {
	Name() string
	Gender() string
	Avatar() string   //url of avatar
	HomePage() string //url of homepage
	Emails() []string
	Telphones() []string
	Intro() string
}
type MUser interface {
	MisesID() string
	Info() MUserInfo
	GetFollow(appDid string) []string
	LoadKeyStore(passPhrase string) error
	IsRegistered() error

	SetInfo(info MUserInfo) (string, error)
	SetFollow(followingId string, op bool, appDid string) (string, error)

	Signer() MSigner
}

type MUserMgr interface {
	CreateUser(mnemonic string, passPhrase string) (MUser, error)
	ListUsers() []MUser
	AddUser(user MUser)
	SetActiveUser(userDid string) error
	ActiveUser() MUser
}

type MSigner interface {
	MisesID() string
	Sign(msg string) (string, error)
	PubKey() string
	AesKey() ([]byte, error)
}

type MisesAppCmd interface {
	MisesUID() string
	PubKey() string
	TxID() string
	SetTxID(txid string)
	TrackID() string
	SetTrackID(trackerID string)

	WaitTx() bool
	SetWaitTx(wait bool)
}

type MisesAppCmdListener interface {
	OnTxGenerated(cmd MisesAppCmd)
	OnSucceed(cmd MisesAppCmd)
	OnFailed(cmd MisesAppCmd, err error)
}

type MisesEventStreamingListener interface {
	OnTxEvent(*tmtypes.EventDataTx)
	OnNewBlockHeaderEvent(*tmtypes.EventDataNewBlockHeader)
	OnEventStreamingTerminated()
}

type MApp interface {
	MisesID() string //did:misesapp:0123456789abcdef
	Info() MAppInfo

	Init(info MAppInfo, options MSdkOption) error

	SetListener(listener MisesAppCmdListener)

	AddAuth(misesUID string, permissions []string)

	RunAsync(cmd MisesAppCmd, wait bool) error

	RunSync(cmd MisesAppCmd) error

	Signer() MSigner

	NewRegisterUserCmd(uid string, pubkey string, feeGrantedPerDay int64) MisesAppCmd
	NewFaucetCmd(uid string, pubkey string, coin int64) MisesAppCmd

	StartEventStreaming(listener MisesEventStreamingListener) error
	ParseEvent(header *tmtypes.EventDataNewBlockHeader, tx *tmtypes.EventDataTx) (*sdk.TxResponse, error)
}

type MAppInfo interface {
	AppName() string //
	IconURL() string
	HomeURL() string
	Domains() []string //app
	Developer() string
}

type MPublicKey interface {
	KeyDid() string       // key did "did:mises:123456789abcdefghi#keys-1"
	KeyType() string      //Ed25519VerificationKey2020
	KeyMultibase() string //
}

type MAgent interface {
	PublicKeyTx() (MPublicKey, error)
	PublicKeyEnc() (MPublicKey, error)
	VerifySign(data string) error
}
