package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"sandbox/internal/dto"
	"sandbox/internal/response"
	"sandbox/internal/service"
)

// SandboxHandler receives HTTP requests and directs them to the SandboxService
type SandboxHandler struct {
	sandboxService service.SandboxService
}

// NewSandboxHandler creates a new SandboxHandler instance
func NewSandboxHandler(sandboxService service.SandboxService) *SandboxHandler {
	return &SandboxHandler{
		sandboxService: sandboxService,
	}
}

// RunSandbox handles the run sandbox API call
// @Summary Run sandbox container
// @Description Start a new mock sandbox container matching the scenario and user
// @Tags sandbox
// @Accept json
// @Produce json
// @Param request body dto.SandboxRunRequest true "Parameters to run sandbox"
// @Success 200 {object} response.SandboxRunResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /sandbox/run [post]
func (h *SandboxHandler) RunSandbox(c *gin.Context) {
	var req dto.SandboxRunRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	resp, err := h.sandboxService.RunSandbox(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}
