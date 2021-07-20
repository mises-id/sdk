package sdk

import (
	"github.com/mises-id/sdk/user"
)

const DefaultEndpoint string = "http://localhost:1317"
const DefaultChainID string = "mises"
const AddressPrefix string = "mises"

const ErrorInvalidLeaseTime string = "Invalid lease time"
const ErrorKeyIsRequired string = "Key is required"
const ErrorValueIsRequired string = "Value is required"
const ErrorKeyFormat string = "Key format error"

type MSdk interface {
	UserMgr() user.MUserMgr
	TestConnection() error
	SetLogLevel(level int) error
	Login(site string, permissions []string) (string, error)
}

/*
type MUserAuthorization interface {
	UserDid() string       //mises app 的did
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
