package dto

// SandboxRunRequest represents the request body for running a sandbox
type SandboxRunRequest struct {
	Scenario string `json:"scenario" binding:"required" example:"python-basic"`
	User     string `json:"user" binding:"required" example:"test-user"`
}
