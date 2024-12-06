// Package openai provides tools to work with openai models.
package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
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
		log.Printf("ERROR: buildBody - failed to marshal JSON, %v", err)
		return nil, err
	}
	return jsonData, nil
}

func sendRequest(enpointURL string, APIKey string, body []byte) ([]byte, error) {
	bodyReader := bytes.NewReader(body)

	URL, err := url.Parse(enpointURL)
	if err != nil {
		log.Printf("ERROR: sendRequest - failed to parse URL, %v", err)
		return nil, err
	}

	header := http.Header{}
	header.Add("Content-Type", "application/json")
	header.Add("Authorization", fmt.Sprintf("Bearer %s", APIKey))
	// log.Printf("header: %v", header.Get("Authorization"))
	// log.Printf("header: %v", header.Get("Content-Type"))

	req := http.Request{
		Method: http.MethodPost,
		URL:    URL,
		Header: header,
		Body:   io.NopCloser(bodyReader),
	}

	log.Printf("Request Method: %s", req.Method)
	log.Printf("Request URL: %s", req.URL.String())
	log.Printf("Request Headers: %v", req.Header)
	log.Printf("Request Body: %s", string(body))

	client := http.Client{}
	res, err := client.Do(&req)
	if err != nil {
		log.Printf("ERROR: sendRequest - failed to send request, %v", err)
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Printf("ERROR: sendRequest - received non-OK HTTP status: %s", res.Status)
		return nil, fmt.Errorf("ERROR: sendRequest - received non-OK HTTP status: %s", res.Status)
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("ERROR: sendRequest - failed to read response body, %v", err)
		return nil, err
	}
	return resBody, nil
}

// CreateChatCompletions returns one or more predicted completions.
// It takes prompt and model, commmunicates with OpenAI models, and returns one or more predicted completions.
func CreateChatCompletion(openaiAPIKey string, model string, userContent string, systemContent string) (string, error) {
	reqBody, err := buildBody(model, systemContent, userContent)
	if err != nil {
		return "", err
	}

	resBody, err := sendRequest(chatCompletionsURL, openaiAPIKey, reqBody)
	if err != nil {
		return "", err
	}

	res := Response{}
	err = json.Unmarshal(resBody, &res)
	if err != nil {
		log.Printf("ERROR: CreateChatCompletions - failed to unmarshal JSON, %v", err)
		return "", err
	}

	content := res.Choices[0].Message.Content
	if len(res.Choices) == 0 {
		log.Printf("ERROR: CreateChatCompletions - failed to return completions from OpenAI")
		return "", fmt.Errorf("ERROR: CreateChatCompletions - failed to return completions from OpenAI")
	}

	log.Printf("INFO: CreateChatCompletion - openai content: %v", content)
	return content, nil
}
