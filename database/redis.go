package database

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

const (
	redis_host = "redis_cache"
	redis_port = 6379
)

func ConnectRedis() *redis.Client {

	rdb := redis.NewClient(&redis.Options{Addr: fmt.Sprintf("%s:%d", redis_host, redis_port)})

	return rdb
}
