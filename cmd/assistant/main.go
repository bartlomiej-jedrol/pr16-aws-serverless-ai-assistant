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

var (
	service           configuration.Service
	serviceEnvVarName = "SERVICE_NAME"
)

func init() {
	var err error
	service.Name, err = iAws.GetEnvironmentVariable(serviceEnvVarName)
	if err != nil {
		return
	}
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	function := "handler"
	iLog.Info("authorization header", request.Headers["Authorization"], nil, service.Name, function)

	if !api.Authenticate(request.Headers["Authorization"]) {
		return api.BuildResponse(http.StatusForbidden, "unauthorized")
	}

	return api.BuildResponse(http.StatusOK, "response")
}

func main() {
	lambda.Start(handler)
}
