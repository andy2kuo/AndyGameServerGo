package socketserver

import (
	"ak-project-server/logger"
	"context"
	"sync"
)

type ConnContext struct {
	sync.RWMutex
	context.Context

	logger *logger.Logger
}

func (ctx *ConnContext) Log() *logger.Logger {
	return ctx.logger
}
