package socketserver

import (
	"errors"
	"time"
)

var ErrDataEmpty error = errors.New("request data empty")
var ErrClientNotSet error = errors.New("request client null")

func NewSocketRequest(_opCode OperationCode, _cmdCode CommandCode) (request *SocketRequest) {
	request = &SocketRequest{
		uid: ReqUID(time.Now().UnixMilli()),
		opCode:  _opCode,
		cmdCode: _cmdCode,
		reqData: make(ReqData),
	}

	return request
}

// Socket 請求
type SocketRequest struct {
	uid     ReqUID
	reqData ReqData
	opCode  OperationCode
	cmdCode CommandCode
	client  *SocketClient
}

// 取得請求編號
func (req SocketRequest) GetUID() ReqUID {
	return req.uid
}

// 取得請求時間
func (req SocketRequest) GetRequestTime() time.Time {
	return time.UnixMilli(int64(req.uid))
}

// 設置此請求客戶端
func (req *SocketRequest) SetClient(client *SocketClient) {
	req.client = client
}

// 取得請求程序編號
func (req *SocketRequest) OperationCode() OperationCode {
	return req.opCode
}

// 取得請求指令編號
func (req *SocketRequest) CommandCode() CommandCode {
	return req.cmdCode
}

// 確認資料編號是否存在
func (req *SocketRequest) IsExist(code DataCode) bool {
	_, isExist := req.reqData[code]
	return isExist
}

// 取得資料
func (req *SocketRequest) Get(code DataCode) (interface{}, bool) {
	data, isExist := req.reqData[code]

	return data, isExist
}

// 依照資料編號設置資料
func (req *SocketRequest) Set(code DataCode, data interface{}) {
	req.reqData[code] = data
}

// 設置所有資料
func (req *SocketRequest) SetAll(datas ReqData) error {
	if datas == nil {
		return ErrDataEmpty
	}

	req.reqData = datas
	return nil
}

// 回覆此請求
func (req *SocketRequest) Response(reqData ReqData) error {
	if req.client == nil {
		return ErrClientNotSet
	}

	return req.client.Send(req.GetRequestTime(), req.opCode, req.cmdCode, reqData)
}

// 發送資料
func (req *SocketRequest) Send(sendTime time.Time, opCode OperationCode, cmdCode CommandCode, reqData ReqData) error {
	if req.client == nil {
		return ErrClientNotSet
	}

	return req.client.Send(sendTime, opCode, cmdCode, reqData)
}

type OperationCode byte
type CommandCode byte
type DataCode uint16
type ReqUID int64
type ReqData map[DataCode]interface{}
