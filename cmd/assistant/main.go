package main

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	iAws "github.com/bartlomiej-jedrol/go-toolkit/aws"
	"github.com/bartlomiej-jedrol/pr16-aws-serverless-ai-assistant/api"
	"github.com/bartlomiej-jedrol/pr16-aws-serverless-ai-assistant/bloodresults"
	"github.com/bartlomiej-jedrol/pr16-aws-serverless-ai-assistant/configuration"
)

var (
	s3Client      *s3.Client
	bloodParser   bloodresults.Parser
	s3Bucket      = "pr16-assistant-bucket"
	lambdaTmpPath = "/tmp"
)

func init() {
	serviceName, err := iAws.GetEnvironmentVariable("SERVICE_NAME")
	if err != nil {
		return
	}
	configuration.SetServiceName(serviceName)

	cfg, err := iAws.LoadDefaultConfig()
	if err != nil {
		return
	}

	s3Client = s3.NewFromConfig(*cfg)
	bloodParser = *bloodParser.New(s3Client, s3Bucket, lambdaTmpPath)
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// function := "handler"
	if !api.Authenticate(request.Headers["Authorization"]) {
		return api.BuildResponse(http.StatusForbidden, "unauthorized")
	}

	assistantRequest, err := api.Parse(request.Body)
	if err != nil {
		return api.BuildResponse(http.StatusInternalServerError, "bad request")
	}

	if assistantRequest.Text == "parse blood results" && assistantRequest.S3ObjectKey != "" {
		err := bloodParser.Parse(ctx, assistantRequest.S3ObjectKey)
		if err != nil {
			api.BuildResponse(http.StatusInternalServerError, err.Error())
		}
	}

	return api.BuildResponse(http.StatusOK, assistantRequest.Text)
}

func main() {
	lambda.Start(handler)
}
