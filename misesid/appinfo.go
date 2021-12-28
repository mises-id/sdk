package misesid

import "github.com/mises-id/sdk/types"

var _ types.MAppInfo = &MisesAppInfoReadonly{}

type MisesAppInfo struct {
	Name        string   `json:"name,omitempty"`
	IconUrl     string   `json:"icon_url,omitempty"`
	HomePageUrl string   `json:"home_page_url,omitempty"`
	Domains     []string `json:"domains,omitempty"`

	Developer string `json:"developer,omitempty"`
}
type MisesAppInfoReadonly struct {
	MisesAppInfo
}

func (user *MisesAppInfoReadonly) AppName() string {
	return user.MisesAppInfo.Name
}
func (user *MisesAppInfoReadonly) IconURL() string {
	return user.MisesAppInfo.IconUrl
}
func (user *MisesAppInfoReadonly) HomeURL() string {
	return user.MisesAppInfo.HomePageUrl
}
func (user *MisesAppInfoReadonly) Domains() []string {
	return user.MisesAppInfo.Domains
}
func (user *MisesAppInfoReadonly) Developer() string {
	return user.MisesAppInfo.Developer
}

func NewMisesAppInfo(info types.MAppInfo) *MisesAppInfo {
	return &MisesAppInfo{
		Name:        info.AppName(),
		IconUrl:     info.IconURL(),
		HomePageUrl: info.HomeURL(),
		Domains:     info.Domains(),
		Developer:   info.Developer(),
	}
}

func NewMisesAppInfoReadonly(
	name string,
	iconUrl string,
	homeUrl string,
	domains []string,
	developer string) types.MAppInfo {
	info := MisesAppInfo{
		Name:        name,
		IconUrl:     iconUrl,
		HomePageUrl: homeUrl,
		Domains:     domains,
		Developer:   developer,
	}
	return &MisesAppInfoReadonly{info}
}
