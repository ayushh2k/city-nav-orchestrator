package services

import (
	"fmt"
	"log"
	"net/http"
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
	log.Println("Plan Service: Starting full multi-pass generation...")

	geo, err := s.Mcp.GetGeocode(req.City)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to geocode city (Critical)", "details": err.Error()})
		return
	}

	weatherData := getWeatherData(s.Mcp, geo.Lat, geo.Lon, req.Date)
	airData := getAirData(s.Mcp, geo.Lat, geo.Lon, req.Date)
	nearbyData := getNearbyData(s.Mcp, geo.Lat, geo.Lon, req.Preferences)

	log.Println("--- PASS 1: Generating structured itinerary draft...")

	draftPrompt := fmt.Sprintf(`
		You are an expert itinerary generator. Based on the data below, select 4 locations from the VENUE LIST and propose a start time for each.
		Your sole output MUST be a JSON object matching the requested schema. Do NOT include any commentary.
		
		DATA: City: %s, Date: %s, Preferences: %s
		Weather Summary: %s
		Air Quality Summary: %s
		VENUE LIST: %s
		`,
		req.City, req.Date, strings.Join(req.Preferences, ", "),
		weatherData, airData, nearbyData,
	)

	draftItinerary, err := s.Gemini.GenerateStructuredItinerary(draftPrompt)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed Pass 1 (Draft Generation)", "details": err.Error()})
		return
	}
	log.Printf("Draft Itinerary generated successfully with %d stops.", len(draftItinerary.Stops))

	log.Println("--- PASS 2: Calling OSRM tool for travel times...")

	var routePoints []clients.Point
	for _, stop := range draftItinerary.Stops {
		routePoints = append(routePoints, clients.Point{Lat: stop.Lat, Lon: stop.Lon})
	}

	etaReq := clients.EtaRequest{
		Profile: "car",
		Points:  routePoints,
	}

	etaData, err := s.Mcp.GetEta(etaReq)
	etaSummary := ""
	if err != nil {
		log.Printf("Warning: Failed OSRM call: %v", err)
		etaSummary = "Travel times are unavailable due to an API error."
	} else {
		etaSummary = fmt.Sprintf("Travel distance: %.1f km, Duration: %.1f min.", etaData.DistanceKm, etaData.DurationMin)
	}

	log.Println("--- PASS 3: Synthesizing final plan and beginning stream...")

	finalPlanPrompt := fmt.Sprintf(`
		You are the final narrative copilot. Your task is to turn the structured draft and travel data into a professional, easy-to-read, minute-by-minute itinerary.
		
		**INSTRUCTIONS**
		1. **Formatting**: Output ONLY the final Markdown itinerary.
		2. **Trace**: Do NOT include the TRACE or INSTRUCTION blocks in the final output.
		3. **Check**: Incorporate the WEATHER, AIR QUALITY, and TRAVEL TIMES into the narrative.
		
		**CONTEXT**
		- Weather: %s
		- Air Quality: %s
		- Travel Summary (Use this for travel estimates): %s
		
		**DRAFT ITINERARY (Structured Data)**
		%v
	`, weatherData, airData, etaSummary, draftItinerary)
	intent, err := s.Gemini.ClassifyIntent(finalPlanPrompt)

	stream, err := s.Gemini.StreamPlan(finalPlanPrompt)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to initiate Gemini stream", "details": err.Error()})
		return
	}

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Status(http.StatusOK)

	c.Writer.WriteString(fmt.Sprintf("data: [TRACE] Intent: %s\n", intent))
	c.Writer.WriteString(fmt.Sprintf("data: [TRACE] Geocoded: %s\n", geo.DisplayName))
	c.Writer.WriteString(fmt.Sprintf("data: [TRACE] Weather Status: %s\n", weatherData))
	c.Writer.WriteString(fmt.Sprintf("data: [TRACE] OSRM Status: %s\n\n", etaSummary))
	c.Writer.Flush()

	for token := range stream {
		c.Writer.WriteString(fmt.Sprintf("data: %s\n", token))
		c.Writer.Flush()
	}

	c.Writer.WriteString("data: [END]\n\n")
	c.Writer.Flush()
}

func getWeatherData(mcp *clients.McpClient, lat, lon float64, date string) string {
	forecast, err := mcp.GetForecast(lat, lon, date)
	if err != nil {
		log.Printf("Warning: Failed to get weather: %v", err)
		return "Weather data unavailable."
	}
	return fmt.Sprintf("Max Temp: %.1f°C (Precip Prob: %.0f%%)", forecast.TempC, forecast.PrecipProb)
}
func getAirData(mcp *clients.McpClient, lat, lon float64, date string) string {
	aqi, err := mcp.GetAQI(lat, lon, date)
	if err != nil {
		log.Printf("Warning: Failed to get AQI: %v", err)
		return "Air quality data unavailable."
	}
	return fmt.Sprintf("PM2.5: %.1f µg/m³, PM10: %.1f µg/m³, Category: %s",
		aqi.PM25, aqi.PM10, aqi.Category)
}

func getNearbyData(mcp *clients.McpClient, lat, lon float64, preferences []string) string {
	var allVenues []string

	venueMap := make(map[string]bool)

	for _, pref := range preferences {
		query := strings.TrimSpace(pref)
		if query == "walkable" {
			continue
		}

		log.Printf("MCP Client: Getting nearby venues for preference: %s", query)
		nearby, err := mcp.GetNearby(lat, lon, query)
		if err != nil {
			log.Printf("Warning: Failed to get nearby venues for %s: %v", query, err)
			continue
		}

		for _, place := range nearby.Places {
			if !venueMap[place.Name] {
				venueMap[place.Name] = true
				allVenues = append(allVenues, fmt.Sprintf("- %s (%.6f, %.6f)", place.Name, place.Lat, place.Lon))
			}
		}
	}

	if len(allVenues) == 0 {
		return "No nearby venues found based on preferences."
	}

	return strings.Join(allVenues, "\n")
}
