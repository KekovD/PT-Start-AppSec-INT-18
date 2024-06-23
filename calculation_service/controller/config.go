package controller

import (
	"calculation_service/service"
	"fmt"
	sw "github.com/RussellLuo/slidingwindow"
	"os"
	"strconv"
	"time"
)

var (
	limiter   *sw.Limiter
	datastore service.Datastore
)

func ParseEnvVariables() (string, int, int, time.Duration, time.Duration, time.Duration, int64, string) {
	redisHost := os.Getenv("REDIS_HOST")
	redisPortStr := os.Getenv("REDIS_PORT")
	redisDatabaseStr := os.Getenv("REDIS_DATABASE")
	redisTtlStr := os.Getenv("REDIS_TTL")
	syncIntervalStr := os.Getenv("REDIS_SYNC_INTERVAL")
	intervalStr := os.Getenv("INTERVAL")
	limitStr := os.Getenv("LIMIT")
	syncWindowKey := os.Getenv("SYNC_WINDOW_KEY")

	redisDatabase, err := strconv.Atoi(redisDatabaseStr)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse REDIS_DATABASE: %v", err))
	}

	redisPort, err := strconv.Atoi(redisPortStr)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse REDIS_PORT: %v", err))
	}

	redisTtl, err := time.ParseDuration(redisTtlStr)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse REDIS_TTL: %v", err))
	}

	syncInterval, err := time.ParseDuration(syncIntervalStr)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse REDIS_SYNC_INTERVAL: %v", err))
	}

	interval, err := time.ParseDuration(intervalStr)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse INTERVAL: %v", err))
	}

	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse LIMIT: %v", err))
	}

	return redisHost, redisPort, redisDatabase, redisTtl, syncInterval, interval, limit, syncWindowKey
}
