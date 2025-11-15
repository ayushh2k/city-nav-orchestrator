package clients

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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

func (gc *GeminiClient) StreamPlan(prompt string) (<-chan string, error) {
	ctx := context.Background()

	stream := make(chan string)

	go func() {
		defer close(stream)

		for resp, err := range gc.Client.Models.GenerateContentStream(ctx,
			gc.Model,
			genai.Text(prompt),
			nil,
		) {
			if err != nil {
				log.Printf("Gemini Stream Error: %v", err)
				return
			}

			text := resp.Text()
			if err != nil {
				log.Printf("Warning: Could not extract text from stream part: %v", err)
				continue
			}

			if text != "" {
				stream <- text
			}
		}
	}()

	return stream, nil
}

func (gc *GeminiClient) GenerateStructuredItinerary(prompt string) (*DraftItinerary, error) {
	ctx := context.Background()

	responseSchema := &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"stops": {
				Type: genai.TypeArray,
				Items: &genai.Schema{
					Type: genai.TypeObject,
					Properties: map[string]*genai.Schema{
						"name":       {Type: genai.TypeString, Description: "The name of the venue."},
						"lat":        {Type: genai.TypeNumber, Description: "Latitude of the venue."},
						"lon":        {Type: genai.TypeNumber, Description: "Longitude of the venue."},
						"start_time": {Type: genai.TypeString, Description: "The suggested start time for the visit (HH:MM)."},
					},
					Required: []string{"name", "lat", "lon", "start_time"},
				},
			},
		},
	}

	config := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema:   responseSchema,
	}

	resp, err := gc.Client.Models.GenerateContent(ctx,
		gc.Model,
		genai.Text(prompt),
		config,
	)
	if err != nil {
		return nil, err
	}

	rawJSON := resp.Text()
	if err != nil {
		return nil, err
	}

	var draft DraftItinerary
	if err := json.Unmarshal([]byte(rawJSON), &draft); err != nil {
		log.Printf("Failed to unmarshal draft itinerary: %s", rawJSON)
		return nil, fmt.Errorf("gemini returned unparsable JSON: %v", err)
	}

	return &draft, nil
}
