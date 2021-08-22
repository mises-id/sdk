package user

import (
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/btcec"
	"github.com/mises-id/sdk/misesid"
	"github.com/mises-id/sdk/types"
)

var _ types.MUser = &MisesUser{}

const (
	ErrorMsgWrongPassword = "wrong password"
)

type MisesUser struct {
	mid          string
	privKey      string
	pubKey       string
	privateKey   *btcec.PrivateKey
	publicKey    *btcec.PublicKey
	uinfo        MisesUserInfo
	isRegistered bool
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
		return fmt.Errorf(ErrorMsgWrongPassword)
	}

	ctext, err := hex.DecodeString(ks.Crypto.Ciphertext)
	if err != nil {
		return fmt.Errorf(ErrorMsgWrongPassword)
	}
	iv, err := hex.DecodeString(ks.Crypto.CipherParams.Iv)
	if err != nil {
		return fmt.Errorf(ErrorMsgWrongPassword)
	}

	privKey, err := misesid.AesDecrypt(ctext, s, iv)
	if err != nil {
		return fmt.Errorf(ErrorMsgWrongPassword)
	}

	privateKey, publicKey := btcec.PrivKeyFromBytes(btcec.S256(), privKey)

	pubKeyByte := publicKey.SerializeUncompressed()

	if ks.PubKey != hex.EncodeToString(pubKeyByte) {
		return fmt.Errorf(ErrorMsgWrongPassword)
	}

	user.privateKey = privateKey
	user.publicKey = publicKey
	user.privKey = hex.EncodeToString(privKey)

	return nil
}

func (user MisesUser) MisesID() string {
	return user.mid
}

func (user MisesUser) PubKEY() string {
	return user.pubKey
}

func (user MisesUser) PrivKEY() string {
	return user.privKey
}

func (user MisesUser) PrivateKey() *btcec.PrivateKey {
	return user.privateKey
}

func (user MisesUser) PublicKey() *btcec.PublicKey {
	return user.publicKey
}

func (user *MisesUser) Info() types.MUserInfo {
	uInfo, err := GetUInfo(user, user.MisesID())
	if err != nil {
		return &MisesUserInfoReadonly{user.uinfo}
	}

	user.uinfo = *uInfo
	return &MisesUserInfoReadonly{*uInfo}
}

func (user *MisesUser) SetInfo(info types.MUserInfo) (string, error) {

	minfo := NewMisesUserInfo(info)
	session, err := SetUInfo(user, minfo)
	if err != nil {
		return "", err
	}

	user.uinfo = *minfo
	return session, nil
}

func (user *MisesUser) GetFollow(appid string) []string {

	f, err := GetFollowing(user, user.MisesID())
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

	return SetFollowing(user, followingId, operator)
}

func (user *MisesUser) IsRegistered() error {
	if user.isRegistered {
		return nil
	}
	_, err := GetUser(user, user.MisesID())
	if err != nil {
		return err
	}
	user.isRegistered = true
	return nil

}

func (user *MisesUser) Register(appID string) (string, error) {
	return CreateUser(user)
}
