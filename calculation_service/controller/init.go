package controller

import (
	"calculation_service/model"
	"encoding/json"
	"time"

	"calculation_service/service"
	sw "github.com/RussellLuo/slidingwindow"
)

func init() {
	errorResponse, _ = json.Marshal(model.ErrorResponse{Error: "Too many requests"})

	datastore = service.NewRedisDatastore()
	size := time.Second
	limiter, _ = sw.NewLimiter(size, 5, func() (sw.Window, sw.StopFunc) {
		return sw.NewSyncWindow("test", sw.NewBlockingSynchronizer(datastore, 500*time.Millisecond))
	})
}
