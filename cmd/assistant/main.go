package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	iAws "github.com/bartlomiej-jedrol/go-toolkit/aws"
	iLog "github.com/bartlomiej-jedrol/go-toolkit/log"
	"github.com/bartlomiej-jedrol/pr16-aws-serverless-ai-assistant/api"
	"github.com/bartlomiej-jedrol/pr16-aws-serverless-ai-assistant/configuration"
	"github.com/bartlomiej-jedrol/pr16-aws-serverless-ai-assistant/telegram"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	awsCfg *aws.Config
)

func init() {
	function := "init"
	serviceName, err := iAws.GetEnvironmentVariable("SERVICE_NAME")
	if err != nil {
		return
	}
	configuration.SetServiceName(serviceName)

	dbHost, err := iAws.GetEnvironmentVariable("DB_HOST")
	if err != nil {
		return
	}
	dbUser, err := iAws.GetEnvironmentVariable("DB_USER")
	if err != nil {
		return
	}
	dbPass, err := iAws.GetEnvironmentVariable("DB_PASSWORD")
	if err != nil {
		return
	}
	dbName, err := iAws.GetEnvironmentVariable("DB_NAME")
	if err != nil {
		return
	}

	awsCfg, err = iAws.LoadDefaultConfig()
	if err != nil {
		return
	}

	psqlConnStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=require",
		dbHost, 5432, dbUser, dbPass, dbName)

	pool, err := pgxpool.New(context.Background(), psqlConnStr)
	if err != nil {
		iLog.Error("failed to create connection pool", nil, err, serviceName, function)
		return
	}
	defer pool.Close()

	var greeting string
	err = pool.QueryRow(context.Background(), "SELECT 'Hello, RDS!'").Scan(&greeting)
	if err != nil {
		iLog.Error("failed to run query row", nil, err, serviceName, function)
		return
	}

	iLog.Info(greeting, nil, nil, serviceName, function)
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
