package controller

import (
	"calculation_service/model"
	"calculation_service/service"
	"encoding/json"
	fhttp "github.com/valyala/fasthttp"
)

func RequestHandler(ctx *fhttp.RequestCtx) {
	if !isMethodAllowed(ctx) {
		return
	}

	if !isRequestAllowed(ctx) {
		return
	}

	data, err := parseRequestBody(ctx)
	if err != nil {
		return
	}

	processRequest(ctx, data)
}

func isMethodAllowed(ctx *fhttp.RequestCtx) bool {
	if string(ctx.Method()) != fhttp.MethodPost {
		ctx.SetStatusCode(fhttp.StatusMethodNotAllowed)
		ctx.SetContentType("application/json")
		ctx.SetBody([]byte(`{"error": "Method not allowed"}`))
		return false
	}
	return true
}

func isRequestAllowed(ctx *fhttp.RequestCtx) bool {
	if allowed := limiter.Allow(); !allowed {
		ctx.SetStatusCode(fhttp.StatusPaymentRequired)
		ctx.SetContentType("application/json")
		ctx.SetBody([]byte(`{"error": "Too many requests"}`))
		return false
	}
	return true
}

func parseRequestBody(ctx *fhttp.RequestCtx) (model.RequestData, error) {
	var data model.RequestData
	if err := json.Unmarshal(ctx.PostBody(), &data); err != nil {
		ctx.SetStatusCode(fhttp.StatusBadRequest)
		ctx.SetContentType("application/json")
		ctx.SetBody([]byte(`{"error": "Invalid request payload"}`))
		return data, err
	}
	return data, nil
}

func processRequest(ctx *fhttp.RequestCtx, data model.RequestData) {
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

	ctx.SetStatusCode(fhttp.StatusOK)
	ctx.SetContentType("application/json")
	ctx.SetBody(successResponse)
}
