package database

import "github.com/go-redis/redis/v8"

func init() {
	tempRedisClient = make(map[string]*redis.Client)
}
