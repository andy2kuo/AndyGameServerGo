package socketserver

import "github.com/andy2kuo/AndyGameServerGo/logger"

// 流程器
type Operation interface {
	GetOperationCode() OperationCode
	Command(*SocketRequest) error
	OnOperationInit(*SocketServer, *logger.Logger) error
	OnServerStart() error
	OnServerClose() error
}
