package types

import (
	"os"
	"path/filepath"
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
	DefaultChainID        string = "mises"
	DefaultPassPhrase     string = "mises.site"
	AddressPrefix         string = "mises"
	MisesIDPrefix                = "did:mises:"
	MisesAppIDPrefix             = "did:misesapp:"
	ErrorInvalidLeaseTime string = "Invalid lease time"
	ErrorKeyIsRequired    string = "Key is required"
	ErrorValueIsRequired  string = "Value is required"
	ErrorKeyFormat        string = "Key format error"
)

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

	WaitTx() bool
	SetWaitTx(wait bool)
}

type MisesAppCmdListener interface {
	OnTxGenerated(cmd MisesAppCmd)
	OnSucceed(cmd MisesAppCmd)
	OnFailed(cmd MisesAppCmd, err error)
}

type MApp interface {
	MisesID() string //did:misesapp:0123456789abcdef
	Info() MAppInfo

	Init(info MAppInfo, chainID string, passPhrase string) error

	SetListener(listener MisesAppCmdListener)

	AddAuth(misesUID string, permissions []string)

	RunAsync(cmd MisesAppCmd, wait bool) error

	RunSync(cmd MisesAppCmd) error

	Signer() MSigner

	NewRegisterUserCmd(uid string, pubkey string, feeGrantedPerDay int64) MisesAppCmd
	NewFaucetCmd(uid string, pubkey string, coin int64) MisesAppCmd
}

type MAppInfo interface {
	AppName() string //
	IconURL() string
	HomeURL() string
	Domains() []string //app
	Developer() string
}

/*
type MUserAuthorization interface {
	UserDid() string       //mises app did
	AppDid() string        //to
	Permissions() []string //user_info_r,  user_info_w, user_relation_r, user_relation_w
	ExpireTimestamp() int  //
	AppAuthorization() MAppAuth
}

type MAppInfo interface {
	ApppDid() string //did:mises:0123456789abcdef
	AppName() string //
	IconDid() string //udid of icon file did:mises:0123456789abcdef/icon
	IconThumb() []byte
	Domain() string //app
	Developer() string
}

type MAppAuth interface {
	AppInfo() MAppInfo
	MisesId()
	Permissions() []string //
	ExpireTimestamp() int  //
}

type MAppMgr interface {
	AddApp(app MAppAuth, removable bool) (MApp, error)
	ListApps() ([]MApp, error)
	RemoveApp(appDid string) error
}

type MEvent interface {
	EventID() int
}
type MEventCallback func(event MEvent) error

type MApp interface {
	AppDID() (string, error)
	AppDomain() string
	SetAppDomain(string)
	IsRegistered() (bool, error)
	Register(MAppInfo string, appDid string) error
	AddAuth(misesid string, permission []string) (MAppAuth, error)
	//	GenerateAuthorization(permisions []string) (MAppAuthorization, error)

	AddEventListener(userDid string, userAuth MUserAuthorization, callback MEventCallback)
	RemoveUserEventListener(userDid string, userAuth MUserAuthorization, callback MEventCallback) error
	ListFollow(whomDid string) (dids []string)
	Commit() error
	Cancel() error
	Agent() (MAgent, error)
}
*/
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
