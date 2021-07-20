package bip39

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/tyler-smith/assert"
)

func TestNewEntropy(t *testing.T) {
	entropy, err := NewEntropy(128)
	assert.Nil(t, err)
	e := big.NewInt(0).SetBytes(entropy)
	fmt.Printf("entropy is : %x\n", e)
}

func TestNewNemonic(t *testing.T) {
	entropy, err := NewEntropy(128)
	assert.Nil(t, err)

	mnemonic, err := NewMnemonic(entropy)
	assert.Nil(t, err)
	fmt.Printf("mnemonic is: %s\n", mnemonic)

	words, b := splitMnemonic(mnemonic)
	if !b {
		fmt.Printf("invalid mnemonic, %s\n", mnemonic)
		return
	}
	for _, w := range words {
		fmt.Printf("%s ", w)
	}
	fmt.Printf("\n")
}

func TestRestoreEntropy(t *testing.T) {
	entropy, err := NewEntropy(128)
	assert.Nil(t, err)

	mnemonic, err := NewMnemonic(entropy)
	assert.Nil(t, err)

	restoreEntropy, err := RestoreEntropy(mnemonic)
	assert.Nil(t, err)

	e1 := big.NewInt(0).SetBytes(entropy)
	e2 := big.NewInt(0).SetBytes(restoreEntropy)

	fmt.Printf("origin entropy is : %x\n", e1)
	fmt.Printf("restored entropy is : %x\n", e2)
}

func TestNewSeed(t *testing.T) {
	entropy, err := NewEntropy(128)
	assert.Nil(t, err)

	mnemonic, err := NewMnemonic(entropy)
	assert.Nil(t, err)

	seed, err := NewSeed(mnemonic, "TREZOR")
	assert.Nil(t, err)

	s := big.NewInt(0).SetBytes(seed)
	fmt.Printf("seed is : %x\n", s)
}
