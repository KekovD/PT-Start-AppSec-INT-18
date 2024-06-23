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
	if !IsMethodAllowed(ctx) {
		return
	}

	if !IsRequestAllowed(ctx) {
		return
	}

	data, err := ParseRequestBody(ctx)
	if err != nil {
		return
	}

	ProcessRequest(ctx, data)
}

func IsMethodAllowed(ctx *fhttp.RequestCtx) bool {
	if string(ctx.Method()) != fhttp.MethodPost {
		ctx.SetStatusCode(fhttp.StatusMethodNotAllowed)
		ctx.SetContentType("application/json")
		ctx.SetBody([]byte(`{"error": "Method not allowed"}`))
		return false
	}
	return true
}

func IsRequestAllowed(ctx *fhttp.RequestCtx) bool {
	if allowed := limiter.Allow(); !allowed {
		log.Printf("Too many requests at %s\n", time.Now().Format(time.RFC3339Nano))
		ctx.SetStatusCode(fhttp.StatusPaymentRequired)
		ctx.SetContentType("application/json")
		ctx.SetBody([]byte(`{"error": "Too many requests"}`))
		return false
	}
	return true
}

func ParseRequestBody(ctx *fhttp.RequestCtx) (model.RequestData, error) {
	var data model.RequestData
	if err := json.Unmarshal(ctx.PostBody(), &data); err != nil {
		log.Printf("Error decoding JSON: %v\n", err)
		ctx.SetStatusCode(fhttp.StatusBadRequest)
		ctx.SetContentType("application/json")
		ctx.SetBody([]byte(`{"error": "Invalid request payload"}`))
		return data, err
	}
	return data, nil
}

func ProcessRequest(ctx *fhttp.RequestCtx, data model.RequestData) {
	E := data.E

	X := service.CalculateValue(data.Values[0], data.Values[1], data.Values[2], E)
	Y := service.CalculateValue(data.Values[3], data.Values[4], data.Values[5], E)

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
