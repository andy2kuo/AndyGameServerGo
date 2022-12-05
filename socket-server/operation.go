package socketserver

type Operation interface {
	GetOperationCode() OperationCode
	Command(*SocketRequest) error
	OnOperationInit() error
	OnServerInit(*SocketServer) error
	OnClose() error
}
