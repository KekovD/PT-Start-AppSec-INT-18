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
	X       float64 `json:"x"`
	Y       float64 `json:"y"`
	IsEqual string  `json:"is_equal"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type Values [6]float64

func (v *Values) UnmarshalJSON(data []byte) error {
	var tmp []float64
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	if len(tmp) != 6 {
		return fmt.Errorf("incorrect number of values: expected 6, got %d", len(tmp))
	}
	for i := range tmp {
		v[i] = tmp[i]
	}
	return nil
}

type RequestData struct {
	Values Values `json:"values"`
	E      int    `json:"e"`
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
			log.Printf("Error decoding JSON: %v\n", err)
			ctx.SetStatusCode(fhttp.StatusBadRequest)
			ctx.SetContentType("application/json")
			ctx.SetBody([]byte(`{"error": "Invalid request payload"}`))
			return
		}

		X := round((data.Values[0]/data.Values[1])*data.Values[2], data.E)
		Y := round((data.Values[3]/data.Values[4])*data.Values[5], data.E)
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

func round(val float64, precision int) float64 {
	factor := math.Pow(10, float64(precision))
	return math.Round(val*factor) / factor
}
