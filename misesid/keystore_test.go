package misesid

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	"github.com/tyler-smith/assert"
)

func TestReadKeyStoreFile(t *testing.T) {
	DeleteKeyStoreFile()
	NewKeyStore(t)
	ks, err := ReadKeyStoreFile()
	if err == nil {
		fmt.Printf("keystore version is: %d\n", ks.Version)
		fmt.Printf("keystore id is: %s\n", ks.MId)
		fmt.Printf("keystore address is: %s\n", ks.PubKey)
		fmt.Printf("keystore kdf is: %s\n", ks.Crypto.Kdf)
		fmt.Printf("keystore dklen is: %d\n", ks.Crypto.KdfParams.Dklen)
		fmt.Printf("keystore salt is: %s\n", ks.Crypto.KdfParams.Salt)
		fmt.Printf("keystore n is: %d\n", ks.Crypto.KdfParams.N)
		fmt.Printf("keystore r is: %d\n", ks.Crypto.KdfParams.R)
		fmt.Printf("keystore p is: %d\n", ks.Crypto.KdfParams.P)
		fmt.Printf("keystore cipher is: %s\n", ks.Crypto.Cipher)
		fmt.Printf("keystore ciphertext is: %s\n", ks.Crypto.Ciphertext)
		fmt.Printf("keystore iv is: %s\n", ks.Crypto.CipherParams.Iv)
		fmt.Printf("keystore mac is: %s\n", ks.Crypto.Mac)
	}

}

func NewKeyStore(t *testing.T) *KeyStore {
	ks := &KeyStore{}
	ks.InitKdfParam()
	s, err := ks.Scrypt("123456")
	assert.NoError(t, err)
	privateKey := "d88169685ceeaaaef5c619bcdaca3b44c323372e0d777720d4a28efe3c658aa2"
	privKey := []byte(privateKey)
	ciphertext, iv, err := AesEncrypt(privKey, s)
	assert.NoError(t, err)
	mac, err := GenMac(s, ciphertext)
	assert.NoError(t, err)
	ks.MId = ""
	ks.PubKey = ""
	ks.Version = Ver
	ks.Crypto.Ciphertext = hex.EncodeToString(ciphertext)
	ks.Crypto.Cipher = CipherMethod
	ks.Crypto.Kdf = KdfMethod
	ks.Crypto.Mac = hex.EncodeToString(mac)
	ks.Crypto.CipherParams.Iv = hex.EncodeToString(iv)

	err = ks.WriteKeyStoreFile()
	assert.NoError(t, err)

	return ks
}

func TestEncDec(t *testing.T) {
	ks := NewKeyStore(t)
	ciphertext, err := hex.DecodeString(ks.Crypto.Ciphertext)
	assert.NoError(t, err)
	iv, err := hex.DecodeString(ks.Crypto.CipherParams.Iv)
	assert.NoError(t, err)

	ds, err := ks.Scrypt("123456")
	assert.NoError(t, err)

	dcryptKey, err := AesDecrypt(ciphertext, ds, iv)
	assert.NoError(t, err)
	fmt.Printf("dcrypt key is: %x, len is: %d\n\n", big.NewInt(0).SetBytes(dcryptKey), len(dcryptKey))
}
