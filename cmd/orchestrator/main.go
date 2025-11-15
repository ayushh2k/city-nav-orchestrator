package main

import (
	"fmt"
	"log"
	"orchestrator/internal/config"
	"orchestrator/internal/routes"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Note: .env file not found, loading from system env")
	}

	cfg := config.LoadConfig()

	router := routes.SetupRouter(cfg)

	fmt.Println("Starting Orchestrator server on http://localhost:8080")
	router.Run(":8080")
}
