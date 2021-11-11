package user

import (
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcutil/bech32"

	"github.com/btcsuite/btcd/btcec"
	"github.com/ebfe/keccak"
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
	var u MisesUser

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
	pubKeyByte := pubKey.SerializeUncompressed()

	k := keccak.New256()
	k.Write(pubKeyByte)
	publicKey := k.Sum(nil)

	mid, err := ConvertAndEncode(
		types.AddressPrefix,
		publicKey[len(publicKey)-20:],
	)
	if err != nil {
		return nil, err
	}
	u.mid = types.MisesIDPrefix + mid
	u.privKey = hex.EncodeToString(privateKey)
	u.pubKey = hex.EncodeToString(pubKeyByte)
	u.privateKey = privKey
	u.publicKey = pubKey
	// privateKey, publicKey & misesId generated, and new MisesUser Created

	// add user to userMgr, set to Active User
	userMgr.AddUser(&u)
	userMgr.SetActiveUser(u.mid)

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
	ks.MId = u.mid
	ks.PubKey = u.pubKey
	ks.Version = misesid.Ver
	ks.Crypto.Ciphertext = hex.EncodeToString(ciphertext)
	ks.Crypto.Cipher = misesid.CipherMethod
	ks.Crypto.Kdf = misesid.KdfMethod
	ks.Crypto.Mac = hex.EncodeToString(mac)
	ks.Crypto.CipherParams.Iv = hex.EncodeToString(iv)

	if err = ks.WriteKeyStoreFile(); err != nil {
		return nil, err
	}

	return &u, nil
}

func ConvertAndEncode(hrp string, data []byte) (string, error) {
	converted, err := bech32.ConvertBits(data, 8, 5, true)
	if err != nil {
		return "", fmt.Errorf("encoding bech32 failed: %w", err)
	}

	return bech32.Encode(hrp, converted)
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
