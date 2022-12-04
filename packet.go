package socketserver

import (
	"ak-project-server/common/request"
	commandCode "ak-project-server/common/request/command"
	operationCode "ak-project-server/common/request/operation"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"sync"
)

func NewPacket() (p *Packer) {
	p = &Packer{
		buffer:        make([]byte, 0),
		tempRequests:  make(map[int]*request.SocketRequest),
		nowDataLength: 0,
		nowIndex:      0,
		maxIndex:      0,
	}

	return p
}

type Packer struct {
	sync.RWMutex

	buffer       []byte
	tempRequests map[int]*request.SocketRequest

	nowDataLength int32

	nowIndex int
	maxIndex int
}

func (p *Packer) Done() bool {
	p.RLock()
	defer p.RUnlock()

	return len(p.tempRequests) > 0
}

func (p *Packer) Get() (req *request.SocketRequest) {
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
				if err = p.transferData(packet_data); err == nil {
					p.buffer = p.buffer[p.nowDataLength+4:]
					p.nowDataLength = 0
				} else {
					return fmt.Errorf("socket client transfer data fail. Error => %v", err)
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

	opCode := operationCode.OperationCode(data[data_index])
	data_index++

	// Comand Code

	cmdCode := commandCode.CommandCode(data[data_index])
	data_index++

	var reqData request.CmdData
	err = json.Unmarshal(data[data_index:], &reqData)

	req := request.NewSocketRequest(opCode, cmdCode, reqData)
	p.tempRequests[p.maxIndex] = req
	p.maxIndex++

	return err
}

func (p *Packer) PackData(opCode operationCode.OperationCode, cmdCode commandCode.CommandCode, rep request.CmdData) (data []byte, err error) {
	data = make([]byte, 0)
	err = nil
	var totalLength int32 = 2

	var jsonData []byte
	jsonData, err = json.Marshal(rep)
	if err != nil {
		return data, err
	}

	totalLength += int32(len(jsonData))
	var b []byte
	buf := bytes.NewBuffer(b)
	buf.Reset()
	binary.Write(buf, binary.LittleEndian, totalLength)

	data = append(data, buf.Bytes()...)
	data = append(data, byte(opCode))
	data = append(data, byte(cmdCode))
	data = append(data, jsonData...)

	return data, err
}
