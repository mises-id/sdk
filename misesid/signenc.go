package misesid

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/btcsuite/btcd/btcec"
	"github.com/cosmos/btcutil/bech32"
	"github.com/mises-id/sdk/types"
	"golang.org/x/crypto/ripemd160"
)

func CheckMisesID(misesID string, pubKeyStr string) error {
	publicKeyBytes, err := hex.DecodeString(pubKeyStr)
	if err != nil {
		return err
	}
	if len(publicKeyBytes) != btcec.PubKeyBytesLenCompressed {
		return fmt.Errorf("pubkey length not 33")
	}

	mid, err := ConvertAndEncode(
		types.AddressPrefix,
		PubKeyAddrBytes(publicKeyBytes),
	)
	if err != nil {
		return err
	}
	if types.MisesIDPrefix+mid != misesID && types.MisesAppIDPrefix+mid != misesID {
		return fmt.Errorf("mises_id[%s] not matching pubkey[%s]", misesID, pubKeyStr)
	}
	return nil

}

// verify msg is sent by user who has the private key
func Verify(msg string, pubKeyStr string, sigStr string) error {
	mhash := sha256.Sum256([]byte(msg))
	publicKeyBytes, err := hex.DecodeString(pubKeyStr)
	if err != nil {
		return fmt.Errorf("can not parse signature")
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

func Encrypt(cuser types.MSigner, msg []byte) (string, string, error) {
	keyByte, err := cuser.AesKey()

	if err != nil {
		return "", "", err
	}
	cipherByte, ivByte, err := AesEncrypt(msg, keyByte)
	if err != nil {
		return "", "", err
	}
	cipher := base64.StdEncoding.EncodeToString(cipherByte)
	iv := base64.StdEncoding.EncodeToString(ivByte)

	return cipher, iv, nil

}

func Decrypt(cuser types.MSigner, encData string, iv string) ([]byte, error) {
	keyByte, err := cuser.AesKey()

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

	msgByte, err := AesDecrypt(cipherByte, keyByte, ivByte)
	if err != nil {
		return nil, err
	}
	return msgByte, nil
}

func PubKeyAddrBytes(pubkey []byte) []byte {
	sha := sha256.Sum256(pubkey)
	hasherRIPEMD160 := ripemd160.New()
	hasherRIPEMD160.Write(sha[:]) // does not error
	pubKeyAddrBytes := hasherRIPEMD160.Sum(nil)
	return pubKeyAddrBytes
}

func ConvertAndEncode(hrp string, data []byte) (string, error) {
	converted, err := bech32.ConvertBits(data, 8, 5, true)
	if err != nil {
		return "", fmt.Errorf("encoding bech32 failed: %w", err)
	}

	return bech32.Encode(hrp, converted)
}
