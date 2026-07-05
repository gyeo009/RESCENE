package response

// SandboxRunResponse represents the response body when successfully starting a sandbox
type SandboxRunResponse struct {
	ContainerID string `json:"container_id" example:"sandbox-test-user"`
	Status      string `json:"status" example:"running"`
}

// ErrorResponse represents the standard response body for errors
type ErrorResponse struct {
	Error string `json:"error" example:"invalid request parameters"`
}
