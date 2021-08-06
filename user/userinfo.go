package user

import "github.com/mises-id/sdk/types"

var _ types.MUserInfo = &MisesUserInfoReadonly{}

type MisesUserInfo struct {
	Name        string   `json:"name,omitempty"`
	Gender      string   `json:"gender,omitempty"`
	AvatarId    string   `json:"avatar_did,omitempty"`
	AvatarThumb []byte   `json:"avatar_thumb,omitempty"`
	HomePage    string   `json:"home_page,omitempty"`
	Emails      []string `json:"emails,omitempty"`
	Telephones  []string `json:"telephones,omitempty"`
	Intro       string   `json:"into,omitempty"`
}
type MisesUserInfoReadonly struct {
	MisesUserInfo
}

func (user *MisesUserInfoReadonly) Name() string {
	return user.MisesUserInfo.Name
}
func (user *MisesUserInfoReadonly) Gender() string {
	return user.MisesUserInfo.Gender
}
func (user *MisesUserInfoReadonly) AvatarDid() string {
	return user.MisesUserInfo.AvatarId
}
func (user *MisesUserInfoReadonly) AvatarThumb() []byte {
	return user.MisesUserInfo.AvatarThumb
}
func (user *MisesUserInfoReadonly) HomePage() string {
	return user.MisesUserInfo.HomePage
}
func (user *MisesUserInfoReadonly) Emails() []string {
	return user.MisesUserInfo.Emails
}
func (user *MisesUserInfoReadonly) Telphones() []string {
	return user.MisesUserInfo.Telephones
}
func (user *MisesUserInfoReadonly) Intro() string {
	return user.MisesUserInfo.Intro
}

func NewMisesUserInfo(info types.MUserInfo) *MisesUserInfo {
	return &MisesUserInfo{
		Name:        info.Name(),
		Gender:      info.Gender(),
		AvatarId:    info.AvatarDid(),
		AvatarThumb: info.AvatarThumb(),
		HomePage:    info.HomePage(),
		Emails:      info.Emails(),
		Telephones:  info.Telphones(),
		Intro:       info.Intro(),
	}
}

func NewMisesUserInfoReadonly(
	name string,
	gender string,
	avatarDid string,
	avatarThumb []byte,
	homePage string,
	emails []string,
	telphones []string,
	intro string) types.MUserInfo {
	info := MisesUserInfo{
		Name:        name,
		Gender:      gender,
		AvatarId:    avatarDid,
		AvatarThumb: avatarThumb,
		HomePage:    homePage,
		Emails:      emails,
		Telephones:  telphones,
		Intro:       intro,
	}
	return &MisesUserInfoReadonly{info}
}
