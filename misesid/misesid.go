package misesid

import (
	"crypto/hmac"
	"crypto/sha512"

	"github.com/btcsuite/btcd/btcec"
	"github.com/ebfe/keccak"
)

type MisesId struct {
	masterKey []byte
	chainCode []byte
	privKey   []byte
	pubKey    []byte
	id        []byte
}

var Mid MisesId

// generate Master Key from seed & password
func Seed2MasterKey(seed []byte) []byte {
	hmac := hmac.New(sha512.New, []byte("Mises seed"))
	_, err := hmac.Write(seed)
	if err != nil {
		panic(err)
	}
	return hmac.Sum(nil)
}

// not used, CreateUser instead
func NewMisesId(seed []byte, password string) {
	Mid.masterKey = Seed2MasterKey(seed)
	privKeyByte := Mid.masterKey[0:31]
	Mid.chainCode = Mid.masterKey[32:63]

	privKey, pubKey := btcec.PrivKeyFromBytes(btcec.S256(), privKeyByte)
	Mid.privKey = privKey.Serialize()
	pubKeyByte := pubKey.SerializeUncompressed()

	k := keccak.New256()
	k.Write(pubKeyByte)
	Mid.pubKey = k.Sum(nil)

	Mid.id = Mid.pubKey[len(Mid.pubKey)-20:]
}
