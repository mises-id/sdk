package sdk

import (
	"github.com/prometheus/common/log"
)

type MSdkOption struct {
	ChainID string
	Debug   bool
}

type misesSdk struct {
	MSdk
	options *MSdkOption
	logger  log.Logger
}

func (ctx *misesSdk) setupLogger() {
	ctx.logger = log.NewNopLogger()
}

func NewMSdk(options *MSdkOption) (MSdk, error) {
	if options.ChainID == "" {
		options.ChainID = DefaultChainID
	}

	ctx := &misesSdk{
		options: options,
	}

	ctx.setupLogger()

	return ctx, nil
}
