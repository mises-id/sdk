package sdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequiresMnemonic(t *testing.T) {
	assert := assert.New(t)
	ctx := &test{}
	if err := ctx.testSetUp(); err != nil {
		t.Fatalf("%s", err)
	}
	defer ctx.testTearDown()
	seed, _ := ctx.SDK.RandomSeed()
	if _, err := newMisesAgent("", seed.(*misesKeySeed)); err == nil {
		t.Fatalf("mnemonic requirement was not validated")
	} else {
		assert.Equal(err.Error(), "mnemonic is required")
	}
}
