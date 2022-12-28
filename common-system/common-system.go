package commonsystem

import (
	"github.com/andy2kuo/AndyGameServerGo/database"
	"github.com/andy2kuo/AndyGameServerGo/logger"
)

type SystemCode byte

// 共用系統
type ICommonSystem interface {
	GetSystemCode() SystemCode
	Start() error
	OnSystemInit(*logger.Logger, *database.MongoConnection, *database.RedisConnection) error
	OnServerStart() error
	OnServerClose() error
}
