package socketserver

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/andy2kuo/AndyGameServerGo/logger"
	"github.com/google/uuid"
)

// 伺服器
type SocketServer struct {
	name          string                   // 伺服器名稱
	port          int                      // 連接埠
	listener      *net.TCPListener         // 伺服器監聽端
	client_list   map[string]*SocketClient // 已連接客戶端列表
	operationCmds map[OperationCode]Operation
	logger        *logger.Logger
	ctx           *ConnContext
	cancel        context.CancelFunc
}

// 啟動
func (server *SocketServer) Start() {

	go func() {
		server.logger.Info("Socket Server Start!")

		for server.listener != nil {
			new_conn, err := server.listener.AcceptTCP()
			if err != nil {
				continue
			}

			new_client_id := fmt.Sprintf("%v-%v-%v", time.Now().Format("20060102"), new_conn.RemoteAddr().String(), uuid.New().String())
			new_client := NewClient(new_client_id, server, server.ctx, new_conn)

			_, is_id_exist := server.client_list[new_client_id]
			if is_id_exist {
				server.logger.Warn(fmt.Sprintf("%v => client id repeated!!", new_client_id))
				server.client_list[new_client_id].Close(DISCONNECT_BY_CLIENT_ID_DUPLICATE)
				delete(server.client_list, new_client_id)
			}

			new_client.StartProcess()
			server.client_list[new_client_id] = new_client
		}
	}()

	osNotify := make(chan os.Signal, 1)
	signal.Notify(osNotify, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

Loop:
	for {
		select {
		case signal := <-osNotify:
			server.logger.Warn(fmt.Sprintf("Get os notify. On signal: %v", signal.String()))
			server.Close()

			// 等待五秒鐘再離開，給正在關閉的服務緩衝時間
			time.Sleep(time.Second * 5)
			break Loop
		default:
			continue
		}
	}
}

// 關閉伺服器
func (server *SocketServer) Close() {
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("%v", r)
			server.logger.Error(err.Error())
		}

		// 最後一定要把監聽關閉
		if server.listener != nil {
			server.listener.Close()
			server.listener = nil
		}
	}()

	// 發送停止通知給底下
	if server.cancel != nil {
		server.cancel()
	}
}

func (server *SocketServer) AddOperation(op Operation) {
	_, isExist := server.operationCmds[op.GetOperationCode()]
	if isExist {
		server.logger.Warn("Op: %v Duplicate!", op.GetOperationCode())
	}

	server.operationCmds[op.GetOperationCode()] = op
}

func (server *SocketServer) RunOperation(req *SocketRequest) {
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("%v", r)
			server.logger.Error(fmt.Sprintf("Recover!! Operation error. Op code = %v, Cmd code = %v, error message => %v", req.OperationCode(), req.CommandCode(), err.Error()))
		}
	}()

	op, isExist := server.operationCmds[req.OperationCode()]
	if isExist {
		err := op.Command(req)
		if err != nil {
			server.logger.Error(fmt.Sprintf("Operation error. Op code = %v, Cmd code = %v, error message => %v", req.OperationCode(), req.CommandCode(), err.Error()))
		}
	} else {
		server.logger.Warn(fmt.Sprintf("Operation not exist. Op code = %v", req.OperationCode()))
	}
}

// 產生新的Socket Server
func NewServer(name string, port int, log *logger.Logger) (server *SocketServer, err error) {
	server = &SocketServer{
		name:          name,
		port:          port,
		client_list:   make(map[string]*SocketClient),
		logger:        log,
		operationCmds: make(map[OperationCode]Operation),
	}

	server.ctx = &ConnContext{
		logger: log,
	}
	server.ctx.Context, server.cancel = context.WithCancel(context.TODO())

	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%v", server.port))
	if err != nil {
		return server, err
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return server, err
	}

	server.listener = listener

	return server, nil
}
