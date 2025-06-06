package api

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	iLog "github.com/bartlomiej-jedrol/go-toolkit/log"
	"github.com/bartlomiej-jedrol/pr16-aws-serverless-ai-assistant/configuration"
	"github.com/bartlomiej-jedrol/pr16-aws-serverless-ai-assistant/telegram"
)

type Request struct {
	Text        string `json:"text"`
	S3ObjectKey string `json:"s3_object_key"`
}

type Response struct {
	Text string
}

func Authenticate(apiKey string) bool {
	ak := strings.TrimPrefix(apiKey, "Bearer ")
	return ak == os.Getenv("ASSISTANT_API_KEY")
}

func BuildResponse(statusCode int, body any) (events.APIGatewayProxyResponse, error) {
	function := "BuildResponse"
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
		iLog.Info("failed to marshal body", nil, err, configuration.ServiceName, function)
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
