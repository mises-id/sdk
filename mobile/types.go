package mobile

type MStringList interface {
	Count() int
	Get(idx int) string
}
type MUserList interface {
	Count() int
	Get(idx int) MUser
}

type MUserInfo interface {
	Name() string
	Gender() string
	AvatarDid() string    //did of avatar file did:mises:0123456789abcdef/avatar
	AavatarThumb() []byte //avatar thumb is a bitmap
	HomePage() string     //url
	Emails() MStringList
	Telphones() MStringList
	Intro() string
}

type MUser interface {
	MisesID() string
	PubKEY() string
	PrivKEY() string
	Info() MUserInfo
	SetInfo(info MUserInfo) (string, error)
	GetFollow(appDid string) MStringList
	SetFollow(followingDid string, op bool, appDid string) (string, error)
	LoadKeyStore(passPhrase string) error
	IsRegistered() error
	Register(appDid string) (string, error)
}

type MUserMgr interface {
	CreateUser(mnemonic string, passPhrase string) (MUser, error)
	ListUsers() MUserList
	SetActiveUser(userDid string, passPhrase string) error
	ActiveUser() MUser
}

type MSdk interface {
	UserMgr() MUserMgr
	SetTestEndpoint(endpoint string) error
	TestConnection() error
	SetLogLevel(level int) error
	SetHomePath(dir string) error
	Login(site string, permissions MStringList) (string, error)
	RandomMnemonics() (string, error)

	CheckSession(sessinID string) (bool, error)
}
