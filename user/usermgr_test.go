package user

import (
	"testing"

	"github.com/mises-id/sdk/bip39"
	"github.com/tyler-smith/assert"
)

func TestCeateUser(t *testing.T) {
	entropy, err := bip39.NewEntropy(128)
	assert.NoError(t, err)
	mnemonics, err := bip39.NewMnemonic(entropy)
	assert.NoError(t, err)

	var ugr MisesUserMgr
	pUgr := &ugr
	passwd := "123456"
	user, err := pUgr.CreateUser(mnemonics, passwd)
	assert.NoError(t, err)

	privKeyCreated := user.PrivKEY()

	err = user.LoadKeyStore(passwd)
	assert.NoError(t, err)

	assert.EqualString(t, privKeyCreated, user.PrivKEY())
}
