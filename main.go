package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	muxAdapter "github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/sukhjit/url-shortener/handler"
)

var (
	awsRegion     string
	dynamoDBTable string
	isLocal       bool
	muxLambdaSvc  *muxAdapter.GorillaMuxAdapter
	router        *mux.Router
	port          string
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
		awsRegion = endpoints.ApSoutheast2RegionID
	}

	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") == "" {
		isLocal = true
	}

	dynamoDBTable = os.Getenv("DYNAMO_DB_TABLE")
}

func main() {
	initEnv()

	router = handler.NewHandler(isLocal, awsRegion, dynamoDBTable)

	if !isLocal {
		lambda.Start(lambdaHandler)
		return
	}

	fmt.Println("Started local server on port:", port)
	log.Fatal(http.ListenAndServe(":"+port, router)) // nolint: gosec
}

// nolint: gocritic
func lambdaHandler(
	_ context.Context,
	req events.APIGatewayProxyRequest,
) (events.APIGatewayProxyResponse, error) {
	if muxLambdaSvc == nil {
		muxLambdaSvc = muxAdapter.New(router)
	}

	return muxLambdaSvc.Proxy(req)
}
