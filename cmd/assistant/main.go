package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bartlomiej-jedrol/pr16-aws-serverless-ai-assistant/telegram"
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("INFO: handler - HTTPMethod, %+v", request.HTTPMethod)
	log.Printf("INFO: handler - Headers, %+v", request.Headers)
	log.Printf("INFO: handler - Body, %+v", request.Body)

	res := telegram.Response{
		Text: "response",
	}
	body, err := json.Marshal(res)
	if err != nil {
		log.Printf("INFO: handler - failed to marshal body, %v", err)
	}

	r := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(body),
	}
	return r, nil
}

func main() {
	lambda.Start(handler)
}
