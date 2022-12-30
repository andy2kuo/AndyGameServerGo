package socketserver

import (
	"testing"
)

func TestPack(t *testing.T) {

	req := NewSocketRequest(OperationCode(2), CommandCode(98))
	req.SetAll(ReqData{
		DataCode(0): "123",
		DataCode(1): 321,
	})
	t.Log(req)
	t.Log(req.GetUID())
	t.Log(req.GetRequestTime())

	packer := NewPacket(nil)
	bd, _ := packer.PackRequest(req)
	t.Log(bd)

	err := packer.Add(bd)
	if err != nil {
		t.Error(err)
	}

	if packer.Done() {
		_req := packer.Get()
		t.Log(_req)
		t.Log(_req.GetUID())
		t.Log(_req.GetRequestTime())
	} else {
		t.Error("Unpack fail")
	}
}
