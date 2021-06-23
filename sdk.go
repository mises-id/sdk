package sdk

import (
	"github.com/prometheus/common/log"
)

type MSdkOption struct {
	ChainId string
	Debug   bool
}

type MisesSdk struct {
	MSdk
	options *MSdkOption
	logger  log.Logger
}

func (ctx *MisesSdk) setupLogger() {
	ctx.logger = log.NewNopLogger()
}

func NewMSdk(options *MSdkOption) (MSdk, error) {
	if options.ChainId == "" {
		options.ChainId = DEFAULT_CHAIN_ID
	}

	ctx := &MisesSdk{
		options: options,
	}

	ctx.setupLogger()

	return ctx, nil
}
