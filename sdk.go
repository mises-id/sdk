package sdk

import (
	"fmt"
	"net/url"

	"github.com/mises-id/sdk/app"
	"github.com/mises-id/sdk/bip39"
	"github.com/mises-id/sdk/types"
	"github.com/mises-id/sdk/user"
)

var _ types.MSdk = &misesSdk{}

type MSdkOption struct {
	ChainID string
	Debug   bool
}

type misesSdk struct {
	options MSdkOption
	userMgr types.MUserMgr
	app     app.MApp
}

func (ctx *misesSdk) setupLogger() {
}

func NewSdkForUser(options MSdkOption, passPhrase string) types.MSdk {
	if options.ChainID == "" {
		options.ChainID = types.DefaultChainID
	}

	var ctx misesSdk
	ctx.options = options
	ctx.userMgr, ctx.app = MSdkInit(passPhrase)

	return &ctx
}

func NewSdkForApp(options MSdkOption) types.MSdk {
	if options.ChainID == "" {
		options.ChainID = types.DefaultChainID
	}

	var ctx misesSdk
	ctx.options = options
	ctx.userMgr, ctx.app = MSdkInit("")

	if ctx.userMgr.ActiveUser() == nil {
		mnemonics, err := RandomMnemonics()
		if err != nil {
			panic(err)
		}
		admin, err := ctx.userMgr.CreateUser(mnemonics, "")
		admin.Register("appID")
		ctx.userMgr.SetActiveUser(admin.MisesID())
	}

	return &ctx
}

func MSdkInit(passPhrase string) (types.MUserMgr, app.MApp) {
	var userMgr user.MisesUserMgr
	var a app.MisesApp
	var u user.MisesUser

	a.SetAppDomain(app.MisesDiscover)

	err := u.LoadKeyStore(passPhrase)
	if u.MisesID() != "" {
		userMgr.AddUser(&u)
	}
	if err == nil {
		userMgr.SetActiveUser(u.MisesID())
	}

	return &userMgr, &a
}

func RandomMnemonics() (string, error) {
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		return "", err
	}

	mnemonics, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", err
	}

	return mnemonics, nil
}

func (sdk *misesSdk) Login(site string, permission []string) (string, error) {
	var valid bool = false
	for _, domain := range sdk.app.AppDomains() {
		if site == domain {
			valid = true
			break
		}
	}
	if !valid {
		return "", fmt.Errorf("only mises discover supported")
	}
	auser := sdk.userMgr.ActiveUser()
	if auser == nil {
		return "", fmt.Errorf("no active user")
	}

	sdk.app.AddAuth(auser.MisesID(), permission)

	// sign user's misesid, publicKey using his privateKey, return the signed result
	signed, nonce, err := user.Sign(auser, auser.MisesID())
	if err != nil {
		return "", err
	}
	v := url.Values{}
	v.Add("mises_id", auser.MisesID())
	v.Add("nonce", nonce)
	v.Add("sig", signed)

	return v.Encode(), nil
}
func (sdk *misesSdk) VerifyLogin(auth string) (string, error) {
	auser := sdk.userMgr.ActiveUser()
	if auser == nil {
		return "", fmt.Errorf("no active user")
	}
	v, err := url.ParseQuery(auth)
	if err != nil {
		return "", err
	}
	misesID := v.Get("mises_id")
	sigStr := v.Get("sig")
	nonce := v.Get("nonce")
	pubKeyStr, err := user.GetUser(auser, misesID)
	if err == nil {
		return "", err
	}

	err = user.Verify(misesID+"&"+nonce, pubKeyStr, sigStr)
	if err == nil {
		return "", err
	}
	return misesID, nil
}

func (sdk *misesSdk) UserMgr() types.MUserMgr {
	return sdk.userMgr
}

func (w *misesSdk) TestConnection() error {
	return nil
}

func (w *misesSdk) SetLogLevel(level int) error {
	return nil
}
