package api

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/bartlomiej-jedrol/pr16-aws-serverless-ai-assistant/telegram"
)

func Authenticate(apiKey string) bool {
	ak := strings.TrimPrefix(apiKey, "Bearer ")
	return ak == os.Getenv("ASSISTANT_API_KEY")
}

func BuildResponse(statusCode int, body any) (events.APIGatewayProxyResponse, error) {
	res := telegram.Response{}
	switch v := body.(type) {
	case error:
		res.Text = v.Error()
	case string:
		res.Text = v
	default:
		res.Text = ""
	}

	r, err := json.Marshal(res)
	if err != nil {
		log.Printf("INFO: handler - failed to marshal body, %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       `{"text": "internal server error}`,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       string(r),
	}, nil
}
