package clients

import (
	"context"
	"errors"
	"log"
	"strings"

	"google.golang.org/genai"
)

type GeminiClient struct {
	Client *genai.Client
	Model  string
}

func NewGeminiClient(apiKey string) (*GeminiClient, error) {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		return nil, err
	}

	return &GeminiClient{
		Client: client,
		Model:  "gemini-2.5-flash",
	}, nil
}

func (gc *GeminiClient) ClassifyIntent(userInput string) (string, error) {
	log.Println("Gemini Client: Classifying intent...")
	ctx := context.Background()

	prompt := `
		You are an intent classifier for a travel planner.
		Your job is to read the user's request and classify it into one of three categories:
		"plan_day", "refine_plan", "compare_options"

		User Request: "Plan 10:00-18:00 in Kyoto on 2025-12-12. Prefer temples and walkable."
		Classification: "plan_day"

		User Request: "Refine: add a specialty coffee stop near the second venue."
		Classification: "refine_plan"

		User Request: "Compare two options if it rains after 3pm."
		Classification: "compare_options"

		User Request: "` + userInput + `"
		Classification: `

	resp, err := gc.Client.Models.GenerateContent(ctx,
		gc.Model,
		genai.Text(prompt),
		nil,
	)
	if err != nil {
		return "", err
	}

	intent := resp.Text()

	if intent == "" {
		return "", errors.New("gemini returned an empty response")
	}

	return strings.TrimSpace(intent), nil
}
