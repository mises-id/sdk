package user

import (
	"testing"

	"github.com/mises-id/sdk/misesid"
	"github.com/tyler-smith/assert"
)

func TestLoadKeyStoreFile(t *testing.T) {
	var user MisesUser
	var passPhrase = "123456"

	pu := &user

	err := pu.LoadKeyStore(passPhrase)
	assert.NoError(t, err)
	assert.EqualString(t, misesid.Ks.PubKey, user.pubKey)
}
