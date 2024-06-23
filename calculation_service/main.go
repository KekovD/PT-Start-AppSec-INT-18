package main

import (
	"calculation_service/controller"
	"fmt"
	"github.com/valyala/fasthttp"
	"os"
	"strconv"
)

func main() {
	controller.Init()

	calcPortStr := os.Getenv("CALCULATION_PORT")
	calcPort, err := strconv.Atoi(calcPortStr)

	if err != nil {
		panic(fmt.Sprintf("CALCULATION_PORT is not a number: %s", calcPortStr))
	}

	fmt.Println("Server is running on port", calcPort, "...")

	requestHandler := func(ctx *fasthttp.RequestCtx) {
		controller.RequestHandler(ctx)
	}

	if err := fasthttp.ListenAndServe(fmt.Sprintf(":%d", calcPort), requestHandler); err != nil {
		fmt.Printf("Error in ListenAndServe: %s", err)
	}
}
