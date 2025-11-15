package config

import (
	"log"
	"os"
)

type Config struct {
	GeminiAPIKey     string
	McpServerBaseURL string
	McpServerAPIKey  string
}

func LoadConfig() *Config {
	geminiKey := os.Getenv("GEMINI_API_KEY")
	if geminiKey == "" {
		log.Fatal("FATAL: GEMINI_API_KEY environment variable not set.")
	}

	mcpURL := os.Getenv("MCP_SERVER_BASE_URL")
	if mcpURL == "" {
		log.Fatal("FATAL: MCP_SERVER_BASE_URL environment variable not set.")
	}

	mcpKey := os.Getenv("MCP_SERVER_API_KEY")
	if mcpKey == "" {
		log.Fatal("FATAL: MCP_SERVER_API_KEY environment variable not set.")
	}

	return &Config{
		GeminiAPIKey:     geminiKey,
		McpServerBaseURL: mcpURL,
		McpServerAPIKey:  mcpKey,
	}
}
