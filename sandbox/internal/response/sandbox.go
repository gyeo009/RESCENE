package response

// SandboxCreateResponse는 샌드박스가 정상적으로 생성되었을 때 반환하는 응답 구조체입니다.
type SandboxCreateResponse struct {
	SandboxID   string `json:"sandbox_id" example:"sandbox-uuid-123"`    // 생성된 고유 샌드박스(작업 공간) 식별자
	ContainerID string `json:"container_id" example:"container-uuid-123"` // 구동된 런타임 컨테이너 식별자 ID
}

// SandboxRunResponse는 샌드박스 명령어 실행 결과를 담고 있는 응답 구조체입니다.
type SandboxRunResponse struct {
	Output string `json:"output" example:"hello world\n"` // 명령어 실행 후 표준출력/에러 결과
}

// FileListResponse는 샌드박스 내부의 파일 목록을 조회한 응답 구조체입니다.
type FileListResponse struct {
	Files []string `json:"files"` // 작업 공간 안의 모든 파일명 리스트
}

// FileContentResponse는 특정 파일의 파일명과 텍스트 내용을 응답하는 구조체입니다.
type FileContentResponse struct {
	Filename string `json:"filename" example:"main.py"` // 파일명
	Content  string `json:"content" example:"print(1)"` // 파일 텍스트 콘텐츠
}

// ErrorResponse는 API 에러 발생 시 반환하는 표준 에러 응답 구조체입니다.
type ErrorResponse struct {
	Error string `json:"error" example:"invalid request parameters"` // 에러 상세 메시지
}

// SuccessResponse는 요청이 성공적으로 완료되었음을 나타내는 기본 응답 구조체입니다.
type SuccessResponse struct {
	Status string `json:"status" example:"success"` // 성공 여부 상태 ("success")
}
