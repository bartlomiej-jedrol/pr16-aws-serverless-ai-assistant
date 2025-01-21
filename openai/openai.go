// Package openai provides tools to work with openai models.
package openai

import (
	"encoding/json"

	iHTTP "github.com/bartlomiej-jedrol/go-toolkit/http"
	iLog "github.com/bartlomiej-jedrol/go-toolkit/log"
	"github.com/bartlomiej-jedrol/pr16-aws-serverless-ai-assistant/configuration"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Payload struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Choice struct {
	Message `json:"message"`
}

type Response struct {
	Choices []Choice `json:"choices"`
}

var (
	chatCompletionsURL = "https://api.openai.com/v1/chat/completions"
)

func buildMessages(systemContent string, userContent string) []Message {
	messages := []Message{}
	messages = append(messages, Message{
		Role:    "system",
		Content: systemContent})
	messages = append(messages, Message{
		Role:    "user",
		Content: userContent})
	return messages
}

func buildBody(model string, systemContent string, userContent string) ([]byte, error) {
	messages := buildMessages(systemContent, userContent)
	payload := Payload{
		Model:    model,
		Messages: messages,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		iLog.Error("failed to marshal JSON", nil, err, configuration.ServiceName, "buildBody")
		return nil, err
	}
	return jsonData, nil
}

// CreateChatCompletions returns one or more predicted completions.
// It takes prompt and model, commmunicates with OpenAI models, and returns one or more predicted completions.
func CreateChatCompletion(openaiAPIKey string, model string, userContent string, systemContent string) (string, error) {
	function := "CreateChatCompletion"

	reqBody, err := buildBody(model, systemContent, userContent)
	if err != nil {
		return "", err
	}

	resBody, err := iHTTP.SendHTTPRequest(chatCompletionsURL, openaiAPIKey, reqBody)
	if err != nil {
		return "", err
	}

	res := Response{}
	err = json.Unmarshal(resBody, &res)
	if err != nil {
		iLog.Error("failed to unmarshal JSON", nil, err, configuration.ServiceName, function)
		return "", err
	}

	if len(res.Choices) == 0 {
		iLog.Error("failed to return completions from OpenAI", nil, nil, configuration.ServiceName, function)
		return "", err
	}

	content := res.Choices[0].Message.Content
	iLog.Info("openai content", content, nil, configuration.ServiceName, function)
	return content, nil
}
