package misesid

import (
	"fmt"
	"log"
	"math/big"
	"testing"
)

func TestReadKeyStoreFile(t *testing.T) {
	err := ReadKeyStoreFile()
	if err != nil {
		fmt.Printf("keystore version is: %d\n", Ks.Version)
		fmt.Printf("keystore id is: %s\n", Ks.MId)
		fmt.Printf("keystore address is: %s\n", Ks.PubKey)
		fmt.Printf("keystore kdf is: %s\n", Ks.Crypto.Kdf)
		fmt.Printf("keystore dklen is: %d\n", Ks.Crypto.KdfParams.Dklen)
		fmt.Printf("keystore salt is: %s\n", Ks.Crypto.KdfParams.Salt)
		fmt.Printf("keystore n is: %d\n", Ks.Crypto.KdfParams.N)
		fmt.Printf("keystore r is: %d\n", Ks.Crypto.KdfParams.R)
		fmt.Printf("keystore p is: %d\n", Ks.Crypto.KdfParams.P)
		fmt.Printf("keystore cipher is: %s\n", Ks.Crypto.Cipher)
		fmt.Printf("keystore ciphertext is: %s\n", Ks.Crypto.Ciphertext)
		fmt.Printf("keystore iv is: %s\n", Ks.Crypto.CipherParams.Iv)
		fmt.Printf("keystore mac is: %s\n", Ks.Crypto.Mac)
	}

	err = WriteKeyStoreFile()
	if err != nil {
		fmt.Printf("keystorew write succeed\n\n")
	}
}

func TestEncDec(t *testing.T) {
	InitKdfParam()

	s, err := Scrypt("123456")
	fmt.Printf("s key is: %x, len is: %d\n", big.NewInt(0).SetBytes(s), len(s))
	privateKey := "d88169685ceeaaaef5c619bcdaca3b44c323372e0d777720d4a28efe3c658aa2"
	privKey := []byte(privateKey)
	fmt.Printf("private key is: %x, len is: %d\n", big.NewInt(0).SetBytes(privKey), len(privKey))

	var ciphertext, iv []byte
	if ciphertext, iv, err = AesEncrypt(privKey, s); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ciphertext is: %x, len is: %d\n", big.NewInt(0).SetBytes(ciphertext), len(ciphertext))

	var mac []byte
	if mac, err = GenMac(s, ciphertext); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("mac is: %x, len is: %d\n", big.NewInt(0).SetBytes(mac), len(mac))

	ds, err := Scrypt("123456")
	fmt.Printf("ds key is: %x, len is: %d\n", big.NewInt(0).SetBytes(ds), len(ds))

	var dmac []byte
	if dmac, err = GenMac(ds, ciphertext); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ds'mac is: %x, len is: %d\n", big.NewInt(0).SetBytes(dmac), len(dmac))

	var dcryptKey []byte
	if dcryptKey, err = AesDecrypt(ciphertext, ds, iv); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("dcrypt key is: %x, len is: %d\n\n", big.NewInt(0).SetBytes(dcryptKey), len(dcryptKey))
}
