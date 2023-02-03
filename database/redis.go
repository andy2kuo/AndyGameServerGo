package database

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/go-redis/redis/v8"
)

var ErrDuplicateRedisConn error = errors.New("duplicate redis connection")
var ErrRedisConnNotExist error = errors.New("redis connection not exist")

type RedisConnSetting struct {
	Name         string `default:"db name"`
	Address      string `default:"127.0.0.1"`
	Port         int    `default:"6379"`
	DB           int    `default:"0"`
	PoolSize     int    `default:"10"`
	MaxRetries   int    `default:"3"`
	DialTimeout  int    `default:"5"`
	ReadTimeout  int    `default:"3"`
	WriteTimeout int    `default:"3"`
	PoolTimeout  int    `default:"5"`
	MinIdleConns int    `default:"10"`
	Comment      string `default:"Comment for db"`
}

func NewRedisConnection(_redisSetting interface{}) (_redisConn *RedisConnection, err error) {
	_redisConn = &RedisConnection{
		tempRedisClient: make(map[string]*redis.Client),
	}

	var setting_value reflect.Value = reflect.ValueOf(_redisSetting)
	if setting_value.Kind() == reflect.Pointer {
		setting_value = setting_value.Elem()
	}

	for i := 0; i < setting_value.NumField(); i++ {
		var connSetting RedisConnSetting
		if setting_value.Field(i).CanConvert(reflect.TypeOf(RedisConnSetting{})) {
			connSetting = setting_value.Field(i).Interface().(RedisConnSetting)

			err = _redisConn.addConnection(connSetting)
			if err != nil {
				return _redisConn, err
			}
		}

	}

	return _redisConn, nil
}

type RedisConnection struct {
	tempRedisClient map[string]*redis.Client
}

func (conn *RedisConnection) addConnection(_config RedisConnSetting) error {
	_, isExist := conn.tempRedisClient[_config.Name]
	if isExist {
		return ErrDuplicateRedisConn
	}

	cli := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%v:%v", _config.Address, _config.Port),
		DB:           _config.DB,
		PoolSize:     _config.PoolSize,
		MaxRetries:   _config.MaxRetries,
		DialTimeout:  time.Second * time.Duration(_config.DialTimeout),
		ReadTimeout:  time.Second * time.Duration(_config.ReadTimeout),
		WriteTimeout: time.Second * time.Duration(_config.WriteTimeout),
		PoolTimeout:  time.Second * time.Duration(_config.PoolTimeout),
		MinIdleConns: _config.MinIdleConns,
	})

	_, err := cli.Ping(context.Background()).Result()
	if err != nil {
		return err
	}

	conn.tempRedisClient[_config.Name] = cli
	return nil
}

func (conn *RedisConnection) GetRedis(name string) (cli *redis.Client, err error) {
	cli, isExist := conn.tempRedisClient[name]
	if !isExist {
		return nil, ErrRedisConnNotExist
	}

	return cli, nil
}

func (conn *RedisConnection) Disconnect(name string) (err error) {
	cli, isExist := conn.tempRedisClient[name]
	if !isExist || cli == nil {
		return nil
	}

	delete(conn.tempRedisClient, name)
	return cli.Close()
}

func (conn *RedisConnection) DisconnectAll() (err error) {
	for name := range conn.tempRedisClient {
		err = conn.Disconnect(name)
		if err != nil {
			return fmt.Errorf("redis '%v' close fail. msg => %v", name, err.Error())
		}
	}

	return nil
}
