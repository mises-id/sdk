package types

import (
	"os"
	"path/filepath"

	"github.com/btcsuite/btcd/btcec"
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
	AddressPrefix         string = "mises"
	MisesIDPrefix                = "did:mises:"
	ErrorInvalidLeaseTime string = "Invalid lease time"
	ErrorKeyIsRequired    string = "Key is required"
	ErrorValueIsRequired  string = "Value is required"
	ErrorKeyFormat        string = "Key format error"
)

type MSdk interface {
	UserMgr() MUserMgr
	TestConnection() error
	SetLogLevel(level int) error
	Login(site string, permissions []string) (string, error)
	VerifyLogin(auth string) (string, error)
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
	PubKEY() string
	PrivKEY() string
	PrivateKey() *btcec.PrivateKey
	PublicKey() *btcec.PublicKey
	Info() MUserInfo
	GetFollow(appDid string) []string
	LoadKeyStore(passPhrase string) error
	IsRegistered() error

	SetInfo(info MUserInfo) (string, error)
	SetFollow(followingId string, op bool, appDid string) (string, error)
	Register(appDid string) (string, error)
}

type MUserMgr interface {
	CreateUser(mnemonic string, passPhrase string) (MUser, error)
	ListUsers() []MUser
	AddUser(user MUser)
	SetActiveUser(userDid string) error
	ActiveUser() MUser
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
