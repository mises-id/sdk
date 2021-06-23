package sdk

import (
	"fmt"
)

// Debugf level formatted messagctx.logger.
func (ctx *MisesAgent) Debugf(msg string, v ...interface{}) {
	if ctx.Debug {
		ctx.logger.Debug(fmt.Sprintf(msg, v...))
	}
}

// Infof level formatted messagctx.logger.
func (ctx *MisesAgent) Infof(msg string, v ...interface{}) {
	if ctx.Debug {
		ctx.logger.Info(fmt.Sprintf(msg, v...))
	}
}

// Warnf level formatted messagctx.logger.
func (ctx *MisesAgent) Warnf(msg string, v ...interface{}) {
	if ctx.Debug {
		ctx.logger.Info(fmt.Sprintf(msg, v...))
	}
}

// Errorf level formatted messagctx.logger.
func (ctx *MisesAgent) Errorf(msg string, v ...interface{}) {
	if ctx.Debug {
		ctx.logger.Error(fmt.Sprintf(msg, v...))
	}
}

// Fatalf level formatted messagctx.logger.
func (ctx *MisesAgent) Fatalf(msg string, v ...interface{}) {
	if ctx.Debug {
		ctx.logger.Error(fmt.Sprintf(msg, v...))
	}
}
