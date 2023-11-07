package openai

import (
	"encoding/json"
	"fmt"
	"regexp"

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
		Model:     "gpt-4-1106-preview",
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
		return nil, fmt.Errorf("openai return [%d]: %s", res.StatusCode(), res.Error())
	}

	fmt.Println("body response: ", string(res.Body()), res.Status())

	var r ChatResponse
	if err := json.Unmarshal(res.Body(), &r); err != nil {
		return nil, fmt.Errorf("failed to parse the openai response. Error: %s", err.Error())
	}

	return &r, nil
}

func (c Client) Completions(prompt string) (*CompletionResponse, error) {
	body := &CompletionPayload{
		Model:       "text-davinci-003",
		Prompt:      cleanPrompt(prompt),
		MaxTokens:   100,
		Temperature: 1,
		N:           1,
	}

	b, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal the payload. %s", err.Error())
	}

	res, err := c.HttpClient.NewRequest().
		SetHeader("Content-Type", "application/json").
		SetBody(b).
		Post("https://api.openai.com/v1/completions")
	if err != nil {
		return nil, fmt.Errorf("failed to query openai. Error: %s", err.Error())
	}

	var r CompletionResponse
	if err := json.Unmarshal(res.Body(), &r); err != nil {
		return nil, fmt.Errorf("failed to parse the openai response. Error: %s", err.Error())
	}

	return &r, nil
}

func cleanPrompt(p string) string {
	regex := regexp.MustCompile(`<@\d+>`)
	return regex.ReplaceAllString(p, "")
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

// -----

type CompletionPayload struct {
	Model       string  `json:"model"`
	Prompt      string  `json:"prompt"`
	MaxTokens   int     `json:"max_tokens"`
	Temperature float32 `json:"temperature"`
	TopP        float32 `json:"top_p"`
	N           int     `json:"n"`
}

type CompletionPayloadMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
	Name    string `json:"name,omitempty"`
}

type CompletionResponse struct {
	ID      string             `json:"id"`
	Object  string             `json:"object"`
	Created int64              `json:"created"`
	Model   string             `json:"model"`
	Usage   CompletionUsage    `json:"usage"`
	Choices []CompletionChoice `json:"choices"`
}

type CompletionUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type CompletionChoice struct {
	Text         string `json:"text"`
	FinishReason string `json:"finish_reason"`
	Index        int    `json:"index"`
}
