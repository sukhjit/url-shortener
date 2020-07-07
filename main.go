package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	muxAdapter "github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/sukhjit/url-shortener/handler"
)

var (
	awsRegion    string
	isLambda     bool
	muxLambdaSvc *muxAdapter.GorillaMuxAdapter
	router       *mux.Router
	port         string
)

func initEnv() {
	// load env file
	_ = godotenv.Load()

	port = os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	awsRegion = os.Getenv("AWS_REGION")
	if awsRegion == "" {
		awsRegion = "ap-southeast-2"
	}

	isLambda = false
	if os.Getenv("WEB") == "" {
		isLambda = true
	}
}

func main() {
	initEnv()

	router = handler.NewHandler()

	if isLambda {
		lambda.Start(lambdaHandler)
	} else {
		fmt.Println("Started local server on port:", port)

		http.ListenAndServe(":"+port, router)
	}
}

func lambdaHandler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if muxLambdaSvc == nil {
		muxLambdaSvc = muxAdapter.New(router)
	}

	return muxLambdaSvc.Proxy(req)
}
