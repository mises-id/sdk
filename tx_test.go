package sdk

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	assert := assert.New(t)

	ctx := &Test{}
	if err := ctx.TestSetUp(); err != nil {
		t.Fatalf("%s", err)
	}
	defer ctx.TestTearDown()

	if err := ctx.Agent.Create(ctx.Key1, ctx.Value1, TestGasInfo(), nil); err != nil {
		t.Fatalf("%s", err)
	}
	if v, err := ctx.Agent.Read(ctx.Key1); err != nil {
		t.Fatalf("%s", err)
	} else {
		assert.Equal(v, ctx.Value1)
	}
}

func TestCreateWithLeaseInfo(t *testing.T) {
	// assert := assert.New(t)

	ctx := &Test{}
	if err := ctx.TestSetUp(); err != nil {
		t.Fatalf("%s", err)
	}
	defer ctx.TestTearDown()

	if err := ctx.Agent.Create(ctx.Key1, ctx.Value1, TestGasInfo(), &LeaseInfo{Seconds: 60}); err != nil {
		t.Fatalf("%s", err)
	}
}

func TestCreateValidatesGasInfo(t *testing.T) {
	assert := assert.New(t)

	ctx := &Test{}
	if err := ctx.TestSetUp(); err != nil {
		t.Fatalf("%s", err)
	}
	defer ctx.TestTearDown()

	err := ctx.Agent.Create(ctx.Key1, ctx.Value1, &GasInfo{MaxFee: 1}, nil)
	assert.True(err != nil) // todo check details
}

func TestCreatesFailsIfKeyContainsSlash(t *testing.T) {
	assert := assert.New(t)

	ctx := &Test{}
	if err := ctx.TestSetUp(); err != nil {
		t.Fatalf("%s", err)
	}
	defer ctx.TestTearDown()

	err := ctx.Agent.Create("123/", ctx.Value1, TestGasInfo(), nil)
	assert.True(err != nil)
	assert.True(strings.Contains(err.Error(), "Key cannot contain a slash"))
}
