/*
	App is a Community Site which supports misesid,
	there is only one app "mises.site" during the first stage(MVP)
*/
package app

type MisesApp struct {
	MApp
	appDId    string
	appDomain string
	auths     []MisesAuth
}

type MisesAuth struct {
	UId                 string
	ExpirationInSeconds int
	Permissions         []string
}

const (
	MisesDiscover = "mises.site"
	defaultExpirationInSeconds = 120
)

type MApp interface {
	AppDID() string
	AppDomain() string
	SetAppDomain(string)

	IsRegistered() bool
	Register(MAppInfo string, appDid string) error

	AddAuth(misesId string, permissions []string)
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
	return app.appDId
}

func (app *MisesApp) AppDomain() string {
	return app.appDomain
}

func (app *MisesApp) SetAppDomain(domain string) {
	app.appDomain = domain
}

func (app *MisesApp) AddAuth(misesId string, permissions []string) {
	var auth MisesAuth
	auth.UId = misesId
	auth.ExpirationInSeconds = defaultExpirationInSeconds // default is 120 seconds
	auth.Permissions = permissions

	app.auths = append(app.auths, auth)
}
