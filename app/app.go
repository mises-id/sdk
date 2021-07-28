/*
	App is a Community Site which supports misesid,
	there is only one app "mises.site" during the first stage(MVP)
*/
package app

type MisesApp struct {
	MApp
	appDID    string
	appDomain string
	auth      []MisesAuth
}

type MisesAuth struct {
	Uid        string
	ExpireTime int
	Permission []string
}

var MisesDiscover = "mises.site"
var expireTime = 120

type MApp interface {
	AppDID() string
	AppDomain() string
	SetAppDomain(string)

	IsRegistered() bool
	Register(MAppInfo string, appDid string) error

	AddAuth(misesid string, permission []string)
}

type MAppInfo interface {
	ApppDid() string //did:mises:0123456789abcdef
	AppName() string //
	IconDid() string //udid of icon file did:mises:0123456789abcdef/icon
	IconThumb() []byte
	Domain() string //app
	Developer() string
}

func (app *MisesApp) AppDID() string {
	return app.appDID
}

func (app *MisesApp) AppDomain() string {
	return app.appDomain
}

func (app *MisesApp) SetAppDomain(domain string) {
	app.appDomain = domain
}

func (app *MisesApp) AddAuth(misesid string, permission []string) {
	var auth MisesAuth
	auth.Uid = misesid
	auth.ExpireTime = expireTime // default is 120 seconds
	auth.Permission = permission

	app.auth = append(app.auth, auth)
}
