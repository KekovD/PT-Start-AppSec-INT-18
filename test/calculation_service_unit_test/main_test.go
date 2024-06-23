package main_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/valyala/fasthttp"
	"os"
	"strconv"
	"testing"
)

type MockController struct {
	mock.Mock
}

func (m *MockController) RequestHandler(ctx *fasthttp.RequestCtx) {
	m.Called(ctx)
}

func TestMainRun(t *testing.T) {
	calcPortStr := "8080"
	os.Setenv("CALCULATION_PORT", calcPortStr)
	calcPort, err := strconv.Atoi(calcPortStr)

	assert.NoError(t, err)
	assert.Equal(t, 8080, calcPort)

	mockController := new(MockController)

	requestHandler := func(ctx *fasthttp.RequestCtx) {
		mockController.RequestHandler(ctx)
	}

	go func() {
		err := fasthttp.ListenAndServe(fmt.Sprintf(":%d", calcPort), requestHandler)
		assert.NoError(t, err)
	}()

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()

	req.SetRequestURI(fmt.Sprintf("http://localhost:%d", calcPort))

	mockController.On("RequestHandler", mock.AnythingOfType("*fasthttp.RequestCtx")).Run(func(args mock.Arguments) {
		ctx := args.Get(0).(*fasthttp.RequestCtx)
		ctx.SetStatusCode(fasthttp.StatusOK)
	}).Once()

	err = fasthttp.Do(req, resp)
	assert.NoError(t, err)
	assert.Equal(t, fasthttp.StatusOK, resp.StatusCode())

	fasthttp.ReleaseRequest(req)
	fasthttp.ReleaseResponse(resp)

	mockController.AssertExpectations(t)
}
