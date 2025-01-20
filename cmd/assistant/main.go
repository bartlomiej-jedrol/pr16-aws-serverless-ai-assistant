package main

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	iLog "github.com/bartlomiej-jedrol/go-toolkit/log"
	"github.com/bartlomiej-jedrol/pr16-aws-serverless-ai-assistant/api"
)

var service = "pr16-assistant"

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	service := "handler"
	enpoint := "assistant"
	iLog.Info("authorization header", request.Headers["Authorization"], nil, "pr16-assistant", service, &enpoint, nil)

	if !api.Authenticate(request.Headers["Authorization"]) {
		return api.BuildResponse(http.StatusForbidden, "unauthorized")
	}

	return api.BuildResponse(http.StatusOK, "response")
}

func main() {
	lambda.Start(handler)
}
