package sdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccount(t *testing.T) {
	assert := assert.New(t)

	ctx := &test{}
	if err := ctx.testSetUp(); err != nil {
		t.Fatalf("%s", err)
	}
	defer ctx.testTearDown()

	if account, err := ctx.Agent.Account(); err != nil {
		t.Fatalf("%s", err)
	} else {
		assert.True(account.AccountNumber > 0)
		assert.True(account.Sequence > 0)
	}
}

func TestVersion(t *testing.T) {
	assert := assert.New(t)

	ctx := &test{}
	if err := ctx.testSetUp(); err != nil {
		t.Fatalf("%s", err)
	}
	defer ctx.testTearDown()

	if v, err := ctx.Agent.Version(); err != nil {
		t.Fatalf("%s", err)
	} else {
		assert.True(v != "")
	}
}
