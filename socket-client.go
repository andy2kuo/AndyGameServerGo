package socketserver

import (
	"ak-project-server/common/errcode"
	"ak-project-server/common/request"
	commandCode "ak-project-server/common/request/command"
	operationCode "ak-project-server/common/request/operation"
	"ak-project-server/logger"
	"errors"
	"io"
	"sync"

	"fmt"
	"net"
	"time"
)

type ClientEvent func(*SocketClient)
var ErrConnectionNull error = errors.New("connection empty")
var ErrClientHadLogin error = errors.New("client is login")

// 客戶端
type SocketClient struct {
	sync.RWMutex

	id              int           // 客戶端編號
	connectTime     time.Time     // 連線時間
	lastConnectTime time.Time     // 最後連線時間
	time_out        time.Duration // 超時時間
	connection      *net.TCPConn  // 客戶端連接口
	conn_ctx        *ConnContext  // 123
	logger          *logger.Logger
	packer          *Packer
	server          *SocketServer

	isLogin bool // 是否登入
}

// 開始客戶端進程
func (client *SocketClient) StartProcess() {
	client.connection.SetKeepAlive(true)
	client.connection.SetKeepAlivePeriod(client.time_out)

	// Getting the file handle of the socket
	sockFile, sockErr := client.connection.File()
	if sockErr == nil {
		//var err error
		//// got socket file handle. Getting descriptor.
		//fd := int(sockFile.Fd())
		//// 心跳封包發送間隔時間
		//err = syscall.SetsockoptInt(fd, syscall.IPPROTO_TCP, syscall.TCP_KEEPCNT, 3)
		//if err != nil {
		//	client.logger.Warn("on setting keepalive probe count", err.Error())
		//}
		//// 重試間隔時間
		//err = syscall.SetsockoptInt(fd, syscall.IPPROTO_TCP, syscall.TCP_KEEPINTVL, 5)
		//if err != nil {
		//	client.logger.Warn("on setting keepalive retry interval", err.Error())
		//}
		// 最後一定要關閉此Socket連線檔案，此關閉不會影響連線
		sockFile.Close()
	} else {
		client.logger.Warn("on setting socket keepalive", sockErr.Error())
	}

	// 接收封包
	go func() {
		// 緩衝接收區
		buffer := make([]byte, 2048)

		for client.connection != nil {
			time.Sleep(time.Millisecond)

			nowTime := time.Now().UTC()
			_getDataLength, readErr := client.connection.Read(buffer)
			if readErr == nil {
				if _getDataLength > 0 {
					client.lastConnectTime = nowTime
					err := client.packer.Add(buffer[:_getDataLength])
					if err != nil {
						client.logger.Error(fmt.Sprintf("Socket add buffer fail. %v", err.Error()))
					}
				}
			} else {
				if errors.Is(readErr, io.EOF) {
					client.Close(errcode.DISCONNECT_BY_CLIENT_STOP)
					break
				}

				client.logger.Error("Socket read fail. %v", readErr.Error())
				continue
			}

			for client.packer.Done() {
				req := client.packer.Get()
				go client.server.OperationResponse(client, req)
			}
		}
	}()

	go func() {
		for {
			time.Sleep(client.time_out)

			if client.connection != nil {
				if time.Now().UTC().Sub(client.lastConnectTime) > client.time_out {
					client.logger.Warn("Client time out")
					client.Close(errcode.DISCONNECT_BY_CLIENT_TIME_OUT)
					break
				}
			} else {
				break
			}
		}
	}()
}

// 關閉客戶端連線
func (client *SocketClient) Close(errCode errcode.ErrorCode) {
	defer func() {
		if client.connection != nil {
			client.connection.Close()
			client.connection = nil
			client.logger.Info("Client Close. Reason:", errCode.Err())
		}
	}()

	//
	if client.isLogin {
		client.SetLogout()
	}
}

// 發送封包
func (client *SocketClient) Send(data []byte) error {
	client.Lock()
	defer client.Unlock()

	if client.connection != nil {
		_, err := client.connection.Write(data)
		if err != nil {
			if !errors.Is(err, io.EOF) {
				client.logger.Error(fmt.Sprintf("Client send data fail, error => %v", err.Error()))
				return err
			} else {
				client.logger.Error("Client send data fail, connection error => EOF")
				return err
			}
		}

		return nil
	} else {
		client.logger.Error("Client send data fail, connection error => Empty")
		return ErrConnectionNull
	}
}

// 發送封包回覆
func (client *SocketClient) SendResponse(opCode operationCode.OperationCode, cmdCode commandCode.CommandCode, rep request.CmdData) error {
	data, err := client.packer.PackData(opCode, cmdCode, rep)
	if err == nil {
		client.Send(data)
	} else {
		client.logger.Error(fmt.Sprintf("Client pack reponse fail, error => %v", err.Error()))
	}

	return err
}

// 設定此連線玩家登入資料
func (client *SocketClient) SetLogin() error {
	client.Lock()
	defer client.Unlock()

	if !client.isLogin {
		
		client.isLogin = true
	} else {
		return ErrClientHadLogin
	}

	return nil
}

func (client *SocketClient) UpdateLoginInfo() error {
	client.Lock()
	defer client.Unlock()

	if !client.isLogin {
		

		client.isLogin = true
	} else {
		return ErrClientHadLogin
	}

	return nil
}

// 清空此連線玩家登入資料
func (client *SocketClient) SetLogout() error {
	client.Lock()
	defer client.Unlock()

	client.isLogin = false

	return nil
}

// 產生新的客戶端
func NewClient(id int, server *SocketServer, ctx *ConnContext, conn *net.TCPConn) *SocketClient {
	new_client := &SocketClient{
		id:              id,
		connectTime:     time.Now().UTC(),
		lastConnectTime: time.Now().UTC(),
		time_out:        time.Second * 30,
		connection:      conn,
		conn_ctx:        ctx,
		packer:          NewPacket(),
		server:          server,
		logger:          ctx.Log(),
	}

	return new_client
}
