package socketserver

import (
	"testing"
	"time"

	"github.com/andy2kuo/AndyGameServerGo/logger"
)

func TestForTest(t *testing.T) {
	t.Log(123)
}

func TestSocketServer(t *testing.T) {
	startTime := time.Now().UnixNano()

	logger := logger.NewLogger("test", "local-test", 0)
	server, err := NewServer(logger)
	if err != nil {
		t.Error(err)
		return
	}

	go server.Start()
	select {
	case <-time.After(time.Minute):
		server.close()
		break
	default:
		time.Sleep(time.Microsecond)
	}

	t.Log("test over")
	t.Log((time.Now().UnixNano() - startTime) / int64(time.Second))
}
