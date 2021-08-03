package mobile

import (
	"os"
	"github.com/mises-id/sdk"
)

var _ MSdk = &mSdkWrapper{}

type mSdkWrapper struct {
	sdk.MSdk
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
	return sdk.RandomMnemonics()
}
func (w *mSdkWrapper) SetHomePath(dir string) error {
	err := os.Chdir(dir)
	if err != nil {
			panic(err)
	}
	return nil
}
func NewMSdk() MSdk {
	opt := sdk.MSdkOption{}
	ret := sdk.NewSdkForUser(opt, "")
	return &mSdkWrapper{ret}
}
