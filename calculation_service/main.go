package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"time"

	sw "github.com/RussellLuo/slidingwindow"
	"github.com/go-redis/redis"
	fhttp "github.com/valyala/fasthttp"
)

type SuccessResponse struct {
	X       float64 `json:"X"`
	Y       float64 `json:"Y"`
	IsEqual string  `json:"IsEqual"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type RequestData struct {
	X1 float64 `json:"X1"`
	X2 float64 `json:"X2"`
	X3 float64 `json:"X3"`
	Y1 float64 `json:"Y1"`
	Y2 float64 `json:"Y2"`
	Y3 float64 `json:"Y3"`
	E  int     `json:"E"`
}

type RedisDatastore struct {
	client redis.Cmdable
	ttl    time.Duration
}

func NewRedisDatastore(client redis.Cmdable, ttl time.Duration) *RedisDatastore {
	return &RedisDatastore{client: client, ttl: ttl}
}

func (d *RedisDatastore) fullKey(key string, start int64) string {
	return fmt.Sprintf("%s@%d", key, start)
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

var (
	errorResponse []byte
	limiter       *sw.Limiter
)

func init() {
	errorResponse, _ = json.Marshal(ErrorResponse{Error: "Too many requests"})

	client := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
		DB:   1,
	})

	store := NewRedisDatastore(client, 2*time.Second)

	size := time.Second

	limiter, _ = sw.NewLimiter(size, 5, func() (sw.Window, sw.StopFunc) {
		return sw.NewSyncWindow("test", sw.NewBlockingSynchronizer(store, 500*time.Millisecond))
	})
}

func main() {
	requestHandler := func(ctx *fhttp.RequestCtx) {
		if allowed := limiter.Allow(); !allowed {
			log.Printf("Too many requests at %s\n", time.Now().Format(time.RFC3339Nano))
			ctx.SetStatusCode(fhttp.StatusPaymentRequired)
			ctx.SetContentType("application/json")
			ctx.SetBody(errorResponse)
			return
		}

		var data RequestData
		if err := json.Unmarshal(ctx.PostBody(), &data); err != nil {
			ctx.SetStatusCode(fhttp.StatusBadRequest)
			ctx.SetContentType("application/json")
			errorResponse, _ := json.Marshal(ErrorResponse{Error: "Invalid request payload"})
			ctx.SetBody(errorResponse)
			return
		}

		X := math.Round((data.X1/data.X2)*data.X3*math.Pow(10, float64(data.E))) / math.Pow(10, float64(data.E))
		Y := math.Round((data.Y1/data.Y2)*data.Y3*math.Pow(10, float64(data.E))) / math.Pow(10, float64(data.E))
		IsEqual := "F"
		if X == Y {
			IsEqual = "T"
		}

		successResponse, _ := json.Marshal(SuccessResponse{
			X:       X,
			Y:       Y,
			IsEqual: IsEqual,
		})

		log.Printf("Request successful at %s\n", time.Now().Format(time.RFC3339Nano))
		ctx.SetStatusCode(fhttp.StatusOK)
		ctx.SetContentType("application/json")
		ctx.SetBody(successResponse)
	}

	fmt.Println("Server is running on port 8080...")
	if err := fhttp.ListenAndServe(":8080", requestHandler); err != nil {
		fmt.Printf("Error in ListenAndServe: %s", err)
	}
}
