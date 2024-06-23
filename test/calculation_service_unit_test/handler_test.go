package calculation_service_unit

import (
	"calculation_service/controller"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func TestRequestHandler(t *testing.T) {
	var ctx fasthttp.RequestCtx
	ctx.Request.Header.SetMethod(fasthttp.MethodGet)

	controller.RequestHandler(&ctx)

	assert.Equal(t, fasthttp.StatusMethodNotAllowed, ctx.Response.StatusCode())
}
