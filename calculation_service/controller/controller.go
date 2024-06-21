package controller

import (
	"encoding/json"
	"log"
	"time"

	"calculation_service/model"
	"calculation_service/service"
	sw "github.com/RussellLuo/slidingwindow"
	fhttp "github.com/valyala/fasthttp"
)

var (
	errorResponse []byte
	limiter       *sw.Limiter
	datastore     service.Datastore
)

func init() {
	errorResponse, _ = json.Marshal(model.ErrorResponse{Error: "Too many requests"})

	// Инициализация Redis и sliding window limiter
	datastore = service.NewRedisDatastore()
	size := time.Second
	limiter, _ = sw.NewLimiter(size, 5, func() (sw.Window, sw.StopFunc) {
		return sw.NewSyncWindow("test", sw.NewBlockingSynchronizer(datastore, 500*time.Millisecond))
	})
}

func RequestHandler(ctx *fhttp.RequestCtx) {
	if allowed := limiter.Allow(); !allowed {
		log.Printf("Too many requests at %s\n", time.Now().Format(time.RFC3339Nano))
		ctx.SetStatusCode(fhttp.StatusPaymentRequired)
		ctx.SetContentType("application/json")
		ctx.SetBody(errorResponse)
		return
	}

	var data model.RequestData
	if err := json.Unmarshal(ctx.PostBody(), &data); err != nil {
		log.Printf("Error decoding JSON: %v\n", err)
		ctx.SetStatusCode(fhttp.StatusBadRequest)
		ctx.SetContentType("application/json")
		ctx.SetBody([]byte(`{"error": "Invalid request payload"}`))
		return
	}

	X := service.CalculateX(data)
	Y := service.CalculateY(data)
	IsEqual := "F"
	if X == Y {
		IsEqual = "T"
	}

	successResponse, _ := json.Marshal(model.SuccessResponse{
		X:       X,
		Y:       Y,
		IsEqual: IsEqual,
	})

	log.Printf("Request successful at %s\n", time.Now().Format(time.RFC3339Nano))
	ctx.SetStatusCode(fhttp.StatusOK)
	ctx.SetContentType("application/json")
	ctx.SetBody(successResponse)
}
