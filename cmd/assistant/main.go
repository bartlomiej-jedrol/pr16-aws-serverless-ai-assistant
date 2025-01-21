package main

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	iAws "github.com/bartlomiej-jedrol/go-toolkit/aws"
	iLog "github.com/bartlomiej-jedrol/go-toolkit/log"
	"github.com/bartlomiej-jedrol/pr16-aws-serverless-ai-assistant/api"
	"github.com/bartlomiej-jedrol/pr16-aws-serverless-ai-assistant/configuration"
)

func init() {
	serviceName, err := iAws.GetEnvironmentVariable("SERVICE_NAME")
	if err != nil {
		return
	}
	configuration.SetServiceName(serviceName)
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	function := "handler"
	iLog.Info("body", request.Body, nil, configuration.ServiceName, function)

	if !api.Authenticate(request.Headers["Authorization"]) {
		return api.BuildResponse(http.StatusForbidden, "unauthorized")
	}

	assistantRequest, err := api.ParseRequest(request.Body)
	if err != nil {
		return api.BuildResponse(http.StatusInternalServerError, "internal server error")
	}

	if assistantRequest.Text == "parse blood results" {
		return api.BuildResponse(http.StatusOK, "got parse blood results")
	}
	return api.BuildResponse(http.StatusOK, "response")
}

func main() {
	lambda.Start(handler)
}
