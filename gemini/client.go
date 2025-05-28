package gemini

import (
	"context"

	"google.golang.org/genai"
)

type Client struct {
	ctx   context.Context
	c     *genai.Client
	model string
}

func NewClient(token string) *Client {
	ctx := context.Background()

	gClient, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  token,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		panic(err)
	}

	return &Client{
		ctx:   ctx,
		c:     gClient,
		model: "gemini-2.5-flash-preview-05-20",
	}
}

func (client Client) GenerateContent(content []*genai.Content) (*genai.GenerateContentResponse, error) {
	temp := float32(0.75)

	return client.c.Models.GenerateContent(
		client.ctx,
		client.model,
		content,
		&genai.GenerateContentConfig{
			Temperature: &temp,
			SystemInstruction: &genai.Content{
				Role: "system",
				Parts: []*genai.Part{
					genai.NewPartFromText(`Tu es un assistant opérant sur un canal discord.
					Ne fait que des réponses de maximum 1500 caractères.
					Ne réponds qu'au dernier sujet de la discussion en cours, ou qu'à la dernière question.`),
				},
			},
		},
	)
}
