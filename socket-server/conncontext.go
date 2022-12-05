package socketserver

import (
	"github.com/andy2kuo/AndyGameServerGo/logger"
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
