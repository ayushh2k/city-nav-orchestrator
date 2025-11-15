package handlers

import (
	"net/http"
	"orchestrator/internal/dtos"
	"orchestrator/internal/services"

	"github.com/gin-gonic/gin"
)

type PlanHandler struct {
	PlanService *services.PlanService
}

func NewPlanHandler(service *services.PlanService) *PlanHandler {
	return &PlanHandler{
		PlanService: service,
	}
}

func (h *PlanHandler) HandlePlanRequest(c *gin.Context) {
	var req dtos.PlanRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.PlanService.GeneratePlan(c, req)
}
