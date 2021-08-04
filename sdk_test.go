package sdk

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/mises-id/sdk/misesid"
	"github.com/mises-id/sdk/user"
	"github.com/tyler-smith/assert"
)

func TestNewSdkForUesr(t *testing.T) {
	mo := MSdkOption{"test", true}

	// test NewSdkForUser
	s := NewSdkForUser(mo, "123456")

	ugr := s.UserMgr()

	u := ugr.ActiveUser()
	if u != nil {
		fmt.Printf("user's mid is %s\n", u.MisesID())
		fmt.Printf("user's privKey is %s\n", u.PrivKEY())
		fmt.Printf("user's pubKey is %s\n", u.PubKEY())

		// test Login, sign & verify
		permissions := []string{"user_info_r", "user_info_w"}
		auth, err := s.Login("mises.site", permissions)
		assert.NoError(t, err)

		fmt.Printf("auth string is: %s\n", auth)

		v, err := url.ParseQuery(auth)
		assert.NoError(t, err)
		misesID := v.Get("mises_id")
		sigStr := v.Get("sig")
		nonce := v.Get("nonce")

		err = user.Verify(misesID+"&"+nonce, u.PubKEY(), sigStr)
		assert.NoError(t, err)
		if err == nil {
			fmt.Printf("signature is verified\n")
		} else {
			fmt.Printf("Signature verification is failed\n")
		}
	}

	fmt.Printf("keystore version is: %d\n", misesid.Ks.Version)
	fmt.Printf("keystore id is: %s\n", misesid.Ks.MId)
	fmt.Printf("keystore address is: %s\n", misesid.Ks.PubKey)
	fmt.Printf("keystore kdf is: %s\n", misesid.Ks.Crypto.Kdf)
	fmt.Printf("keystore dklen is: %d\n", misesid.Ks.Crypto.KdfParams.Dklen)
	fmt.Printf("keystore salt is: %s\n", misesid.Ks.Crypto.KdfParams.Salt)
	fmt.Printf("keystore n is: %d\n", misesid.Ks.Crypto.KdfParams.N)
	fmt.Printf("keystore r is: %d\n", misesid.Ks.Crypto.KdfParams.R)
	fmt.Printf("keystore p is: %d\n", misesid.Ks.Crypto.KdfParams.P)
	fmt.Printf("keystore cipher is: %s\n", misesid.Ks.Crypto.Cipher)
	fmt.Printf("keystore ciphertext is: %s\n", misesid.Ks.Crypto.Ciphertext)
	fmt.Printf("keystore iv is: %s\n", misesid.Ks.Crypto.CipherParams.Iv)
	fmt.Printf("keystore mac is: %s\n", misesid.Ks.Crypto.Mac)

	// test CreateUser
	mnemonics, err := RandomMnemonics()
	assert.NoError(t, err)
	fmt.Printf("mnemonics is: %s\n", mnemonics)

	user, err := ugr.CreateUser(mnemonics, "123456")
	assert.NoError(t, err)

	user.MisesID()
}
