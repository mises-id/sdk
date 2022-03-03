package misesid

import "github.com/mises-id/sdk/types"

var _ types.MUserInfo = &MisesUserInfoReadonly{}

type MisesUserInfo struct {
	Name        string   `json:"name,omitempty"`
	Gender      string   `json:"gender,omitempty"`
	AvatarUrl   string   `json:"avatar_url,omitempty"`
	HomePageUrl string   `json:"home_page_url,omitempty"`
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
func (user *MisesUserInfoReadonly) Avatar() string {
	return user.MisesUserInfo.AvatarUrl
}
func (user *MisesUserInfoReadonly) HomePage() string {
	return user.MisesUserInfo.HomePageUrl
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
		AvatarUrl:   info.Avatar(),
		HomePageUrl: info.HomePage(),
		Emails:      info.Emails(),
		Telephones:  info.Telphones(),
		Intro:       info.Intro(),
	}
}

func NewMisesUserInfoReadonly(
	name string,
	gender string,
	avatar string,
	homePage string,
	emails []string,
	telphones []string,
	intro string) types.MUserInfo {
	info := MisesUserInfo{
		Name:        name,
		Gender:      gender,
		AvatarUrl:   avatar,
		HomePageUrl: homePage,
		Emails:      emails,
		Telephones:  telphones,
		Intro:       intro,
	}
	return &MisesUserInfoReadonly{info}
}
