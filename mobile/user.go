package mobile

import (
	"github.com/mises-id/sdk/misesid"
	"github.com/mises-id/sdk/types"
)

var _ MUserInfo = &mUserInfoWrapper{}
var _ MUser = &mUserWrapper{}
var _ MUserMgr = &mUserMgrWrapper{}

type mUserInfoWrapper struct {
	info types.MUserInfo
}
type mUserWrapper struct {
	types.MUser
}
type mUserMgrWrapper struct {
	types.MUserMgr
}

func (w *mUserInfoWrapper) Name() string {
	return w.info.Name()
}
func (w *mUserInfoWrapper) Gender() string {
	return w.info.Gender()
}
func (w *mUserInfoWrapper) Avatar() string {
	return w.info.Avatar()
}
func (w *mUserInfoWrapper) HomePage() string {
	return w.info.HomePage()
}
func (w *mUserInfoWrapper) Emails() MStringList {
	return &mStringListWrapper{w.info.Emails()}
}
func (w *mUserInfoWrapper) Telphones() MStringList {
	return &mStringListWrapper{w.info.Telphones()}
}
func (w *mUserInfoWrapper) Intro() string {
	return w.info.Intro()
}

func (w *mUserWrapper) MisesID() string {
	return w.MUser.MisesID()
}
func (w *mUserWrapper) Info() MUserInfo {
	return &mUserInfoWrapper{info: w.MUser.Info()}
}
func (w *mUserWrapper) SetInfo(info MUserInfo) (string, error) {
	minfo := misesid.NewMisesUserInfoReadonly(
		info.Name(),
		info.Gender(),
		info.Avatar(),
		info.HomePage(),
		mStringListToSlice(info.Emails()),
		mStringListToSlice(info.Telphones()),
		info.Intro(),
	)
	return w.MUser.SetInfo(minfo)
}
func (w *mUserWrapper) GetFollow(appDid string) MStringList {
	return &mStringListWrapper{w.MUser.GetFollow(appDid)}
}
func (w *mUserWrapper) SetFollow(followingID string, op bool, appDid string) (string, error) {
	return w.MUser.SetFollow(followingID, op, appDid)
}
func (w *mUserWrapper) LoadKeyStore(passPhrase string) error {
	return w.MUser.LoadKeyStore(passPhrase)
}
func (w *mUserWrapper) IsRegistered() error {
	return w.MUser.IsRegistered()
}

func (w *mUserMgrWrapper) CreateUser(mnemonic string, passPhrase string) (MUser, error) {
	u, err := w.MUserMgr.CreateUser(mnemonic, passPhrase)
	if err != nil {
		return nil, err
	}
	return &mUserWrapper{u}, nil
}
func (w *mUserMgrWrapper) ListUsers() MUserList {
	wus := []MUser{}
	for _, user := range w.MUserMgr.ListUsers() {
		wus = append(wus, &mUserWrapper{user})
	}

	return &mUserListWrapper{wus}
}
func (w *mUserMgrWrapper) SetActiveUser(userDid string, passPhrase string) error {
	users := w.MUserMgr.ListUsers()
	for _, cu := range users {
		if userDid == cu.MisesID() {
			err := cu.LoadKeyStore(passPhrase)
			if err != nil {
				return err
			}
			err = w.MUserMgr.SetActiveUser(userDid)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
func (w *mUserMgrWrapper) ActiveUser() MUser {
	u := w.MUserMgr.ActiveUser()
	if u == nil {
		return nil
	}
	return &mUserWrapper{u}
}

func NewMUserInfo(
	name string,
	gender string,
	avatar string,
	homePage string,
	emails MStringList,
	telphones MStringList,
	intro string) MUserInfo {
	info := misesid.NewMisesUserInfoReadonly(
		name,
		gender,
		avatar,
		homePage,
		mStringListToSlice(emails),
		mStringListToSlice(telphones),
		intro,
	)
	return &mUserInfoWrapper{info}
}
