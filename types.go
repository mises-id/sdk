package sdk

type MSdk interface {
	AppMgr() (MAppMgr, error)
	UserMgr() (MUserMgr, error)
	TestConnection() error
	SetLogLevel(level int) error

	RandomSeed() (MKeySeed, error)
	RestoreSeed(mnemonic string, passPhrase string) (MKeySeed, error)
}

type MKeySeed interface {
	Mnemonic() string
	PassPhrase() string
}

type MUserInfo interface {
	Name() string
	Gender() string
	AvatarDid() string    //did of avatar file did:mises:0123456789abcdef/avatar
	AavatarThumb() []byte //avatar thumb is a bitmap
	HomePage() string     //url
	Emails() []string
	Telphones() []string
	Intro() string
}
type MUserAuthorization interface {
	UserDid() string       //mises app 的did
	AppDid() string        //to
	Permissions() []string //user_info_r,  user_info_w, user_relation_r, user_relation_w
	ExpireTimestamp() int  //
	AppAuthorization() MAppAuthorization
}

type MUserMgr interface {
	AddUser(seed MKeySeed) (MUser, error)
	ListUsers() ([]MUser, error)
	SetActiveUser(userDid string) error
	ActiveUser() (MUser, error)
	RemoveUser(userDid string, seed MKeySeed) error
}

type MUser interface {
	MisesID() (string, error)
	Info(appDid string) (MUserInfo, error)
	IsRegistered() (bool, error)
	Register(info MUserInfo, appDid string) error
	SetInfo(info MUserInfo, appDid string) error
	Follow(whomDid string, appDid string) error
	UnFollow(whomDid string, appDid string) error
	AddToBlackList(whomDid string, appDid string) error
	RemoveFromBlockList(whomDid string, appDid string) error
	ApproveAuthorization(appDid string, permissions []string, expireIn int) (MUserAuthorization, error)
	RevokeAuthorization(appDid string, permissions []string) (MUserAuthorization, error)
	GetAuthorization(appDid string) (MUserAuthorization, error)
}

type MAppInfo interface {
	ApppDid() string //did:mises:0123456789abcdef
	AppName() string //
	IconDid() string //udid of icon file did:mises:0123456789abcdef/icon
	IconThumb() []byte
	Domain() string //app
	Developer() string
}

type MAppAuthorization interface {
	AppInfo() MAppInfo
	Permissions() []string //
	ExpireTimestamp() int  //
}

type MAppMgr interface {
	AddApp(app MAppAuthorization, removable bool) (MApp, error)
	ListApps() ([]MApp, error)
	RemoveApp(appDid string) error
}

type MEvent interface {
	EventID() int
}
type MEventCallback func(event MEvent) error

type MApp interface {
	AppDID() (string, error)
	IsRegistered() (bool, error)
	Register(MAppInfo string, appDid string) error
	GenerateAuthorization(permisions []string) (MAppAuthorization, error)

	AddEventListener(userDid string, userAuth MUserAuthorization, callback MEventCallback)
	RemoveUserEventListener(userDid string, userAuth MUserAuthorization, callback MEventCallback) error
	ListFollow(whomDid string) (dids []string)
	ListBlackListed(whomDid string) (dids []string)
	Commit() error
	Cancel() error
	Agent() (MAgent, error)
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
