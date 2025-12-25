package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"hal/env"
)

type Client struct {
	http     http.Client
	token    string
	model    string
	endpoint string
}

func NewClient(env *env.Env) *Client {
	return &Client{
		http:     http.Client{Timeout: 10 * time.Second},
		token:    env.LlmToken,
		model:    env.LlmModel,
		endpoint: env.LlmEndpoint,
	}
}

func (client Client) GenerateContent(message string) (string, error) {
	payload := map[string]any{
		"model":  client.model,
		"stream": false,
		"prompt": message,
	}

	payloadBytes, _ := json.Marshal(payload)

	res, err := client.http.Post(
		client.endpoint,
		"application/json",
		bytes.NewBuffer(payloadBytes),
	)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return "", fmt.Errorf("LLM returned status code: [%d]", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read the llm response body: %w", err)
	}

	var llmResponse LLMResponse
	if err := json.Unmarshal(body, &llmResponse); err != nil {
		return "", fmt.Errorf("failed to unmarshal llm response: %w", err)
	}

	return llmResponse.Response, nil
}

type LLMResponse struct {
	Response string `json:"response"`
}
