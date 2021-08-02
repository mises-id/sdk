package user

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/btcsuite/btcd/btcec"
)

var sep = "&"

// sign msg using user's private key
func Sign(cuser MUser, msg string) (string, string, error) {
	privKey := cuser.PrivateKey()
	if privKey == nil {
		return "", "", fmt.Errorf("private key or public key not available")
	}

	tb, err := time.Now().MarshalText()
	if err != nil {
		return "", "", err
	}
	t := hex.EncodeToString(tb)
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
func Verify(signed string) bool {
	pubKey, mhash, sig, err := parseSigned(signed)
	if err != nil {
		return false
	}

	return sig.Verify(mhash, pubKey)
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
