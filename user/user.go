package user

import (
	"encoding/hex"
	"encoding/json"
	"strings"

	"github.com/btcsuite/btcd/btcec"
	"github.com/mises-id/sdk/misesid"
)

type MUser interface {
	MisesID() string
	PubKEY() string
	PrivKEY() string
	PrivateKey() *btcec.PrivateKey
	PublicKey() *btcec.PublicKey
	Info() MisesUserInfo
	SetInfo(info MisesUserInfo) string
	GetFollow(appDid string) []string
	SetFollow(followingId string, op bool, appDid string) string
	LoadKeyStore(passPhrase string) error
	IsRegistered() (bool, error)
	Register(info MisesUserInfo, appDid string) error
}

// type MUserInfo interface {
// 	Name() string
// 	Gender() string
// 	AvatarDid() string   //did of avatar file did:mises:0123456789abcdef/avatar
// 	AvatarThumb() []byte //avatar thumb is a bitmap
// 	HomePage() string    //url
// 	Emails() []string
// 	Telphones() []string
// 	Intro() string
// }

type MisesUserInfo struct {
	Name        string
	Gender      string
	AvatarId    string
	AvatarThumb []byte
	HomePage    string
	Emails      []string
	Telephones  []string
	Intro       string
}

type MisesUser struct {
	MUser
	mid        string
	privKey    string
	pubKey     string
	privateKey *btcec.PrivateKey
	publicKey  *btcec.PublicKey
	uinfo      MisesUserInfo
	//	isRegister bool
}

// read keystore file, decode private key
func (user *MisesUser) LoadKeyStore(passPhrase string) error {
	err := misesid.ReadKeyStoreFile()
	if err != nil {
		return err
	}

	s, err := misesid.Scrypt(passPhrase)
	if err != nil {
		return err
	}

	ctext, err := hex.DecodeString(misesid.Ks.Crypto.Ciphertext)
	if err != nil {
		return err
	}

	privKey, err := misesid.AesDecrypt(ctext, s)
	if err != nil {
		return err
	}

	user.privateKey, user.publicKey = btcec.PrivKeyFromBytes(btcec.S256(), privKey)
	user.mid = misesid.Ks.MId
	user.pubKey = misesid.Ks.PubKey
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

func (user *MisesUser) Info() MisesUserInfo {
	uib, err := GetUInfo(user, user.MisesID())
	if err != nil {
		return user.uinfo
	}

	var ui MisesUserInfo
	err = json.Unmarshal(uib, &ui)
	if err != nil {
		return user.uinfo
	}

	user.uinfo = ui
	return user.uinfo
}

func (user *MisesUser) SetInfo(info MisesUserInfo) string {
	session, err := SetUInfo(user, info)
	if err != nil {
		return ""
	}

	user.uinfo = info
	return session
}

func (user *MisesUser) GetFollow(appid string) []string {
	f, err := GetFollowing(user, user.MisesID())
	if err != nil {
		return nil
	}

	return ParseFollowing(f)
}

func ParseFollowing(f []byte) []string {
	follows := string(f)
	return strings.Fields(follows)
}

func (user *MisesUser) SetFollow(followingId string, op bool, appid string) string {
	var operator string
	if op {
		operator = "follow"
	} else {
		operator = "unfollow"
	}

	session, err := SetFollowing(user, followingId, operator)
	if err != nil {
		return ""
	}

	return session
}
