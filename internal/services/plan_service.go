package services

import (
	"fmt"
	"log"
	"orchestrator/internal/clients"
	"orchestrator/internal/dtos"
	"strings"

	"github.com/gin-gonic/gin"
)

type PlanService struct {
	Gemini *clients.GeminiClient
	Mcp    *clients.McpClient
}

func NewPlanService(gemini *clients.GeminiClient, mcp *clients.McpClient) *PlanService {
	return &PlanService{
		Gemini: gemini,
		Mcp:    mcp,
	}
}

func (s *PlanService) GeneratePlan(c *gin.Context, req dtos.PlanRequest) {
	log.Println("Plan Service: Generating plan...")

	fullPrompt := fmt.Sprintf("Plan a day in %s on %s. Preferences: %s",
		req.City,
		req.Date,
		strings.Join(req.Preferences, ", "),
	)

	intent, err := s.Gemini.ClassifyIntent(fullPrompt)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to classify intent", "details": err.Error()})
		return
	}
	log.Printf("Intent classified as: %s", intent)

	geo, err := s.Mcp.GetGeocode(req.City)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to geocode city", "details": err.Error()})
		return
	}
	log.Printf("Geocoded city: %s", geo.DisplayName)

	forecast, err := s.Mcp.GetForecast(geo.Lat, geo.Lon, req.Date)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get forecast", "details": err.Error()})
		return
	}
	log.Printf("Forecast: %s", forecast.Summary)

	c.JSON(200, gin.H{
		"intent":    intent,
		"city_info": geo,
		"weather":   forecast,
		"plan":      "This is a placeholder plan. We will implement the full Gemini-powered plan next.",
	})
}
