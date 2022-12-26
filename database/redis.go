package database

import (
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

var tempRedisClient map[string]*redis.Client

type RedisConnection struct {
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

func GetRedisClient(_config RedisConnection) (cli *redis.Client, err error) {
	cli, isExist := tempRedisClient[_config.Name]
	if isExist {
		return cli, nil
	}

	cli = redis.NewClient(&redis.Options{
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

	tempRedisClient[_config.Name] = cli

	return cli, nil
}
