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

// DestroySandbox는 샌드박스 종료 및 삭제 API 요청을 처리합니다.
// @Summary 샌드박스 컨테이너 종료 및 삭제
// @Description 실행 중인 샌드박스 컨테이너를 중지하고 리소스를 제거합니다.
// @Tags sandbox
// @Accept json
// @Produce json
// @Param request body dto.SandboxDestroyRequest true "샌드박스 종료 파라미터"
// @Success 200 {object} map[string]string "성공 메시지"
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /sandbox/stop [post]
func (h *SandboxHandler) DestroySandbox(c *gin.Context) {
	var req dto.SandboxDestroyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	if err := h.sandboxService.DestroySandbox(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "sandbox destroyed successfully"})
}

// ExecuteCommand는 샌드박스 내에서 명령어를 실행하는 API 요청을 처리합니다.
// @Summary 샌드박스 내 명령어 실행
// @Description 실행 중인 샌드박스 내부에서 지정된 명령을 동기식으로 실행하고 출력을 반환합니다.
// @Tags sandbox
// @Accept json
// @Produce json
// @Param request body dto.SandboxExecRequest true "명령어 실행 파라미터"
// @Success 200 {object} response.SandboxExecResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /sandbox/execute [post]
func (h *SandboxHandler) ExecuteCommand(c *gin.Context) {
	var req dto.SandboxExecRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	resp, err := h.sandboxService.ExecuteCommand(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// CopyFile은 샌드박스 내 특정 경로로 파일을 복사하는 API 요청을 처리합니다.
// @Summary 샌드박스 내부로 파일 복사
// @Description 샌드박스 내부의 지정된 경로로 텍스트/스크립트 파일을 복사합니다.
// @Tags sandbox
// @Accept json
// @Produce json
// @Param request body dto.SandboxCopyRequest true "파일 복사 파라미터"
// @Success 200 {object} map[string]string "성공 메시지"
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /sandbox/copy [post]
func (h *SandboxHandler) CopyFile(c *gin.Context) {
	var req dto.SandboxCopyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	if err := h.sandboxService.CopyFile(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "file copied successfully"})
}

// GetLogs는 샌드박스의 표준 출력/에러 로그를 수집하여 반환하는 API 요청을 처리합니다.
// @Summary 샌드박스 컨테이너 로그 조회
// @Description 지정된 샌드박스의 stdout/stderr 로그를 조회합니다.
// @Tags sandbox
// @Produce json
// @Param container_id query string true "컨테이너 식별자 ID"
// @Success 200 {object} response.SandboxLogsResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /sandbox/logs [get]
func (h *SandboxHandler) GetLogs(c *gin.Context) {
	containerID := c.Query("container_id")
	if containerID == "" {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error: "container_id query parameter is required",
		})
		return
	}

	resp, err := h.sandboxService.GetLogs(c.Request.Context(), containerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}
