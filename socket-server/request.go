package socketserver

import "errors"

var ErrDataEmpty error = errors.New("Request data empty")

func NewSocketRequest(client *SocketClient) (request *SocketRequest) {
	request = &SocketRequest{
		opCode:  0,
		cmdCode: 0,
		reqData: make(ReqData),
		client:  client,
	}

	return request
}

// Socket 請求
type SocketRequest struct {
	reqData ReqData
	opCode  OperationCode
	cmdCode CommandCode
	client  *SocketClient
}

// 設置程序編號
func (req *SocketRequest) SetOperationCode(code OperationCode) {
	req.opCode = code
}

// 取得請求程序編號
func (req *SocketRequest) OperationCode() OperationCode {
	return req.opCode
}

// 設置指令編號
func (req *SocketRequest) SetCommandCode(code CommandCode) {
	req.cmdCode = code
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

type OperationCode byte
type CommandCode byte
type DataCode uint16
type ReqData map[DataCode]interface{}
