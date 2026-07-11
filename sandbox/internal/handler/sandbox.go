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

// CreateSandbox는 새로운 샌드박스 작업 공간 및 실행 컨테이너를 생성합니다.
// @Summary 샌드박스 환경 생성
// @Description 지정된 이미지를 활용해 격리된 실행 환경 컨테이너를 올리고, 호스트에 마운트할 작업 폴더를 생성합니다.
// @Tags sandbox
// @Accept json
// @Produce json
// @Param request body dto.SandboxCreateRequest true "샌드박스 이미지/시나리오 정보"
// @Success 200 {object} response.SandboxCreateResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /sandbox [post]
func (h *SandboxHandler) CreateSandbox(c *gin.Context) {
	var req dto.SandboxCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	resp, err := h.sandboxService.CreateSandbox(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteSandbox는 지정된 샌드박스 실행 환경 및 디렉토리를 완전히 파괴합니다.
// @Summary 샌드박스 완전 삭제 및 정리
// @Description 샌드박스 컨테이너 작동을 멈추고 삭제한 뒤, 호스트의 작업 디렉토리를 비웁니다.
// @Tags sandbox
// @Produce json
// @Param id path string true "샌드박스 식별 ID"
// @Success 200 {object} response.SuccessResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /sandbox/{id} [delete]
func (h *SandboxHandler) DeleteSandbox(c *gin.Context) {
	id := c.Param("id")
	if err := h.sandboxService.DeleteSandbox(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// ListFiles는 샌드박스 내부 작업 공간에 저장된 모든 소스파일 명단을 상대경로 리스트로 받습니다.
// @Summary 샌드박스 내 파일 목록 조회
// @Description 사용자가 샌드박스 내부에서 작성/수정한 파일 전체 경로들을 나열합니다.
// @Tags files
// @Produce json
// @Param id path string true "샌드박스 식별 ID"
// @Success 200 {object} response.FileListResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /sandbox/{id}/files [get]
func (h *SandboxHandler) ListFiles(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.sandboxService.ListFiles(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ReadFile은 샌드박스 내부의 개별 파일 텍스트 내용을 전달받습니다.
// @Summary 파일 내용 상세 조회
// @Description 특정 샌드박스의 개별 파일 내용을 텍스트로 읽어옵니다.
// @Tags files
// @Produce json
// @Param id path string true "샌드박스 식별 ID"
// @Param filename query string true "조회하려는 파일의 상대경로"
// @Success 200 {object} response.FileContentResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /sandbox/{id}/files/content [get]
func (h *SandboxHandler) ReadFile(c *gin.Context) {
	id := c.Param("id")
	filename := c.Query("filename")
	if filename == "" {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error: "filename query parameter is required",
		})
		return
	}

	resp, err := h.sandboxService.ReadFile(id, filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// WriteFile은 샌드박스 작업 디렉토리에 새로운 파일을 쓰거나 기존 파일을 덮어씁니다.
// @Summary 파일 생성 및 저장
// @Description 샌드박스 작업 디렉토리에 임의의 파일을 생성하거나 파일의 내용을 변경합니다.
// @Tags files
// @Accept json
// @Produce json
// @Param id path string true "샌드박스 식별 ID"
// @Param request body dto.FileWriteRequest true "작성할 파일 경로 및 본문 텍스트"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /sandbox/{id}/files [put]
func (h *SandboxHandler) WriteFile(c *gin.Context) {
	id := c.Param("id")
	var req dto.FileWriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	if err := h.sandboxService.WriteFile(id, &req); err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// DeleteFile은 샌드박스 내부의 특정 파일을 삭제합니다.
// @Summary 파일 삭제
// @Description 샌드박스 작업 폴더 내부의 특정 파일을 삭제합니다.
// @Tags files
// @Produce json
// @Param id path string true "샌드박스 식별 ID"
// @Param filename query string true "삭제하려는 파일의 상대경로"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /sandbox/{id}/files [delete]
func (h *SandboxHandler) DeleteFile(c *gin.Context) {
	id := c.Param("id")
	filename := c.Query("filename")
	if filename == "" {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error: "filename query parameter is required",
		})
		return
	}

	if err := h.sandboxService.DeleteFile(id, filename); err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// RunSandbox는 Docker 클라이언트 exec를 호출하여 컨테이너 내에서 명령을 수행하고 실행 결과를 반환합니다.
// @Summary 샌드박스 내 명령어 실행
// @Description 구동 중인 샌드박스 컨테이너 내부에서 동기로 명령을 Exec하여 그 출력을 가져옵니다.
// @Tags sandbox
// @Accept json
// @Produce json
// @Param id path string true "샌드박스 식별 ID"
// @Param request body dto.SandboxRunRequest true "실행 명령 정보"
// @Success 200 {object} response.SandboxRunResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /sandbox/{id}/run [post]
func (h *SandboxHandler) RunSandbox(c *gin.Context) {
	id := c.Param("id")
	var req dto.SandboxRunRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	resp, err := h.sandboxService.RunCommand(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// StopSandbox는 구동 중인 컨테이너 프로세스를 정지합니다.
// @Summary 샌드박스 컨테이너 중지
// @Description 작동 중인 샌드박스 컨테이너를 중지 상태로 변경합니다.
// @Tags sandbox
// @Produce json
// @Param id path string true "샌드박스 식별 ID"
// @Success 200 {object} response.SuccessResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /sandbox/{id}/stop [post]
func (h *SandboxHandler) StopSandbox(c *gin.Context) {
	id := c.Param("id")
	if err := h.sandboxService.StopSandbox(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// GetLogs는 샌드박스 컨테이너에서 발생한 로그 전체를 조회합니다.
// @Summary 샌드박스 컨테이너 로그 조회
// @Description 샌드박스 컨테이너의 표준 출력 및 에러 출력을 조회합니다.
// @Tags sandbox
// @Produce json
// @Param id path string true "샌드박스 식별 ID"
// @Success 200 {object} response.SandboxRunResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /sandbox/{id}/logs [get]
func (h *SandboxHandler) GetLogs(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.sandboxService.GetLogs(c.Request.Context(), id)
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
