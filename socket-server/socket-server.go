package socketserver

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	config "github.com/andy2kuo/AndyGameServerGo/cfg"
	commonsystem "github.com/andy2kuo/AndyGameServerGo/common-system"
	"github.com/andy2kuo/AndyGameServerGo/database"
	"github.com/andy2kuo/AndyGameServerGo/logger"
)

// 伺服器
type SocketServer struct {
	listener    *net.TCPListener         // 伺服器監聽端
	client_list map[string]*SocketClient // 已連接客戶端列表
	systems     map[commonsystem.SystemCode]commonsystem.ICommonSystem
	operations  map[OperationCode]IOperation
	logger      *logger.Logger
	ctx         context.Context
	cancel      context.CancelFunc
	serialNum   uint64
	mongoConn   *database.MongoConnection
	redisConn   *database.RedisConnection
	env         string

	AppSetting *AppSetting
}

func (server SocketServer) Environment() string {
	return server.env
}

// 啟動
func (server *SocketServer) Start() {
	var cross__day_time time.Time

	go func() {
		server.logger.Info("Socket Server Start!")
		year, month, day := time.Now().Date()
		cross__day_time = time.Date(year, month, day+1, 0, 0, 0, 0, time.Now().Location())

		if len(server.systems) > 0 {
			for _, sys := range server.systems {
				if err := sys.OnServerStart(); err != nil {
					server.logger.Error(fmt.Sprintf("System Start error. Sys code = %v, error message => %v", sys.GetSystemCode(), err.Error()))
				}
			}
		}

		if len(server.operations) > 0 {
			for _, op := range server.operations {
				if err := op.OnServerStart(); err != nil {
					server.logger.Error(fmt.Sprintf("Operation Start error. Op code = %v, error message => %v", op.GetOperationCode(), err.Error()))
				}
			}
		}

		for server.listener != nil {
			new_conn, err := server.listener.AcceptTCP()
			if err != nil {
				continue
			}

			if time.Now().After(cross__day_time) {
				cross__day_time = cross__day_time.AddDate(0, 0, 1)
				server.serialNum = 0
			}

			new_client_id := fmt.Sprintf("socket-%v-%v-%v", time.Now().Format("20060102"), new_conn.RemoteAddr().String(), server.serialNum)
			new_client := NewClient(new_client_id, server, server.ctx, new_conn)

			_, is_id_exist := server.client_list[new_client_id]
			if is_id_exist {
				server.logger.Warn(fmt.Sprintf("%v => client id repeated!!", new_client_id))
				server.client_list[new_client_id].Close(ErrClientIDDuplicate)
				delete(server.client_list, new_client_id)
			}

			new_client.StartProcess()
			server.client_list[new_client_id] = new_client

			server.serialNum++
		}
	}()

	osNotify := make(chan os.Signal, 1)
	signal.Notify(osNotify, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

Loop:
	for {
		select {
		case signal := <-osNotify:
			server.logger.Warn(fmt.Sprintf("Get os notify. On signal: %v", signal.String()))
			server.close()

			// 等待五秒鐘再離開，給正在關閉的服務緩衝時間
			time.Sleep(time.Second * 5)
			break Loop
		default:
			continue
		}
	}
}

// 關閉伺服器
func (server *SocketServer) close() {
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

	if len(server.systems) > 0 {
		for _, sys := range server.systems {
			if err := sys.OnServerClose(); err != nil {
				server.logger.Error(fmt.Sprintf("System Close error. Sys code = %v, error message => %v", sys.GetSystemCode(), err.Error()))
			}
		}
	}

	if len(server.operations) > 0 {
		for _, op := range server.operations {
			if err := op.OnServerClose(); err != nil {
				server.logger.Error(fmt.Sprintf("Operation Close error. Op code = %v, error message => %v", op.GetOperationCode(), err.Error()))
			}
		}
	}
}

// 加入流程器
func (server *SocketServer) AddOperation(op IOperation) error {
	_, isExist := server.operations[op.GetOperationCode()]
	if isExist {
		server.logger.Warn(fmt.Sprintf("Op: %v Duplicate!", op.GetOperationCode()))
	}

	server.operations[op.GetOperationCode()] = op
	return op.OnOperationInit(server, server.logger)
}

// 加入共用系統
func (server *SocketServer) AddSubSystem(sys commonsystem.ICommonSystem) error {
	_, isExist := server.systems[sys.GetSystemCode()]
	if isExist {
		server.logger.Warn(fmt.Sprintf("Sys: %v Duplicate!", sys.GetSystemCode()))
	}

	server.systems[sys.GetSystemCode()] = sys
	return sys.OnSystemInit(server.logger, server.mongoConn, server.redisConn)
}

// 取得共用系統
func (server *SocketServer) GetCommonSystem(sysCode commonsystem.SystemCode) commonsystem.ICommonSystem {
	sys, isExist := server.systems[sysCode]
	if isExist {
		return sys
	}

	server.logger.Warn(fmt.Sprintf("Sys: Get System %v fail", sysCode))
	return nil
}

// 當有新的客戶端連線進入時
func (server *SocketServer) OnClientConnect(client *SocketClient) {
	if len(server.operations) > 0 {
		for _, op := range server.operations {
			if err := op.OnClientConnect(client); err != nil {
				server.logger.Error(fmt.Sprintf("Operation error on client connect notify. Op code = %v, error message => %v", op.GetOperationCode(), err.Error()))
			}
		}
	}
}

// 當有客戶端斷線離開時
func (server *SocketServer) OnClientDisconnect(client *SocketClient) {
	if len(server.operations) > 0 {
		for _, op := range server.operations {
			if err := op.OnClientDisconnect(client); err != nil {
				server.logger.Error(fmt.Sprintf("Operation error on client disconnect notify. Op code = %v, error message => %v", op.GetOperationCode(), err.Error()))
			}
		}
	}
}

// 當有用戶事件通知時
func (server *SocketServer) OnEventNotify(client *SocketClient, sysEvent OperationEvent) {
	if len(server.operations) > 0 {
		for _, op := range server.operations {
			if err := op.OnEventNotify(client, sysEvent); err != nil {
				server.logger.Error(fmt.Sprintf("Operation error on client event notify. Op code = %v, error message => %v", op.GetOperationCode(), err.Error()))
			}
		}
	}
}

// 執行流程
func (server *SocketServer) RunOperation(req *SocketRequest) {
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("%v", r)
			server.logger.Error(fmt.Sprintf("Recover!! Operation error. Op code = %v, Cmd code = %v, error message => %v", req.OperationCode(), req.CommandCode(), err.Error()))
		}
	}()

	op, isExist := server.operations[req.OperationCode()]
	if isExist {
		timeoutChannel := make(chan bool)

		go func() {
			err := op.Command(req)

			if err != nil {
				server.logger.Error(fmt.Sprintf("Operation error. Op code = %v, Cmd code = %v, error message => %v", req.OperationCode(), req.CommandCode(), err.Error()))
			}

			timeoutChannel <- true
		}()

		select {
		case <-timeoutChannel:
			// 正常執行
			break
		case <-time.After(time.Duration(server.AppSetting.Operation.RunMaxTime) * time.Second):
			// 流程執行超時
			server.logger.Error(fmt.Sprintf("Operation time out for %v secs. Op code = %v, Cmd code = %v", server.AppSetting.Operation.RunMaxTime, req.OperationCode(), req.CommandCode()))
			break
		}

	} else {
		server.logger.Warn(fmt.Sprintf("Operation not exist. Op code = %v", req.OperationCode()))
	}
}

// 產生新的Socket Server
func NewServer(env string, log *logger.Logger, _mongoConn *database.MongoConnection, _redisConn *database.RedisConnection) (server *SocketServer, err error) {
	server = &SocketServer{
		env:         env,
		client_list: make(map[string]*SocketClient),
		logger:      log,
		operations:  make(map[OperationCode]IOperation),
		serialNum:   0,
		mongoConn:   _mongoConn,
		redisConn:   _redisConn,
	}

	var _setting *AppSetting = &AppSetting{}
	_setting_err := config.GetConfig(env, _setting)
	if config.IsCreateNew(_setting_err) {
		server.logger.Info("Create new application setting file")
	}

	server.AppSetting = _setting
	server.ctx, server.cancel = context.WithCancel(context.TODO())

	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%v", server.AppSetting.Server.Port))
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
