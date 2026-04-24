package database

import (
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

var (
	redis_host = os.Getenv("REDIS_HOST")
	redis_port = os.Getenv("REDIS_PORT")
)

func ConnectRedis() *redis.Client {

	rdb := redis.NewClient(&redis.Options{Addr: fmt.Sprintf("%s:%s", redis_host, redis_port)})

	return rdb
}
