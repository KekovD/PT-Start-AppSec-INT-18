package controller

import (
	"calculation_service/model"
	"calculation_service/service"
	"encoding/json"
	"fmt"
	sw "github.com/RussellLuo/slidingwindow"
	"os"
)

func init() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Initialization error: %v\n", r)
			os.Exit(1)
		}
	}()

	errorResponse, _ = json.Marshal(model.ErrorResponse{Error: "Too many requests"})

	redisHost, redisPort, redisDatabase, redisTtl, syncInterval, interval, limit, syncWindowKey := parseEnvVariables()

	datastore = service.NewRedisDatastore(redisHost, redisPort, redisDatabase, redisTtl)

	limiter, _ = sw.NewLimiter(interval, limit, func() (sw.Window, sw.StopFunc) {
		return sw.NewSyncWindow(syncWindowKey, sw.NewBlockingSynchronizer(datastore, syncInterval))
	})
}
