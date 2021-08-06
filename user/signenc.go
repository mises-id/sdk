package user

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/btcsuite/btcd/btcec"
	"github.com/mises-id/sdk/misesid"
	"github.com/mises-id/sdk/types"
)

var sep = "&"

// sign msg using user's private key
func Sign(cuser types.MUser, msg string) (string, string, error) {
	privKey := cuser.PrivateKey()
	if privKey == nil {
		return "", "", fmt.Errorf("private key or public key not available")
	}

	t := strconv.FormatInt(time.Now().UTC().Unix(), 10)
	msg = msg + sep + t
	/*
		dt, err := hex.DecodeString(msg)
		if err != nil {
			return "", err
		}
	*/
	mhash := sha256.Sum256([]byte(msg))

	sig, err := privKey.Sign(mhash[:])
	if err != nil {
		return "", "", err
	}

	derString := sig.Serialize()

	signed := hex.EncodeToString(derString)

	return signed, t, nil
}

// verify msg is sent by user who has the private key
func Verify(msg string, pubKeyStr string, sigStr string) error {
	mhash := sha256.Sum256([]byte(msg))
	publicKeyBytes, err := hex.DecodeString(pubKeyStr)
	if err != nil {
		return err
	}
	pubKey, err := btcec.ParsePubKey(publicKeyBytes, btcec.S256())
	if err != nil {
		return fmt.Errorf("can not parse public key")
	}
	sigByte, err := hex.DecodeString(sigStr)
	if err != nil {
		return err
	}
	sig, err := btcec.ParseDERSignature(sigByte, btcec.S256())
	if err != nil {
		return fmt.Errorf("can not parse signature")
	}

	if !sig.Verify(mhash[:], pubKey) {
		return fmt.Errorf("wrong signature")
	}
	return nil
}

func parseSigned(signed string) (*btcec.PublicKey, []byte, *btcec.Signature, error) {
	param := strings.Fields(signed)
	num := len(param)
	if num != 3 {
		return nil, nil, nil, fmt.Errorf("incorrect num of fields of signed string")
	}

	// parse publicKey, 1. parse derstring; 2.hex.DecodeString to []byte
	publicKey, err := hex.DecodeString(param[0])
	if err != nil {
		return nil, nil, nil, err
	}

	pubKey, err := btcec.ParsePubKey(publicKey, btcec.S256())
	if err != nil {
		return nil, nil, nil, fmt.Errorf("can not parse public key")
	}

	mhash, err := hex.DecodeString(param[1])
	if err != nil {
		return nil, nil, nil, err
	}

	sigByte, err := hex.DecodeString(param[2])
	if err != nil {
		return nil, nil, nil, err
	}

	sig, err := btcec.ParseDERSignature(sigByte, btcec.S256())
	if err != nil {
		return nil, nil, nil, fmt.Errorf("can not parse signature")
	}

	return pubKey, mhash, sig, nil
}

func AesKey(cuser types.MUser) ([]byte, error) {
	privKey := cuser.PrivKEY()
	if privKey == "" {
		return nil, fmt.Errorf("private key or public key not available")
	}
	mhash := sha256.Sum256([]byte(privKey))
	return mhash[:], nil
}
func Encrypt(cuser types.MUser, msg []byte) (string, string, error) {
	keyByte, err := AesKey(cuser)

	if err != nil {
		return "", "", err
	}
	cipherByte, ivByte, err := misesid.AesEncrypt(msg, keyByte)
	if err != nil {
		return "", "", err
	}
	cipher := base64.StdEncoding.EncodeToString(cipherByte)
	iv := base64.StdEncoding.EncodeToString(ivByte)

	return cipher, iv, nil

}

func Decrypt(cuser types.MUser, encData string, iv string) ([]byte, error) {
	keyByte, err := AesKey(cuser)

	if err != nil {
		return nil, err
	}

	cipherByte, err := base64.StdEncoding.Strict().DecodeString(encData)
	if err != nil {
		return nil, err
	}
	ivByte, err := base64.StdEncoding.Strict().DecodeString(iv)
	if err != nil {
		return nil, err
	}

	msgByte, err := misesid.AesDecrypt(cipherByte, keyByte, ivByte)
	if err != nil {
		return nil, err
	}
	return msgByte, nil
}
