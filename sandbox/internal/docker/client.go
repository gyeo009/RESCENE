package docker

import "context"

// Client 인터페이스는 Docker 엔진과 상호작용하기 위한 공통 계약을 정의합니다.
// 샌드박스 라이프사이클 관리를 담당하는 기능들이 선언됩니다.
type Client interface {
	// CreateSandbox는 지정된 요청 파라미터를 기반으로 새로운 샌드박스 컨테이너를 생성하고 시작합니다.
	CreateSandbox(ctx context.Context, req CreateSandboxRequest) (*SandboxInfo, error)
	// DestroySandbox는 실행 중인 샌드박스를 종료하고, 호스트 볼륨 등 임시 자원을 정리합니다.
	DestroySandbox(ctx context.Context, containerID string, user string) error
	// ExecuteCommand는 실행 중인 샌드박스 컨테이너 내에서 명령을 수행하고 실행 결과(Stdout/Stderr)를 반환합니다.
	ExecuteCommand(ctx context.Context, containerID string, cmd []string) (string, error)
	// CopyFileToSandbox는 지정된 파일 내용을 타르 스트림으로 구성하여 컨테이너 내부로 전송합니다.
	CopyFileToSandbox(ctx context.Context, containerID string, content []byte, destPath string) error
	// GetSandboxLogs는 컨테이너의 표준 출력 및 표준 에러 스트림 전체를 수집합니다.
	GetSandboxLogs(ctx context.Context, containerID string) (string, error)
	// WaitSandbox는 컨테이너의 종료 상태를 대기하며 상태 코드를 반환합니다.
	WaitSandbox(ctx context.Context, containerID string) (int64, error)
}
