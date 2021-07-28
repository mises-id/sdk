package sdk

import (
	"fmt"

	"github.com/mises-id/sdk/app"
	"github.com/mises-id/sdk/bip39"
	"github.com/mises-id/sdk/user"
	"github.com/tendermint/tendermint/libs/log"
)

type MSdkOption struct {
	ChainID string
	Debug   bool
}

type misesSdk struct {
	MSdk
	options MSdkOption
	userMgr user.MUserMgr
	app     app.MApp
	logger  log.Logger
}

func (ctx *misesSdk) setupLogger() {
	ctx.logger = log.NewNopLogger()
}

func NewSdkForUser(options MSdkOption, passPhrase string) MSdk {
	if options.ChainID == "" {
		options.ChainID = DefaultChainID
	}

	var ctx misesSdk
	ctx.options = options
	ctx.setupLogger()
	ctx.userMgr, ctx.app = MSdkInit(passPhrase)

	return &ctx
}

func MSdkInit(passPhrase string) (user.MUserMgr, app.MApp) {
	var userMgr user.MisesUserMgr
	var a app.MisesApp
	var u user.MisesUser

	a.SetAppDomain(app.MisesDiscover)

	err := u.LoadKeyStore(passPhrase)
	if err != nil {
		return &userMgr, &a
	}

	userMgr.AddUser(&u)
	userMgr.SetActiveUser(u.MisesID())

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
	if site != sdk.app.AppDomain() {
		return "", fmt.Errorf("only mises discover supported")
	}

	sdk.app.AddAuth(sdk.userMgr.ActiveUser().MisesID(), permission)

	// sign user's misesid, publicKey using his privateKey, return the signed result
	_, signed, err := user.Sign(sdk.userMgr.ActiveUser(), sdk.userMgr.ActiveUser().MisesID())
	if err != nil {
		return "", err
	}

	return signed, nil
}

func (sdk *misesSdk) UserMgr() user.MUserMgr {
	return sdk.userMgr
}
