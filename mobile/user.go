package mobile

import (
	"github.com/mises-id/sdk/user"
)

type mUserInfoWrapper struct {
	info user.MisesUserInfo
}
type mUserWrapper struct {
	user.MUser
}
type mUserMgrWrapper struct {
	user.MUserMgr
}

func (w *mUserInfoWrapper) Name() string {
	return w.info.Name
}
func (w *mUserInfoWrapper) Gender() string {
	return w.info.Gender
}
func (w *mUserInfoWrapper) AvatarDid() string {
	return w.info.AvatarId
}
func (w *mUserInfoWrapper) AavatarThumb() []byte {
	return w.info.AvatarThumb
}
func (w *mUserInfoWrapper) HomePage() string {
	return w.info.HomePage
}
func (w *mUserInfoWrapper) Emails() MStringList {
	return &mStringListWrapper{w.info.Emails}
}
func (w *mUserInfoWrapper) Telphones() MStringList {
	return &mStringListWrapper{w.info.Telephones}
}
func (w *mUserInfoWrapper) Intro() string {
	return w.info.Intro
}

func (w *mUserWrapper) MisesID() string {
	return w.MUser.MisesID()
}
func (w *mUserWrapper) PubKEY() string {
	return w.MUser.PubKEY()
}
func (w *mUserWrapper) PrivKEY() string {
	return w.MUser.PrivKEY()
}
func (w *mUserWrapper) Info() MUserInfo {
	return &mUserInfoWrapper{info: w.MUser.Info()}
}
func (w *mUserWrapper) SetInfo(info MUserInfo) string {
	var i user.MisesUserInfo
	return w.MUser.SetInfo(i)
}
func (w *mUserWrapper) GetFollow(appDid string) MStringList {
	return &mStringListWrapper{w.MUser.GetFollow(appDid)}
}
func (w *mUserWrapper) SetFollow(followingId string, op bool, appDid string) string {
	return w.MUser.SetFollow(followingId, op, appDid)
}
func (w *mUserWrapper) LoadKeyStore(passPhrase string) error {
	return w.MUser.LoadKeyStore(passPhrase)
}
func (w *mUserWrapper) IsRegistered() (bool, error) {
	return w.MUser.IsRegistered()
}
func (w *mUserWrapper) Register(info MUserInfo, appDid string) error {
	var i user.MisesUserInfo
	return w.MUser.Register(i, appDid)
}

func (w *mUserMgrWrapper) CreateUser(mnemonic string, passPhrase string) (MUser, error) {
	u, err := w.MUserMgr.CreateUser(mnemonic, passPhrase)
	if err != nil {
		return nil, err
	}
	w.MUserMgr.AddUser(u)
	return &mUserWrapper{u}, nil
}
func (w *mUserMgrWrapper) ListUsers() MUserList {
	wus := []MUser{}
	for _, user := range w.MUserMgr.ListUsers() {
		wus = append(wus, &mUserWrapper{user})
	}

	return &mUserListWrapper{wus}
}
func (w *mUserMgrWrapper) SetActiveUser(userDid string) error {
	return w.MUserMgr.SetActiveUser(userDid)
}
func (w *mUserMgrWrapper) ActiveUser() MUser {
	u := w.MUserMgr.ActiveUser()
	if u == nil {
		return nil
	}
	return &mUserWrapper{u}
}
