package main

import (
	"context"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bartlomiej-jedrol/pr16-aws-serverless-ai-assistant/api"
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("INFO: handler - HTTPMethod, %+v", request.HTTPMethod)
	log.Printf("INFO: handler - Headers, %+v", request.Headers)
	// log.Printf("INFO: handler - Body, %+v", request.Body)

	if !api.Authenticate(request.Headers["Authorization"]) {
		return api.BuildResponse(http.StatusForbidden, "unauthorised")
	}

	return api.BuildResponse(http.StatusOK, "response")
}

func main() {
	lambda.Start(handler)
}
