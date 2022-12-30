package socketserver

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

func NewPacket(client *SocketClient) (p *Packer) {
	p = &Packer{
		buffer:        make([]byte, 0),
		tempRequests:  make(map[uint8]*SocketRequest),
		nowDataLength: 0,
		nowIndex:      0,
		maxIndex:      0,
	}

	return p
}

type Packer struct {
	sync.RWMutex

	buffer       []byte
	tempRequests map[uint8]*SocketRequest

	nowDataLength int32

	nowIndex uint8
	maxIndex uint8
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

func (p *Packer) Get() (req *SocketRequest) {
	p.RLock()
	defer p.RUnlock()

	if len(p.tempRequests) > 0 {
		req = p.tempRequests[p.nowIndex]
		delete(p.tempRequests, p.nowIndex)
		if p.nowIndex < 255 {
			p.nowIndex++
		} else {
			p.nowIndex = 0
		}
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

	data_buff := bytes.NewBuffer(data)

	var uid ReqUID
	bd_uid := data_buff.Next(8)
	if len(bd_uid) < 8 {
		return fmt.Errorf("request unpack fail. invalid uid")
	}
	binary.Read(bytes.NewBuffer(bd_uid), binary.LittleEndian, &uid)

	// Operation Code
	var opCode OperationCode
	bd_opCode := data_buff.Next(1)
	if len(bd_uid) < 1 {
		return fmt.Errorf("request unpack fail. invalid op code")
	}
	binary.Read(bytes.NewBuffer(bd_opCode), binary.LittleEndian, &opCode)

	// Comand Code
	var cmdCode CommandCode
	bd_cmdCode := data_buff.Next(1)
	if len(bd_uid) < 1 {
		return fmt.Errorf("request unpack fail. invalid command code")
	}
	binary.Read(bytes.NewBuffer(bd_cmdCode), binary.LittleEndian, &cmdCode)

	var reqData ReqData
	err = json.Unmarshal(data_buff.Next(data_buff.Len()), &reqData)

	req := NewSocketRequest(opCode, cmdCode)
	req.uid = uid
	req.SetAll(reqData)
	p.tempRequests[p.maxIndex] = req
	if p.maxIndex < 255 {
		p.maxIndex++
	} else {
		p.maxIndex = 0
	}

	return err
}

// 打包檔案
func (p *Packer) PackData(reqTime time.Time, opCode OperationCode, cmdCode CommandCode, reqData ReqData) (byteData []byte, err error) {
	byteData = make([]byte, 0)
	err = nil

	var jsonData []byte
	jsonData, err = json.Marshal(reqData)
	if err != nil {
		return byteData, err
	}

	var totalLength int32 = int32(8 + 1 + 1 + len(jsonData))

	buf := bytes.NewBuffer(make([]byte, 0))
	binary.Write(buf, binary.LittleEndian, totalLength)
	binary.Write(buf, binary.LittleEndian, reqTime.UnixMilli())
	binary.Write(buf, binary.LittleEndian, opCode)
	binary.Write(buf, binary.LittleEndian, cmdCode)
	binary.Write(buf, binary.LittleEndian, jsonData)

	return buf.Bytes(), err
}

func (p *Packer) PackRequest(req *SocketRequest) ([]byte, error) {
	return p.PackData(req.GetRequestTime(), req.opCode, req.cmdCode, req.reqData)
}
