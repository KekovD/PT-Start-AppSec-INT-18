package main

import (
	"calculation_service/controller"
	"fmt"
	"github.com/valyala/fasthttp"
	"os"
	"strconv"
)

func main() {
	calcPortStr := os.Getenv("CALCULATION_PORT")
	calcPort, _ := strconv.Atoi(calcPortStr)

	fmt.Println("Server is running on port", calcPort, "...")

	requestHandler := func(ctx *fasthttp.RequestCtx) {
		controller.RequestHandler(ctx)
	}

	if err := fasthttp.ListenAndServe(fmt.Sprintf(":%d", calcPort), requestHandler); err != nil {
		fmt.Printf("Error in ListenAndServe: %s", err)
	}
}
