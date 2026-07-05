package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"sandbox/internal/dto"
	"sandbox/internal/response"
	"sandbox/internal/service"
)

// SandboxHandler는 HTTP 요청을 수신하여 SandboxService로 비즈니스 로직을 위임하는 핸들러입니다.
type SandboxHandler struct {
	sandboxService service.SandboxService
}

// NewSandboxHandler는 SandboxHandler의 새 인스턴스를 생성합니다.
func NewSandboxHandler(sandboxService service.SandboxService) *SandboxHandler {
	return &SandboxHandler{
		sandboxService: sandboxService,
	}
}

// RunSandbox는 샌드박스 실행 API 요청을 처리합니다.
// @Summary 샌드박스 컨테이너 실행
// @Description 요청된 시나리오와 사용자 정보에 매칭되는 샌드박스 컨테이너를 구동합니다.
// @Tags sandbox
// @Accept json
// @Produce json
// @Param request body dto.SandboxRunRequest true "샌드박스 실행 파라미터"
// @Success 200 {object} response.SandboxRunResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /sandbox/run [post]
func (h *SandboxHandler) RunSandbox(c *gin.Context) {
	var req dto.SandboxRunRequest
	// JSON 본문을 DTO 바인딩 및 유효성 검사 수행
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	// 서비스 계층 호출
	resp, err := h.sandboxService.RunSandbox(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}
