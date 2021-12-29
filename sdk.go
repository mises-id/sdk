package sdk

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/mises-id/sdk/app"
	"github.com/mises-id/sdk/misesid"
	"github.com/mises-id/sdk/types"
	"github.com/mises-id/sdk/user"
)

var _ types.MSdk = &misesSdk{}

type MSdkOption struct {
	ChainID    string //'mises' for the mainnet
	PassPhrase string //8 chars needed, default is 'mises.site'
	Debug      bool
}

type misesSdk struct {
	options MSdkOption
	userMgr types.MUserMgr
	app     types.MApp
}

func (ctx *misesSdk) setupLogger() {
}

func NewSdkForUser(options MSdkOption, passPhrase string) types.MSdk {
	if options.ChainID == "" {
		options.ChainID = types.DefaultChainID
	}

	var ctx misesSdk
	ctx.options = options
	ctx.userMgr = MSdkInitUserMgr(passPhrase)

	return &ctx
}

func NewSdkForApp(options MSdkOption, info types.MAppInfo) (types.MSdk, types.MApp) {
	if options.ChainID == "" {
		options.ChainID = types.DefaultChainID
	}
	if options.PassPhrase == "" {
		options.PassPhrase = types.DefaultPassPhrase
	}

	var ctx misesSdk
	ctx.options = options

	mapp, err := ctx.EnsureApp(info)
	if err != nil {
		panic(err)
	}
	ctx.app = mapp
	return &ctx, mapp
}

func MSdkInitUserMgr(passPhrase string) types.MUserMgr {
	var userMgr user.MisesUserMgr

	var u user.MisesUser
	err := u.LoadKeyStore(passPhrase)
	if u.MisesID() != "" {
		userMgr.AddUser(&u)
	}
	if err == nil {
		userMgr.SetActiveUser(u.MisesID())
	}

	return &userMgr
}

func (sdk *misesSdk) EnsureApp(info types.MAppInfo) (types.MApp, error) {

	var app app.MisesApp
	if err := app.Init(info, sdk.options.ChainID, sdk.options.PassPhrase); err != nil {
		return nil, err
	}
	return &app, nil
}

func (sdk *misesSdk) SetEndpoint(endpoint string) error {
	return misesid.SetTestEndpoint(endpoint)
}

func (sdk *misesSdk) Login(site string, permission []string) (string, error) {
	var valid bool = false
	for _, domain := range sdk.app.Info().Domains() {
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
	nonce := strconv.FormatInt(time.Now().UTC().Unix(), 10)
	sigData := url.Values{}
	sigData.Add("mises_id", auser.MisesID())
	sigData.Add("nonce", nonce)
	signed, err := auser.Signer().Sign(sigData.Encode())
	if err != nil {
		return "", err
	}
	v := url.Values{}
	v.Add("mises_id", auser.MisesID())
	v.Add("nonce", nonce)
	v.Add("sig", signed)

	return v.Encode(), nil
}
func (sdk *misesSdk) VerifyLogin(auth string) (string, string, error) {

	v, err := url.ParseQuery(auth)
	if err != nil {
		return "", "", err
	}
	misesID := v.Get("mises_id")
	sigStr := v.Get("sig")
	nonce := v.Get("nonce")
	pubKeyStr := v.Get("pubkey")

	if err := misesid.CheckMisesID(misesID, pubKeyStr); err != nil {
		return "", "", err
	}

	if err := misesid.Verify("mises_id="+misesID+"&nonce="+nonce, pubKeyStr, sigStr); err != nil {
		return "", "", err
	}
	return misesID, pubKeyStr, nil
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
