package socketserver

import (
	"testing"
)

func TestPack(t *testing.T) {

	req := NewSocketRequest(OperationCode(0), CommandCode(0))
	req.SetAll(ReqData{
		DataCode(0): "123",
		DataCode(1): 321,
	})

	packer := NewPacket(nil)
	bd, _ := packer.PackRequest(req)
	t.Log(bd)

	t.Log()
}
