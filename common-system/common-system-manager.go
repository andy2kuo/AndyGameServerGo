package commonsystem

import (
	"fmt"
	"reflect"

	"github.com/andy2kuo/AndyGameServerGo/database"
	"github.com/andy2kuo/AndyGameServerGo/logger"
	"github.com/andy2kuo/AndyGameServerGo/pubsub"
)

func NewSystemManager(_log *logger.Logger, _mongoConn *database.MongoConnection, _redisConn *database.RedisConnection) *CommonSystemManager {
	return &CommonSystemManager{
		Hub:       pubsub.NewHub(),
		logger:    _log,
		mongoConn: _mongoConn,
		redisConn: _redisConn,
		systems:   make(map[SystemCode]ICommonSystem),
	}
}

type CommonSystemManager struct {
	*pubsub.Hub
	logger    *logger.Logger
	mongoConn *database.MongoConnection
	redisConn *database.RedisConnection

	systems map[SystemCode]ICommonSystem
}

func (m *CommonSystemManager) AddSystem(sys ICommonSystem) error {
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
	for _, sys := range m.systems {
		err := sys.Close()
		if err != nil {
			m.logger.Error(fmt.Sprintf("Sys: %v fail on close", sys.GetSystemCode()))
		}
	}
}
