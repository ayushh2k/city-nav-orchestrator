package routes

import (
	"orchestrator/internal/clients"
	"orchestrator/internal/config"
	"orchestrator/internal/handlers"
	"orchestrator/internal/services"

	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg *config.Config) *gin.Engine {
	router := gin.Default()

	geminiClient, err := clients.NewGeminiClient(cfg.GeminiAPIKey)
	if err != nil {
		panic(err)
	}

	mcpClient := clients.NewMcpClient(cfg.McpServerBaseURL, cfg.McpServerAPIKey)

	planService := services.NewPlanService(geminiClient, mcpClient)

	planHandler := handlers.NewPlanHandler(planService)

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "Orchestrator is running"})
	})

	apiV1 := router.Group("/api/v1")
	{
		apiV1.POST("/plan", planHandler.HandlePlanRequest)
	}

	return router
}
