package telegram

import (
	"encoding/json"

	iLog "github.com/bartlomiej-jedrol/go-toolkit/log"
	"github.com/bartlomiej-jedrol/pr16-aws-serverless-ai-assistant/configuration"
)

type Message struct {
	Text string `json:"text"`
}

type Response struct {
	Text string `json:"text"`
}

func ParseMessage(body string) (*Message, error) {
	function := "ParseRequest"

	msg := Message{}
	err := json.Unmarshal([]byte(body), &msg)
	if err != nil {
		iLog.Error("failed to unmarshal telegram message", nil, err, configuration.ServiceName, function)
		return nil, err
	}
	iLog.Info("telegram message", msg, nil, configuration.ServiceName, function)
	return &msg, nil
}
