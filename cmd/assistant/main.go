package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
)

func handleRequest(ctx context.Context, event json.RawMessage) error {
	// Logging
	log.Printf("INFO: handleRequest - event: %+v", string(event))
	return nil
}

func main() {
	lambda.Start(handleRequest)
}
