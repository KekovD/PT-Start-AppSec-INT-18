package service

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

type Datastore interface {
	Add(key string, start, value int64) (int64, error)
	Get(key string, start int64) (int64, error)
}

type RedisDatastore struct {
	client redis.Cmdable
	ttl    time.Duration
}

func NewRedisDatastore() *RedisDatastore {
	redisHost := os.Getenv("REDIS_HOST")
	redisPortStr := os.Getenv("REDIS_PORT")
	redisPort, _ := strconv.Atoi(redisPortStr)

	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", redisHost, redisPort),
		DB:   1,
	})

	return &RedisDatastore{client: client, ttl: 2 * time.Second}
}
