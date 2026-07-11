package dto

// SandboxRunRequest는 샌드박스 실행 요청 바디의 DTO입니다.
type SandboxRunRequest struct {
	Scenario string `json:"scenario" binding:"required" example:"python-basic"` // 실행할 시나리오 이름 (필수)
	User     string `json:"user" binding:"required" example:"test-user"`         // 요청 사용자 식별자 (필수)
	Runtime  string `json:"runtime" example:"runsc"`                             // 실행할 컨테이너 런타임 (선택)
}

// SandboxDestroyRequest는 샌드박스 종료 및 삭제 요청 바디의 DTO입니다.
type SandboxDestroyRequest struct {
	ContainerID string `json:"container_id" binding:"required" example:"sandbox-test-user"` // 종료할 컨테이너 식별자 ID (필수)
	User        string `json:"user" binding:"required" example:"test-user"`                 // 요청 사용자 식별자 (필수)
}

// SandboxExecRequest는 샌드박스 내부 명령 실행 요청 바디의 DTO입니다.
type SandboxExecRequest struct {
	ContainerID string   `json:"container_id" binding:"required" example:"sandbox-test-user"` // 명령을 실행할 컨테이너 ID (필수)
	Cmd         []string `json:"cmd" binding:"required" example:"[\"echo\", \"hello\"]"`     // 실행할 명령어 배열 (필수)
}

// SandboxCopyRequest는 샌드박스 내부 파일 생성/복사 요청 바디의 DTO입니다.
type SandboxCopyRequest struct {
	ContainerID string `json:"container_id" binding:"required" example:"sandbox-test-user"` // 파일을 복사할 컨테이너 ID (필수)
	Content     string `json:"content" binding:"required" example:"print('hello')"`         // 파일에 저장할 텍스트 내용 (필수)
	DestPath    string `json:"dest_path" binding:"required" example:"/workspace/test.py"`   // 저장할 컨테이너 내부 경로 (필수)
}
