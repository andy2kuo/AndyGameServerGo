package socketserver

import "github.com/andy2kuo/AndyGameServerGo/logger"

type OperationEventCode byte

// 流程器
type Operation interface {
	GetOperationCode() OperationCode
	Command(*SocketRequest) error
	OnOperationInit(*SocketServer, *logger.Logger) error
	OnClientConnect(*SocketClient) error
	OnClientDisconnect(*SocketClient) error
	OnEventNotify(*SocketClient, OperationEvent) error
	OnServerStart() error
	OnServerClose() error
}

// 產生新的流程事件
func NewOperationEvent(code OperationEventCode) OperationEvent {
	return OperationEvent{
		opEventCode: code,
	}
}

// 流程事件
type OperationEvent struct {
	opEventCode OperationEventCode
	eventData    ReqData
}

// 取得流程事件編號
func (opEvent *OperationEvent) GetEventCode() OperationEventCode {
	return opEvent.opEventCode
}

// 確認資料編號是否存在
func (opEvent *OperationEvent) IsExist(code DataCode) bool {
	_, isExist := opEvent.eventData[code]
	return isExist
}

// 取得資料
func (opEvent *OperationEvent) Get(code DataCode) (interface{}, bool) {
	data, isExist := opEvent.eventData[code]

	return data, isExist
}

// 依照資料編號設置資料
func (opEvent *OperationEvent) Set(code DataCode, data interface{}) {
	opEvent.eventData[code] = data
}

// 設置所有資料
func (opEvent *OperationEvent) SetAll(datas ReqData) error {
	if datas == nil {
		return ErrDataEmpty
	}

	opEvent.eventData = datas
	return nil
}