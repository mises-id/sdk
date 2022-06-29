package mobile

import (
	"os"
	"time"

	"github.com/mises-id/sdk"
	"github.com/mises-id/sdk/bip39"
	"github.com/mises-id/sdk/misesid"
	"github.com/mises-id/sdk/types"
)

var _ MSdk = &mSdkWrapper{}

type mSdkWrapper struct {
	types.MSdk
}

func (w *mSdkWrapper) UserMgr() MUserMgr {
	return &mUserMgrWrapper{w.MSdk.UserMgr()}
}
func (w *mSdkWrapper) TestConnection() error {
	return w.MSdk.TestConnection()
}
func (w *mSdkWrapper) SetLogLevel(level int) error {
	return w.MSdk.SetLogLevel(level)
}
func (w *mSdkWrapper) Login(site string, permissions MStringList) (string, error) {
	return w.MSdk.Login(site, mStringListToSlice(permissions))
}
func (w *mSdkWrapper) RandomMnemonics() (string, error) {
	return misesid.RandomMnemonics()
}
func (w *mSdkWrapper) CheckMnemonics(mne string) error {
	_, err := bip39.Mnemonic2ByteArray(mne)
	return err
}

type mSessionResultWrapper struct {
	misesid.WaitResult
}

func (w *mSessionResultWrapper) SessionID() string {
	return w.Session
}
func (w *mSessionResultWrapper) Msg() string {
	if w.ErrMsg != "" {
		return w.ErrMsg
	}
	return w.Result
}
func (w *mSessionResultWrapper) Success() bool {
	return w.ErrMsg == ""
}

func (w *mSdkWrapper) PollSessionResult() MSessionResult {
	wr, err := misesid.PollSessionResult(2 * time.Second)
	if err != nil {
		return nil
	}
	return &mSessionResultWrapper{*wr}
}
func (w *mSdkWrapper) SetTestEndpoint(endpoint string) error {
	return misesid.SetTestEndpoint(endpoint)
}

func NewMSdk() MSdk {
	opt := types.MSdkOption{}
	opt.PassPhrase = ""
	ret := sdk.NewSdkForUser(opt)
	return &mSdkWrapper{ret}
}

func SetHomePath(dir string) error {
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}

	types.NodeHome = dir + ".misestm"
	return nil
}
