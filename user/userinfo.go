package user

import "github.com/mises-id/sdk/types"

var _ types.MUserInfo = &MisesUserInfo{}

type MisesUserInfo struct {
	name        string
	gender      string
	avatarId    string
	avatarThumb []byte
	homePage    string
	emails      []string
	telephones  []string
	intro       string
}

func (info *MisesUserInfo) Name() string {
	return info.name
}

func (info *MisesUserInfo) Gender() string {
	return info.gender
}
func (info *MisesUserInfo) AvatarDid() string {
	return info.avatarId
} //did of avatar file did:mises:0123456789abcdef/avatar
func (info *MisesUserInfo) AvatarThumb() []byte {
	return info.avatarThumb
} //avatar thumb is a bitmap
func (info *MisesUserInfo) HomePage() string {
	return info.homePage
} //url
func (info *MisesUserInfo) Emails() []string {
	return info.emails
}
func (info *MisesUserInfo) Telphones() []string {
	return info.telephones
}
func (info *MisesUserInfo) Intro() string {
	return info.intro
}

func NewMisesUserInfo(info types.MUserInfo) *MisesUserInfo {
	return &MisesUserInfo{
		name:        info.Name(),
		gender:      info.Gender(),
		avatarId:    info.AvatarDid(),
		avatarThumb: info.AvatarThumb(),
		homePage:    info.HomePage(),
		emails:      info.Emails(),
		telephones:  info.Telphones(),
		intro:       info.Intro(),
	}
}

func NewMisesUserInfoRaw(
	name string,
	gender string,
	avatarDid string,
	avatarThumb []byte,
	homePage string,
	emails []string,
	telphones []string,
	intro string) *MisesUserInfo {
	return &MisesUserInfo{
		name:        name,
		gender:      gender,
		avatarId:    avatarDid,
		avatarThumb: avatarThumb,
		homePage:    homePage,
		emails:      emails,
		telephones:  telphones,
		intro:       intro,
	}
}
