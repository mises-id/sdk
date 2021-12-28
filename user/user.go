package user

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/btcec"
	"github.com/mises-id/sdk/misesid"
	"github.com/mises-id/sdk/types"
)

var _ types.MUser = &MisesUser{}

type MisesUser struct {
	mid          string
	privKey      string
	pubKey       string
	privateKey   *btcec.PrivateKey
	publicKey    *btcec.PublicKey
	uinfo        misesid.MisesUserInfo
	isRegistered bool
}

func NewMisesUser(mid string, prikeyBytes []byte) *MisesUser {
	privKey, pubKey := btcec.PrivKeyFromBytes(btcec.S256(), prikeyBytes)
	pubkeyBytes := pubKey.SerializeCompressed()
	u := MisesUser{
		mid:        types.MisesIDPrefix + mid,
		privKey:    hex.EncodeToString(prikeyBytes),
		pubKey:     hex.EncodeToString(pubkeyBytes),
		privateKey: privKey,
		publicKey:  pubKey,
	}
	return &u
}

// read keystore file, decode private key
func (user *MisesUser) LoadKeyStore(passPhrase string) error {
	ks, err := misesid.ReadKeyStoreFile()
	if err != nil {
		return err
	}

	user.mid = ks.MId
	user.pubKey = ks.PubKey

	s, err := ks.Scrypt(passPhrase)
	if err != nil {
		return fmt.Errorf(misesid.ErrorMsgWrongPassword)
	}

	ctext, err := hex.DecodeString(ks.Crypto.Ciphertext)
	if err != nil {
		return fmt.Errorf(misesid.ErrorMsgWrongPassword)
	}
	iv, err := hex.DecodeString(ks.Crypto.CipherParams.Iv)
	if err != nil {
		return fmt.Errorf(misesid.ErrorMsgWrongPassword)
	}

	privKey, err := misesid.AesDecrypt(ctext, s, iv)
	if err != nil {
		return fmt.Errorf(misesid.ErrorMsgWrongPassword)
	}

	privateKey, publicKey := btcec.PrivKeyFromBytes(btcec.S256(), privKey)

	pubKeyByte := publicKey.SerializeUncompressed()

	if ks.PubKey != hex.EncodeToString(pubKeyByte) {
		return fmt.Errorf(misesid.ErrorMsgWrongPassword)
	}

	user.privateKey = privateKey
	user.publicKey = publicKey
	user.privKey = hex.EncodeToString(privKey)

	return nil
}

func (user MisesUser) MisesID() string {
	return user.mid
}

func (user MisesUser) PubKey() string {
	return user.pubKey
}

func (user MisesUser) PrivKey() string {
	return user.privKey
}

func (user MisesUser) PrivateKey() *btcec.PrivateKey {
	return user.privateKey
}

func (user MisesUser) PublicKey() *btcec.PublicKey {
	return user.publicKey
}

func (user *MisesUser) Info() types.MUserInfo {
	uInfo, err := misesid.GetUInfo(user, user.MisesID())
	if err != nil {
		return &misesid.MisesUserInfoReadonly{user.uinfo}
	}

	user.uinfo = *uInfo
	return &misesid.MisesUserInfoReadonly{*uInfo}
}

func (user *MisesUser) SetInfo(info types.MUserInfo) (string, error) {

	minfo := misesid.NewMisesUserInfo(info)
	session, err := misesid.SetUInfo(user, minfo)
	if err != nil {
		return "", err
	}

	user.uinfo = *minfo
	return session, nil
}

func (user *MisesUser) GetFollow(appid string) []string {

	f, err := misesid.GetFollowing(user, user.MisesID())
	if err != nil {
		return []string{}
	}

	return f
}

func (user *MisesUser) SetFollow(followingId string, op bool, appid string) (string, error) {

	var operator string
	if op {
		operator = "follow"
	} else {
		operator = "unfollow"
	}

	return misesid.SetFollowing(user, followingId, operator)
}

func (user *MisesUser) IsRegistered() error {
	if user.isRegistered {
		return nil
	}
	_, err := misesid.GetMisesID(user, user.MisesID())
	if err != nil {
		return err
	}
	user.isRegistered = true
	return nil

}

func (user *MisesUser) Register(appID string) (string, error) {
	return misesid.CreateMisesID(user)
}

// sign msg using user's private key
func (cuser *MisesUser) Sign(msg string) (string, error) {
	privKey := cuser.PrivateKey()
	if privKey == nil {
		return "", fmt.Errorf("private key or public key not available")
	}

	mhash := sha256.Sum256([]byte(msg))

	sig, err := privKey.Sign(mhash[:])
	if err != nil {
		return "", err
	}

	derString := sig.Serialize()

	signed := hex.EncodeToString(derString)

	return signed, nil
}
func (cuser *MisesUser) AesKey() ([]byte, error) {
	privKey := cuser.PrivKey()
	if privKey == "" {
		return nil, fmt.Errorf("private key or public key not available")
	}
	privKeyBytes, err := hex.DecodeString(privKey)
	if err != nil {
		return nil, err
	}
	mhash := sha256.Sum256([]byte(privKeyBytes))
	return mhash[:], nil
}

func (cuser *MisesUser) Signer() types.MSigner {
	return cuser
}
