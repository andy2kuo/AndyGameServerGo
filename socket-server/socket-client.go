package socketserver

import (
	"context"
	"errors"
	"io"
	"sync"

	"fmt"
	"net"
	"time"

	"github.com/andy2kuo/AndyGameServerGo/logger"
)

type ClientEvent func(*SocketClient)
type ClientInfoCode byte

var ErrConnectionNull error = errors.New("connection empty")
var ErrClientHadLogin error = errors.New("client is login")

// 客戶端
type SocketClient struct {
	sync.RWMutex

	id              string        // 客戶端編號
	connectTime     time.Time     // 連線時間
	lastConnectTime time.Time     // 最後連線時間
	time_out        time.Duration // 超時時間
	connection      *net.TCPConn  // 客戶端連接口
	conn_ctx        context.Context
	logger          *logger.Logger
	packer          *Packer
	server          *SocketServer

	customInfo map[ClientInfoCode]interface{}
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
	Loop:
		for client.connection != nil {
			time.Sleep(time.Millisecond)

			select {
			case <-client.conn_ctx.Done():
				break Loop
			default:
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
						client.Close(ErrClientStop)
						break Loop
					}

					client.logger.Error(fmt.Sprintf("Socket read fail. %v", readErr.Error()))
					continue
				}

				for client.packer.Done() {
					req := client.packer.GetWithClient(client)
					go client.server.RunOperation(req)
				}
			}
		}
	}()

	go func() {
	Loop:
		for {
			time.Sleep(client.time_out)

			select {
			case <-client.conn_ctx.Done():
				break Loop
			default:
				if client.connection != nil {
					if time.Now().UTC().Sub(client.lastConnectTime) > client.time_out {
						client.logger.Warn("Client time out")
						client.Close(ErrConnectTimeOut)
						break Loop
					}
				} else {
					break Loop
				}
			}
		}
	}()
}

// 關閉客戶端連線
func (client *SocketClient) Close(err error) {
	defer func() {
		if client.connection != nil {
			client.connection.Close()
			client.connection = nil
			client.logger.Info("Client Close. Reason:", err.Error())
		}
	}()
}

// 發送封包
func (client *SocketClient) Send(opCode OperationCode, cmdCode CommandCode, reqData ReqData) error {
	client.Lock()
	defer client.Unlock()

	byteData, err := client.packer.PackData(opCode, cmdCode, reqData)

	if err != nil {
		return err
	}

	if client.connection != nil {
		_, err := client.connection.Write(byteData)
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

// 設定自訂資料
func (client *SocketClient) Set(code ClientInfoCode, data interface{}) {
	client.Lock()
	defer client.Unlock()

	client.customInfo[code] = data
}

// 取得自訂資料
func (client *SocketClient) Get(code ClientInfoCode) (data interface{}) {
	client.RLock()
	defer client.RUnlock()

	data, isExist := client.customInfo[code]

	if isExist {
		return data
	}

	return nil
}

// 移除指定自訂資料
func (client *SocketClient) Clear(code ClientInfoCode) {
	client.Lock()
	defer client.Unlock()

	_, isExist := client.customInfo[code]

	if isExist {
		delete(client.customInfo, code)
	}
}

// 清空自訂資料
func (client *SocketClient) ClearAll() {
	client.Lock()
	defer client.Unlock()

	client.customInfo = make(map[ClientInfoCode]interface{})
}

// 產生新的客戶端
func NewClient(id string, server *SocketServer, ctx context.Context, conn *net.TCPConn) *SocketClient {
	new_client := &SocketClient{
		id:              id,
		connectTime:     time.Now().UTC(),
		lastConnectTime: time.Now().UTC(),
		time_out:        time.Second * time.Duration(server.AppSetting.Server.TimeOut),
		connection:      conn,
		conn_ctx:        ctx,
		server:          server,
		logger:          server.logger,
		customInfo:      make(map[ClientInfoCode]interface{}),
	}

	new_client.packer = NewPacket(new_client)

	return new_client
}
