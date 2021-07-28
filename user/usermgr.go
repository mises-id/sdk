package user

import (
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/btcec"
	"github.com/ebfe/keccak"
	"github.com/mises-id/sdk/bip39"
	"github.com/mises-id/sdk/misesid"
)

type MUserMgr interface {
	CreateUser(mnemonic string, passPhrase string) (MUser, error)
	ListUsers() []MUser
	AddUser(user MUser)
	SetActiveUser(userDid string) error
	ActiveUser() MUser
}

type MisesUserMgr struct {
	MUserMgr
	activeUser MUser
	users      []MUser
}

// create user & his misesid, private key & uncompressed public key, write a keystore file
func (userMgr *MisesUserMgr) CreateUser(mnemonic string, passPhrase string) (MUser, error) {
	var u MisesUser

	seed, err := bip39.NewSeed(mnemonic, passPhrase)
	if err != nil {
		return nil, err
	}

	masterKey := misesid.Seed2MasterKey(seed, []byte(passPhrase))
	privKeyByte := masterKey[0:32]
	//	chainCode := masterKey[32:]

	privKey, pubKey := btcec.PrivKeyFromBytes(btcec.S256(), privKeyByte)
	privateKey := privKey.Serialize()
	pubKeyByte := pubKey.SerializeUncompressed()

	k := keccak.New256()
	k.Write(pubKeyByte)
	publicKey := k.Sum(nil)

	u.mid = hex.EncodeToString(publicKey[len(publicKey)-20:])
	u.privKey = hex.EncodeToString(privateKey)
	u.pubKey = hex.EncodeToString(pubKeyByte)
	u.privateKey = privKey
	u.publicKey = pubKey
	// privateKey, publicKey & misesId generated, and new MisesUser Created

	// add user to userMgr, set to Active User
	userMgr.AddUser(&u)
	userMgr.SetActiveUser(u.mid)

	// encrypt privateKey, and write keystore file
	misesid.InitKdfParam()

	s, err := misesid.Scrypt(passPhrase)
	if err != nil {
		return nil, err
	}

	var ciphertext []byte
	if ciphertext, err = misesid.AesEncrypt(privKeyByte, s); err != nil {
		return nil, err
	}

	var mac []byte
	if mac, err = misesid.GenMac(s, ciphertext); err != nil {
		return nil, err
	}

	// write keystore file
	misesid.Ks.MId = u.mid
	misesid.Ks.PubKey = u.pubKey
	misesid.Ks.Version = misesid.Ver
	misesid.Ks.Crypto.Ciphertext = hex.EncodeToString(ciphertext)
	misesid.Ks.Crypto.Cipher = misesid.CipherMethod
	misesid.Ks.Crypto.Kdf = misesid.KdfMethod
	misesid.Ks.Crypto.Mac = hex.EncodeToString(mac)

	if err = misesid.WriteKeyStoreFile(); err != nil {
		return nil, err
	}

	return &u, nil
}

func (userMgr *MisesUserMgr) ListUsers() []MUser {
	return userMgr.users
}

func (userMgr *MisesUserMgr) AddUser(user MUser) {
	userMgr.users = append(userMgr.users, user)
}

func (userMgr *MisesUserMgr) SetActiveUser(uid string) error {
	for _, u := range userMgr.users {
		if uid == u.MisesID() {
			userMgr.activeUser = u
			return nil
		}
	}
	return fmt.Errorf("can not find user(%s) in users manager", uid)
}

func (userMgr *MisesUserMgr) ActiveUser() MUser {
	return userMgr.activeUser
}
