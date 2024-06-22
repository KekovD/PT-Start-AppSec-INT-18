package controller

import (
	"calculation_service/service"
	"fmt"
	sw "github.com/RussellLuo/slidingwindow"
	"os"
)

func init() {
	const ttlMultiplier = 2

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Initialization error: %v\n", r)
			os.Exit(1)
		}
	}()

	redisHost, redisPort, redisDatabase, redisTtl, syncInterval, interval, limit, syncWindowKey := parseEnvVariables()

	if redisTtl < interval*ttlMultiplier {
		panic(fmt.Sprintf("redisTtl (%v) must be at least twice as large as interval (%v)", redisTtl, interval))
	}

	datastore = service.NewRedisDatastore(redisHost, redisPort, redisDatabase, redisTtl)

	limiter, _ = sw.NewLimiter(interval, limit, func() (sw.Window, sw.StopFunc) {
		return sw.NewSyncWindow(syncWindowKey, sw.NewBlockingSynchronizer(datastore, syncInterval))
	})
}
