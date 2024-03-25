package openai

import (
	"encoding/json"
	"fmt"

	resty "github.com/go-resty/resty/v2"
)

const MAX_TOKENS = 250

type Client struct {
	HttpClient *resty.Client
}

func NewClient(token string) *Client {
	httpclient := resty.New()
	httpclient.SetAuthToken(token)

	return &Client{HttpClient: httpclient}
}

func (c Client) Chat(messages []*ChatMessage) (*ChatResponse, error) {
	body := &ChatPayload{
		Messages:  messages,
		Model:     "gpt-4",
		MaxTokens: MAX_TOKENS,
	}

	b, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal the payload. %s", err.Error())
	}

	res, err := c.HttpClient.NewRequest().
		SetHeader("Content-Type", "application/json").
		SetBody(b).
		Post("https://api.openai.com/v1/chat/completions")
	if err != nil {
		return nil, fmt.Errorf("failed to query openai. Error: %s", err.Error())
	}
	if res.IsError() {
		return nil, fmt.Errorf("openai returned [%d]: %s", res.StatusCode(), res.Error())
	}

	var r ChatResponse
	if err := json.Unmarshal(res.Body(), &r); err != nil {
		return nil, fmt.Errorf("failed to parse the openai response. Error: %s", err.Error())
	}

	return &r, nil
}

type ChatResponse struct {
	ID      string               `json:"id"`
	Choices []ChatResponseChoice `json:"choices"`
}

type ChatResponseChoice struct {
	Message ChatResponseChoiceMessage `json:"message"`
}

type ChatResponseChoiceMessage struct {
	Content string `json:"content"`
}

type ChatPayload struct {
	Messages  []*ChatMessage `json:"messages"`
	Model     string         `json:"model"`
	MaxTokens int            `json:"max_tokens"`
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
