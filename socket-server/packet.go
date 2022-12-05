package socketserver

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"sync"
)

func NewPacket(client *SocketClient) (p *Packer) {
	p = &Packer{
		buffer:        make([]byte, 0),
		tempRequests:  make(map[int]*SocketRequest),
		nowDataLength: 0,
		nowIndex:      0,
		maxIndex:      0,
	}

	return p
}

type Packer struct {
	sync.RWMutex

	buffer       []byte
	tempRequests map[int]*SocketRequest

	nowDataLength int32

	nowIndex int
	maxIndex int
}

func (p *Packer) Done() bool {
	p.RLock()
	defer p.RUnlock()

	return len(p.tempRequests) > 0
}

func (p *Packer) GetWithClient(client *SocketClient) (req *SocketRequest) {
	p.RLock()
	defer p.RUnlock()

	if len(p.tempRequests) > 0 {
		req = p.tempRequests[p.nowIndex]
		delete(p.tempRequests, p.nowIndex)
		p.nowIndex++
	}

	req.SetClient(client)
	return req
}

func (p *Packer) Get(client *SocketClient) (req *SocketRequest) {
	p.RLock()
	defer p.RUnlock()

	if len(p.tempRequests) > 0 {
		req = p.tempRequests[p.nowIndex]
		delete(p.tempRequests, p.nowIndex)
		p.nowIndex++
	}

	return req
}

func (p *Packer) Add(data []byte) (err error) {
	p.Lock()
	defer p.Unlock()

	err = nil
	if len(data) <= 0 {
		return err
	}

	p.buffer = append(p.buffer, data...)

	for {
		if p.nowDataLength <= 0 {
			if len(p.buffer) >= 4 {
				err := binary.Read(bytes.NewBuffer(p.buffer[0:4]), binary.LittleEndian, &p.nowDataLength)
				if err != nil {
					return err
				}
			} else {
				break
			}
		}

		if p.nowDataLength > 0 {
			if len(p.buffer) >= int(p.nowDataLength) {
				packet_data := p.buffer[4 : p.nowDataLength+4]

				err = p.transferData(packet_data)

				p.buffer = p.buffer[p.nowDataLength+4:]
				p.nowDataLength = 0

				if err != nil {
					err = fmt.Errorf("socket client transfer data fail. Error => %v", err)
					return err
				}
			} else {
				break
			}
		}
	}

	return err
}

func (p *Packer) transferData(data []byte) (err error) {
	err = nil

	var data_index int = 0

	// Operation Code

	opCode := OperationCode(data[data_index])
	data_index++

	// Comand Code

	cmdCode := CommandCode(data[data_index])
	data_index++

	var reqData ReqData
	err = json.Unmarshal(data[data_index:], &reqData)

	req := NewSocketRequest()
	req.SetOperationCode(opCode)
	req.SetCommandCode(cmdCode)
	req.SetAll(reqData)
	p.tempRequests[p.maxIndex] = req
	p.maxIndex++

	return err
}

// 打包檔案
func (p *Packer) PackData(opCode OperationCode, cmdCode CommandCode, reqData ReqData) (byteData []byte, err error) {
	byteData = make([]byte, 0)
	err = nil
	var totalLength int32 = 2

	var jsonData []byte
	jsonData, err = json.Marshal(reqData)
	if err != nil {
		return byteData, err
	}

	totalLength += int32(len(jsonData))
	var b []byte
	buf := bytes.NewBuffer(b)
	buf.Reset()
	binary.Write(buf, binary.LittleEndian, totalLength)

	byteData = append(byteData, buf.Bytes()...)
	byteData = append(byteData, byte(opCode))
	byteData = append(byteData, byte(cmdCode))
	byteData = append(byteData, jsonData...)

	return byteData, err
}
