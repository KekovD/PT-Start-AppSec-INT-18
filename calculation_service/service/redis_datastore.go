package service

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
)

func (d *RedisDatastore) Add(key string, start, value int64) (int64, error) {
	k := d.fullKey(key, start)
	c, err := d.client.IncrBy(k, value).Result()

	if err != nil {
		return 0, err
	}

	err = d.client.Expire(k, d.ttl).Err()
	return c, err
}

func (d *RedisDatastore) Get(key string, start int64) (int64, error) {
	k := d.fullKey(key, start)
	value, err := d.client.Get(k).Int64()

	if errors.Is(err, redis.Nil) {
		return 0, nil
	}

	return value, err
}

func (d *RedisDatastore) fullKey(key string, start int64) string {
	return fmt.Sprintf("%s@%d", key, start)
}
