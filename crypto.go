package sdk

import (
	tmsecp256k1 "github.com/tendermint/tendermint/crypto/secp256k1"
)

type MisesKeySeed struct {
	MKeySeed
	mnemonic   string
	massPhrase string
}

func (c *MisesSdk) RandomSeed() (MKeySeed, error) {
	return nil, nil
}

func (c *MisesSdk) RestoreSeed(mnemonic string, pass_phrase string) (MKeySeed, error) {
	return nil, nil
}

// Generate private key from mnemonic and compute address
func (ctx *MisesKeySeed) genPrivateKey() (*tmsecp256k1.PrivKey, error) {
	return nil, nil
}

// Derive address from the mnemonic
func (ctx *MisesKeySeed) getAddress() (string, error) {

	return "", nil
}
