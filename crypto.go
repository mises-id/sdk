package sdk

import (
	tmsecp256k1 "github.com/tendermint/tendermint/crypto/secp256k1"
)

type misesKeySeed struct {
	MKeySeed
	mnemonic   string
	massPhrase string
}

func (c *misesSdk) RandomSeed() (MKeySeed, error) {
	return nil, nil
}

func (c *misesSdk) RestoreSeed(mnemonic string, passPhrase string) (MKeySeed, error) {
	return nil, nil
}

// Generate private key from mnemonic and compute address
func (ctx *misesKeySeed) getPrivateKey() (*tmsecp256k1.PrivKey, error) {
	return nil, nil
}

// Derive address from the mnemonic
func (ctx *misesKeySeed) getAddress() (string, error) {

	return "", nil
}
