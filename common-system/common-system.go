package commonsystem

import "github.com/andy2kuo/AndyGameServerGo/logger"

type SystemCode byte

// 共用系統
type CommonSystem interface {
	GetSystemCode() SystemCode
	Start() error
	OnSystemInit(*logger.Logger) error
	OnServerStart() error
	OnServerClose() error
}
