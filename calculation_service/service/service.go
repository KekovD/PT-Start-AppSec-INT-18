package service

import (
	"calculation_service/model"
	"fmt"
	"math"
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
	if err == redis.Nil {
		return 0, nil
	}
	return value, err
}

func (d *RedisDatastore) fullKey(key string, start int64) string {
	return fmt.Sprintf("%s@%d", key, start)
}

func CalculateX(data model.RequestData) float64 {
	return round((data.Values[0]/data.Values[1])*data.Values[2], data.E)
}

func CalculateY(data model.RequestData) float64 {
	return round((data.Values[3]/data.Values[4])*data.Values[5], data.E)
}

func round(val float64, precision int) float64 {
	factor := math.Pow(10, float64(precision))
	return math.Round(val*factor) / factor
}
