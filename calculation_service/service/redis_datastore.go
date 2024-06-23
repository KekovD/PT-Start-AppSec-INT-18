package service

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
)

func (d *RedisDatastore) Add(key string, start, value int64) (int64, error) {
	k := d.fullKey(key, start)

	pipe := d.client.Pipeline()
	incrCmd := pipe.IncrBy(k, value)
	expireCmd := pipe.Expire(k, d.ttl)

	_, err := pipe.Exec()
	if err != nil {
		return 0, err
	}

	return incrCmd.Val(), expireCmd.Err()
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
