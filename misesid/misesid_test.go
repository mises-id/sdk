package misesid

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/btcsuite/btcd/btcec"
	"github.com/ebfe/keccak"
	"github.com/mises-id/sdk/bip39"
	"github.com/tyler-smith/assert"
)

func TestNewMisesId(t *testing.T) {
	entropy, err := bip39.NewEntropy(128)
	assert.NoError(t, err)

	mnemonic, err := bip39.NewMnemonic(entropy)
	assert.NoError(t, err)

	seed, err := bip39.NewSeed(mnemonic, "MISES")
	assert.NoError(t, err)

	Mid.masterKey = Seed2MasterKey(seed)
	privKeyByte := Mid.masterKey[0:32]
	Mid.chainCode = Mid.masterKey[32:]
	fmt.Printf("master key is: %x, len is: %d\n", big.NewInt(0).SetBytes(Mid.masterKey), len(Mid.masterKey))
	fmt.Printf("master private key is: %x, len is: %d\n", big.NewInt(0).SetBytes(privKeyByte), len(privKeyByte))
	fmt.Printf("master chain code is: %x, len is: %d\n", big.NewInt(0).SetBytes(Mid.chainCode), len(Mid.chainCode))

	privKey, pubKey := btcec.PrivKeyFromBytes(btcec.S256(), privKeyByte)
	Mid.privKey = privKey.Serialize()
	pubKeyByte := pubKey.SerializeCompressed()
	fmt.Printf("private key is: %x, len is: %d\n", big.NewInt(0).SetBytes(Mid.privKey), len(Mid.privKey))
	fmt.Printf("uncompressed public key is: %x, len is: %d\n", big.NewInt(0).SetBytes(pubKeyByte), len(pubKeyByte))

	k := keccak.New256()
	k.Write(pubKeyByte)
	Mid.pubKey = k.Sum(nil)
	fmt.Printf("pubKey is: %x, len is: %d\n", big.NewInt(0).SetBytes(Mid.pubKey), len(Mid.pubKey))

	Mid.id = Mid.pubKey[len(Mid.pubKey)-20:]
	fmt.Printf("misesid is %x, len is: %d\n", big.NewInt(0).SetBytes(Mid.id), len(Mid.id))
}
