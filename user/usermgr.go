package user

import (
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/btcec"
	"github.com/mises-id/sdk/bip39"
	"github.com/mises-id/sdk/misesid"
	"github.com/mises-id/sdk/types"
)

var _ types.MUserMgr = &MisesUserMgr{}

type MisesUserMgr struct {
	activeUser types.MUser
	users      []types.MUser
}

// create user & his misesid, private key & uncompressed public key, write a keystore file
func (userMgr *MisesUserMgr) CreateUser(mnemonic string, passPhrase string) (types.MUser, error) {
	// in bip39, passPhrase will affect the wallet address,
	// that's not what we need now, so we simply set the passwaor to empty
	seed, err := bip39.NewSeed(mnemonic, "")
	if err != nil {
		return nil, err
	}

	masterKey := misesid.Seed2MasterKey(seed)
	privKeyByte := masterKey[0:32]
	//	chainCode := masterKey[32:]

	privKey, pubKey := btcec.PrivKeyFromBytes(btcec.S256(), privKeyByte)
	privateKey := privKey.Serialize()
	pubKeyByte := pubKey.SerializeCompressed()

	mid, err := misesid.ConvertAndEncode(
		types.AddressPrefix,
		misesid.PubKeyAddrBytes(pubKeyByte),
	)
	if err != nil {
		return nil, err
	}
	u := NewMisesUser(mid, privateKey)
	// privateKey, publicKey & misesId generated, and new MisesUser Created

	// add user to userMgr, set to Active User
	userMgr.AddUser(u)
	_ = userMgr.SetActiveUser(u.MisesID())

	// encrypt privatKey, and write keystore file
	ks := &misesid.KeyStore{}
	ks.InitKdfParam()

	s, err := ks.Scrypt(passPhrase)
	if err != nil {
		return nil, err
	}

	var ciphertext, iv []byte
	if ciphertext, iv, err = misesid.AesEncrypt(privKeyByte, s); err != nil {
		return nil, err
	}

	var mac []byte
	if mac, err = misesid.GenMac(s, ciphertext); err != nil {
		return nil, err
	}

	// write keystore file
	ks.MId = u.MisesID()
	ks.PubKey = u.PubKey()
	ks.Version = misesid.Ver
	ks.Crypto.Ciphertext = hex.EncodeToString(ciphertext)
	ks.Crypto.Cipher = misesid.CipherMethod
	ks.Crypto.Kdf = misesid.KdfMethod
	ks.Crypto.Mac = hex.EncodeToString(mac)
	ks.Crypto.CipherParams.Iv = hex.EncodeToString(iv)

	if err = ks.WriteKeyStoreFile(); err != nil {
		return nil, err
	}

	return u, nil
}

func (userMgr *MisesUserMgr) ListUsers() []types.MUser {
	return userMgr.users
}

func (userMgr *MisesUserMgr) AddUser(user types.MUser) {
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

func (userMgr *MisesUserMgr) ActiveUser() types.MUser {
	return userMgr.activeUser
}
