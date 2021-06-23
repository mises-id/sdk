package sdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequiresMnemonic(t *testing.T) {
	assert := assert.New(t)
	options := &MSdkOption{}
	sdk, _ := NewMSdk(options)
	seed, _ := sdk.RandomSeed()
	if _, err := newMisesAgent("", seed.(*MisesKeySeed)); err == nil {
		t.Fatalf("mnemonic requirement was not validated")
	} else {
		assert.Equal(err.Error(), "mnemonic is required")
	}
}
