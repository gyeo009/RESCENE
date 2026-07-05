package dto

// SandboxRunRequest는 샌드박스 실행 요청 바디의 DTO입니다.
type SandboxRunRequest struct {
	Scenario string `json:"scenario" binding:"required" example:"python-basic"` // 실행할 시나리오 이름 (필수)
	User     string `json:"user" binding:"required" example:"test-user"`         // 요청 사용자 식별자 (필수)
}
