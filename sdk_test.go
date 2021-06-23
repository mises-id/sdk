package sdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequiresChainId(t *testing.T) {
	assert := assert.New(t)

	options := &MSdkOption{}
	if _, err := NewMSdk(options); err == nil {
		t.Fatalf("chain id requirement was not validated")
	} else {
		assert.Equal(err.Error(), "chain id is required")
	}
}
