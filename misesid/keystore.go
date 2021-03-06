/*
	keystore read & write keystore file, encode and decode privateKey
*/
package misesid

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/btcsuite/btcd/btcec"
	"github.com/mises-id/sdk/bip39"
	"golang.org/x/crypto/scrypt"
	"golang.org/x/crypto/sha3"
)

const (
	ErrorMsgWrongPassword = "wrong password"
)

type KdfParams struct {
	Dklen int
	Salt  string
	N     int
	R     int
	P     int
}

type CipherParams struct {
	Iv string
}

type Crypto struct {
	Cipher       string
	Ciphertext   string
	CipherParams CipherParams
	Kdf          string
	KdfParams    KdfParams
	Mac          string
}

type KeyStore struct {
	Version    int
	MId        string
	PubKey     string
	PrivateKey *btcec.PrivateKey
	PublicKey  *btcec.PublicKey
	Crypto     Crypto
}

//var Ks KeyStore
var keyStoreFile = "keystore" // keystore file located in ./config
var Ver = 3
var KdfMethod = "scrypt"
var CipherMethod = "aes-128-ctr"

func DeleteKeyStoreFile() error {
	return os.Remove(keyStoreFile)
}
func ReadKeyStoreFile() (*KeyStore, error) {
	content, err := ioutil.ReadFile(keyStoreFile)
	if err != nil {
		return nil, err
	}

	var ks KeyStore

	err = json.Unmarshal(content, &ks)
	if err != nil {
		return nil, err
	}

	return &ks, nil
}

func (ks *KeyStore) WriteKeyStoreFile() error {
	if ks.Crypto.Mac == "" {
		return fmt.Errorf("invalid mac")
	}

	content, err := json.Marshal(*ks)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(keyStoreFile, content, 0644)
}

// init kdf parameters, must be called before encoding or decoding funcs
func (ks *KeyStore) InitKdfParam() {
	salt := make([]byte, 16)
	_, _ = rand.Read(salt)

	ks.Crypto.KdfParams.Salt = hex.EncodeToString(salt)
	ks.Crypto.KdfParams.Dklen = 16
	ks.Crypto.KdfParams.N = 32768
	ks.Crypto.KdfParams.R = 8
	ks.Crypto.KdfParams.P = 1
}

// compute s decoding key from local password
func (ks *KeyStore) Scrypt(password string) ([]byte, error) {
	salt, err := hex.DecodeString(ks.Crypto.KdfParams.Salt)
	if err != nil {
		return nil, err
	}

	ck, err := scrypt.Key(
		[]byte(password), salt,
		ks.Crypto.KdfParams.N,
		ks.Crypto.KdfParams.R,
		ks.Crypto.KdfParams.P,
		ks.Crypto.KdfParams.Dklen,
	)
	if err != nil {
		return nil, err
	}

	return ck, nil
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)

	return append(ciphertext, padtext...)
}

// encode private key using s key
func AesEncrypt(origData, key []byte) ([]byte, []byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	blockSize := block.BlockSize()
	origData = PKCS5Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, iv)
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)

	return crypted, iv, nil
}

func PKCS5UnPadding(origData []byte) ([]byte, error) {
	length := len(origData)
	unpadding := int(origData[length-1])

	if length >= unpadding {
		return origData[:(length - unpadding)], nil
	}
	return nil, fmt.Errorf("unpadding > len(origData, can not unpadding")
}

func AesDecrypt(crypted, key []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockMode := cipher.NewCBCDecrypter(block, iv)
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData, err = PKCS5UnPadding(origData)
	if err != nil {
		return nil, err
	}

	return origData, nil
}

func GenMac(ck []byte, ciphertext []byte) ([]byte, error) {
	data := append(ck[1:], ciphertext[:]...)

	sha3 := sha3.New256()
	if _, err := sha3.Write(data); err != nil {
		return nil, err
	}

	return sha3.Sum(nil), nil
}

func Void() {

}

func RandomMnemonics() (string, error) {
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		return "", err
	}

	mnemonics, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", err
	}

	return mnemonics, nil
}
