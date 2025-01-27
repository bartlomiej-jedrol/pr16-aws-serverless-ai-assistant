package main

import (
	"context"
	"io"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	iAws "github.com/bartlomiej-jedrol/go-toolkit/aws"
	iLog "github.com/bartlomiej-jedrol/go-toolkit/log"
	"github.com/bartlomiej-jedrol/pr16-aws-serverless-ai-assistant/api"
	"github.com/bartlomiej-jedrol/pr16-aws-serverless-ai-assistant/configuration"
	convertapi "github.com/bartlomiej-jedrol/pr16-aws-serverless-ai-assistant/convert-api"
)

var (
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

	cfg, err := iAws.LoadDefaultConfig()
	if err != nil {
		return
	}

	s3Client = s3.NewFromConfig(*cfg)
}

func parseBloodResults(ctx context.Context, s3ObjKey string) error {
	function := "parseBloodResults"
	s3OjectBody, err := iAws.GetS3Object(ctx, s3Client, s3Bucket, s3ObjKey)
	if err != nil {
		return err
	}

	tmpFile, err := os.CreateTemp(lambdaTmpPath, "blood-results-*.pdf")
	if err != nil {
		iLog.Error("failed to create pdf file", nil, err, configuration.ServiceName, function)
		return err
	}
	defer tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	_, err = io.Copy(tmpFile, *s3OjectBody)
	if err != nil {
		iLog.Error("failed to save pdf file to tmp path", nil, err, configuration.ServiceName, function)
		return err
	}
	iLog.Info("pdf saved to", tmpFile.Name(), nil, configuration.ServiceName, function)

	outputFiles, err := convertapi.ConvertPDFToJPG(tmpFile.Name(), lambdaTmpPath)
	if err != nil {
		iLog.Error("failed to convert PDF to JPG", err, nil, configuration.ServiceName, function)
		return err
	}

	for _, file := range outputFiles {
		defer os.Remove(file)
		iLog.Info("converted file", file, nil, configuration.ServiceName, function)
	}
	return nil
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
		err := parseBloodResults(ctx, assistantRequest.S3ObjectKey)
		if err != nil {
			api.BuildResponse(http.StatusInternalServerError, err.Error())
		}
	}

	return api.BuildResponse(http.StatusOK, assistantRequest.Text)
}

func main() {
	lambda.Start(handler)
}
