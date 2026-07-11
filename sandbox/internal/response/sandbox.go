package response

// SandboxRunResponse는 샌드박스 컨테이너가 성공적으로 구동되었을 때 반환하는 응답 구조체입니다.
type SandboxRunResponse struct {
	ContainerID string `json:"container_id" example:"sandbox-test-user"` // 생성된 컨테이너 식별자 ID
	Status      string `json:"status" example:"running"`                 // 현재 컨테이너 상태
}

// SandboxExecResponse는 명령어 실행 결과를 반환하는 응답 구조체입니다.
type SandboxExecResponse struct {
	Output string `json:"output" example:"hello\n"` // 실행된 명령어의 표준 출력/에러 결과
}

// SandboxLogsResponse는 로그 조회 결과를 반환하는 응답 구조체입니다.
type SandboxLogsResponse struct {
	Logs string `json:"logs" example:"standard output log..."` // 컨테이너 로그 전체 출력
}

// ErrorResponse는 API 에러 발생 시 반환하는 표준 에러 응답 구조체입니다.
type ErrorResponse struct {
	Error string `json:"error" example:"invalid request parameters"` // 에러 상세 메시지
}
