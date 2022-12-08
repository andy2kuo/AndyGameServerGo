package socketserver

import "github.com/andy2kuo/AndyGameServerGo/logger"

type SystemCode byte

// 共用系統
type SubSystem interface {
	GetSystemCode() SystemCode
	Start() error
	OnSystemInit(*SocketServer, *logger.Logger) error
	OnServerStart() error
	OnServerClose() error
	OnClientConnect(*SocketClient) error
	OnClientLogin(*SocketClient) error
	OnClientDisconnect(*SocketClient) error
	OnClientLogout(*SocketClient) error
}
