package socketserver

import "github.com/andy2kuo/AndyGameServerGo/logger"

type SystemCode byte
type SystemEventCode byte

// 共用系統
type CommonSystem interface {
	GetSystemCode() SystemCode
	Start() error
	OnSystemInit(*SocketServer, *logger.Logger) error
	OnServerStart() error
	OnServerClose() error
	OnClientConnect(*SocketClient) error
	OnClientDisconnect(*SocketClient) error
	OnEventNotify(*SocketClient, SystemEvent) error
}

// 產生新的系統事件
func NewSystemEvent(code SystemEventCode) SystemEvent {
	return SystemEvent{
		sysEventCode: code,
	}
}

// 共用系統事件
type SystemEvent struct {
	sysEventCode SystemEventCode
	eventData    ReqData
}

// 取得系統事件編號
func (sysEvent *SystemEvent) GetEventCode() SystemEventCode {
	return sysEvent.sysEventCode
}

// 確認資料編號是否存在
func (sysEvent *SystemEvent) IsExist(code DataCode) bool {
	_, isExist := sysEvent.eventData[code]
	return isExist
}

// 取得資料
func (sysEvent *SystemEvent) Get(code DataCode) (interface{}, bool) {
	data, isExist := sysEvent.eventData[code]

	return data, isExist
}

// 依照資料編號設置資料
func (sysEvent *SystemEvent) Set(code DataCode, data interface{}) {
	sysEvent.eventData[code] = data
}

// 設置所有資料
func (sysEvent *SystemEvent) SetAll(datas ReqData) error {
	if datas == nil {
		return ErrDataEmpty
	}

	sysEvent.eventData = datas
	return nil
}
