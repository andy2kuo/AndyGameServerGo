package commonsystem

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"github.com/andy2kuo/AndyGameServerGo/database"
	"github.com/andy2kuo/AndyGameServerGo/logger"
)

func NewSystemManager(_log *logger.Logger, _mongoConn *database.MongoConnection, _redisConn *database.RedisConnection) *CommonSystemManager {
	return &CommonSystemManager{
		logger:    _log,
		mongoConn: _mongoConn,
		redisConn: _redisConn,
		systems:   make(map[SystemCode]ICommonSystem),
	}
}

type CommonSystemManager struct {
	sync.Mutex
	logger    *logger.Logger
	mongoConn *database.MongoConnection
	redisConn *database.RedisConnection

	systems map[SystemCode]ICommonSystem
}

func (m *CommonSystemManager) AddSystem(sys ICommonSystem) error {
	m.Lock()
	defer m.Unlock()

	if reflect.TypeOf(sys).Kind() != reflect.Ptr {
		return fmt.Errorf("CommonSystemManager: AddSystem should use pointer")
	}

	_, isExist := m.systems[sys.GetSystemCode()]
	if isExist {
		m.logger.Warn(fmt.Sprintf("CommonSystemManager: AddSystem %v Duplicate", sys.GetSystemCode()))
		return nil
	}

	m.systems[sys.GetSystemCode()] = sys
	return sys.Init(m, m.logger, m.mongoConn, m.redisConn)
}

func (m *CommonSystemManager) GetSystem(sysCode SystemCode) ICommonSystem {
	sys, isExist := m.systems[sysCode]
	if !isExist {
		m.logger.Warn(fmt.Sprintf("Sys: %v not found", sysCode))
	}

	return sys
}

func (m *CommonSystemManager) OnServerStart() error {
	for _, sys := range m.systems {
		err := sys.OnServerStart()
		if err != nil {
			m.logger.Error(fmt.Sprintf("Sys: %v fail on server start", sys.GetSystemCode()))
			return err
		}
	}

	return nil
}

func (m *CommonSystemManager) CloseAllSystem() {
	m.Lock()
	defer m.Unlock()

	for _, sys := range m.systems {
		err := sys.Close()
		if err != nil {
			m.logger.Error(fmt.Sprintf("Sys: %v fail on close", sys.GetSystemCode()))
		}
	}
}

func (m *CommonSystemManager) notify(ctx context.Context, event SystemEvent) {
	m.Lock()
	defer m.Unlock()

	select {
	case <-ctx.Done():
		break
	default:
		for _, sys := range m.systems {
			sys.OnSystemEventNotify(event)
		}
	}
}
