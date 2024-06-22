package service

import (
	"fmt"
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

func NewRedisDatastore(
	redisHost string,
	redisPort int,
	redisDatabase int,
	redisTtl time.Duration) *RedisDatastore {
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", redisHost, redisPort),
		DB:   redisDatabase,
	})

	return &RedisDatastore{
		client: client,
		ttl:    redisTtl * time.Second,
	}
}
