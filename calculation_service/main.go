package main

import (
	"calculation_service/controller"
	"fmt"
	"github.com/valyala/fasthttp"
	"os"
	"strconv"
)

func main() {
	goappPortStr := os.Getenv("GOAPP_PORT")
	goappPort, _ := strconv.Atoi(goappPortStr)

	fmt.Println("Server is running on port", goappPort, "...")

	requestHandler := func(ctx *fasthttp.RequestCtx) {
		controller.RequestHandler(ctx)
	}

	if err := fasthttp.ListenAndServe(fmt.Sprintf(":%d", goappPort), requestHandler); err != nil {
		fmt.Printf("Error in ListenAndServe: %s", err)
	}
}
