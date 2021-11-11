package user

import (
	"testing"

	"github.com/mises-id/sdk/bip39"
	"github.com/mises-id/sdk/misesid"
	"github.com/tyler-smith/assert"
)

func CreateRandomUser(t *testing.T) (*MisesUserMgr, *MisesUser) {
	misesid.DeleteKeyStoreFile()
	entropy, err := bip39.NewEntropy(128)
	assert.NoError(t, err)
	mnemonics, err := bip39.NewMnemonic(entropy)
	assert.NoError(t, err)

	var ugr MisesUserMgr
	pUgr := &ugr
	passwd := "123456"
	user, err := pUgr.CreateUser(mnemonics, passwd)
	assert.NoError(t, err)
	muser, _ := user.(*MisesUser)
	return pUgr, muser

}
func TestUserMgrCreate(t *testing.T) {
	_, user := CreateRandomUser(t)

	privKeyCreated := user.PrivKEY()

	err := user.LoadKeyStore("123456")
	assert.NoError(t, err)

	assert.EqualString(t, privKeyCreated, user.PrivKEY())
}

func TestUserMgrList(t *testing.T) {
	misesid.DeleteKeyStoreFile()
	CreateRandomUser(t)
	var ugr MisesUserMgr
	ul := ugr.ListUsers()
	assert.EqualInt(t, 0, len(ul))

	var u MisesUser
	err := u.LoadKeyStore("")
	assert.NotNil(t, err)
	ugr.AddUser(&u)
	ul = ugr.ListUsers()
	assert.EqualInt(t, 1, len(ul))

	actuser := ugr.ActiveUser()
	assert.Nil(t, actuser)
	u.LoadKeyStore("123456")
	ugr.SetActiveUser(u.MisesID())
	actuser = ugr.ActiveUser()
	assert.NotNil(t, actuser)
}

func TestUserMgrSetActivate(t *testing.T) {
	_, user := CreateRandomUser(t)

	err := user.LoadKeyStore("1")
	assert.EqualString(t, err.Error(), ErrorMsgWrongPassword)

	err = user.LoadKeyStore("12")
	assert.EqualString(t, err.Error(), ErrorMsgWrongPassword)
	err = user.LoadKeyStore("123")
	assert.EqualString(t, err.Error(), ErrorMsgWrongPassword)
	err = user.LoadKeyStore("1234")
	assert.EqualString(t, err.Error(), ErrorMsgWrongPassword)
	err = user.LoadKeyStore("12345")
	assert.EqualString(t, err.Error(), ErrorMsgWrongPassword)
	err = user.LoadKeyStore("123456")
	assert.NoError(t, err)
}
