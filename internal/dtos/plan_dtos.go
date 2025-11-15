package dtos

type PlanRequest struct {
	City        string   `json:"city" binding:"required"`
	Date        string   `json:"date" binding:"required"`
	Preferences []string `json:"preferences"`
}

type PlanResponse struct {
	Plan      string   `json:"plan"`
	ToolTrace []string `json:"tool_trace"`
	Error     string   `json:"error,omitempty"`
}
