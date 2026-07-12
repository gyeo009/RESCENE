package dto

// SandboxCreateRequest는 샌드박스 생성 요청 바디의 DTO입니다.
type SandboxCreateRequest struct {
	Scenario string `json:"scenario" binding:"required" example:"python:3.10-alpine"` // 실행할 이미지/시나리오명
}

// SandboxRunRequest는 샌드박스 내부 명령 실행 요청 바디의 DTO입니다.
type SandboxRunRequest struct {
	Cmd []string `json:"cmd" binding:"required" example:"python,main.py"` // 실행할 명령어 및 인자 리스트
}

// FileWriteRequest는 샌드박스 작업 공간 파일 생성/수정 요청 바디의 DTO입니다.
type FileWriteRequest struct {
	Filename string `json:"filename" binding:"required" example:"main.py"`      // 생성/수정할 파일명 및 상대 경로
	Content  string `json:"content" binding:"required" example:"print('hello')"` // 파일에 기록할 텍스트 내용
}
