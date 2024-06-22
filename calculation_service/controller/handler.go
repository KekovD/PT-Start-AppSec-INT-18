package controller

import (
	"encoding/json"
	"log"
	"time"

	"calculation_service/model"
	"calculation_service/service"
	fhttp "github.com/valyala/fasthttp"
)

func RequestHandler(ctx *fhttp.RequestCtx) {
	if string(ctx.Method()) != fhttp.MethodPost {
		ctx.SetStatusCode(fhttp.StatusMethodNotAllowed)
		ctx.SetContentType("application/json")
		ctx.SetBody([]byte(`{"error": "Method not allowed"}`))
		return
	}

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
