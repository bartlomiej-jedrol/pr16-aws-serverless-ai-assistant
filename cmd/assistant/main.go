package main

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	iAws "github.com/bartlomiej-jedrol/go-toolkit/aws"
	"github.com/bartlomiej-jedrol/pr16-aws-serverless-ai-assistant/api"
	"github.com/bartlomiej-jedrol/pr16-aws-serverless-ai-assistant/configuration"
	"github.com/bartlomiej-jedrol/pr16-aws-serverless-ai-assistant/telegram"
)

var (
	awsCfg        *aws.Config
	s3Client      *s3.Client
	s3Bucket      = "pr16-assistant-bucket"
	lambdaTmpPath = "/tmp"
)

func init() {
	serviceName, err := iAws.GetEnvironmentVariable("SERVICE_NAME")
	if err != nil {
		return
	}
	configuration.SetServiceName(serviceName)

	awsCfg, err = iAws.LoadDefaultConfig()
	if err != nil {
		return
	}
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// function := "handler"
	if !api.Authenticate(request.Headers["Authorization"]) {
		return api.BuildResponse(http.StatusForbidden, "unauthorized")
	}

	message, err := telegram.ParseMessage(request.Body)
	if err != nil {
		return api.BuildResponse(http.StatusBadRequest, "bad request")
	}

	return api.BuildResponse(http.StatusOK, message.Text)
}

func main() {
	lambda.Start(handler)
}
