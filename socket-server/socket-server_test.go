package socketserver

import (
	"testing"
	"time"

	config "github.com/andy2kuo/AndyGameServerGo/cfg"
	"github.com/andy2kuo/AndyGameServerGo/database"
	"github.com/andy2kuo/AndyGameServerGo/logger"
)

type TestDatabaseSetting struct {
	Redis database.RedisConnSetting
	Mongo database.MongoConnSetting
}

func (TestDatabaseSetting) Name() string {
	return "TestDB"
}

func TestForTest(t *testing.T) {
	t.Log(123)
}

func TestSocketServer(t *testing.T) {
	return
	startTime := time.Now().UnixNano()
	_config := &TestDatabaseSetting{}
	config.GetConfig(_config)

	logger := logger.NewLogger("test", "local-test", 0)
	_mongoConn, err := database.NewMongoConnection("TestServer", _config.Mongo)
	if err != nil {
		t.Error(err)
		return
	}
	_redisConn, err := database.NewRedisConnection(_config.Redis)
	if err != nil {
		t.Error(err)
		return
	}

	server, err := NewServer(logger, _mongoConn, _redisConn)
	if err != nil {
		t.Error(err)
		return
	}

	go server.Start()
	callback := time.After(time.Minute)
	select {
	case <-callback:
		server.close()
		break
	default:
		time.Sleep(time.Microsecond)
	}

	t.Log("test over")
	t.Log((time.Now().UnixNano() - startTime) / int64(time.Second))
}
