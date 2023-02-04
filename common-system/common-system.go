package commonsystem

import (
	"context"
	"time"

	"github.com/andy2kuo/AndyGameServerGo/database"
	"github.com/andy2kuo/AndyGameServerGo/logger"
)

type SystemCode byte
type SystemEventCode byte
type SystemEvent struct {
	Code SystemCode
	Data interface{}
}

// 共用系統
type ICommonSystem interface {
	GetSystemCode() SystemCode
	Init(*CommonSystemManager, *logger.Logger, *database.MongoConnection, *database.RedisConnection) error
	OnServerStart() error
	OnSystemEventNotify(SystemEvent)
	Close() error
}

type BaseSystem struct {
	ctx       context.Context
	cancel    context.CancelFunc
	manager   *CommonSystemManager
	logger    *logger.Logger
	mongoConn *database.MongoConnection
	redisConn *database.RedisConnection
}

func (b *BaseSystem) GetSystem(sysCode SystemCode) ICommonSystem {
	return b.manager.GetSystem(sysCode)
}

func (b *BaseSystem) Init(_manager *CommonSystemManager, _logger *logger.Logger, _mongoConn *database.MongoConnection, _redisConn *database.RedisConnection) error {
	b.manager = _manager
	b.logger = _logger
	b.mongoConn = _mongoConn
	b.redisConn = _redisConn

	b.ctx = context.TODO()
	b.ctx, b.cancel = context.WithCancel(b.ctx)

	return nil
}

func (b *BaseSystem) Logger() *logger.Logger {
	return b.logger
}

func (b *BaseSystem) MongoConn() *database.MongoConnection {
	return b.mongoConn
}

func (b *BaseSystem) RedisConn() *database.RedisConnection {
	return b.redisConn
}

func (b *BaseSystem) Context() context.Context {
	return b.ctx
}

func (b *BaseSystem) Cancel() context.CancelFunc {
	return b.cancel
}

func (b *BaseSystem) OnServerStart() error {
	return nil
}

func (b *BaseSystem) Start(interval time.Duration, operation func()) {
	go func() {

		for {
			select {
			case <-b.ctx.Done():
				return
			default:
				operation()
				time.Sleep(interval)
			}
		}
	}()
}

func (b *BaseSystem) Close() error {
	b.cancel()

	return nil
}

func (b *BaseSystem) Notify(event SystemEvent) {
	b.manager.notify(b.ctx, event)
}
