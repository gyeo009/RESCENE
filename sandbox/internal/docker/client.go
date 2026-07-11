package docker

import "context"

// Client는 Docker 컨테이너 실행환경 관리를 위한 인터페이스입니다.
// 이 인터페이스는 Docker SDK 개념과 유사하게 유지하며, 파일 조작 등 호스트 FS 로직은 포함하지 않습니다.
type Client interface {
	// CreateContainer는 새 컨테이너를 생성하고 생성된 컨테이너 ID를 반환합니다.
	// name: 컨테이너에 부여할 고유 이름 (공백일 경우 자동 생성)
	// image: 사용할 Docker 이미지명
	// cmd: 기본 컨테이너 구동 명령어 (예: ["tail", "-f", "/dev/null"])
	// binds: 호스트 볼륨 바인드 마운트 목록 (예: ["/host/path:/workspace"])
	// env: 컨테이너 환경변수 목록 (예: ["ENV=production"])
	CreateContainer(ctx context.Context, name string, image string, cmd []string, binds []string, env []string) (string, error)

	// StartContainer는 생성된 컨테이너를 구동합니다.
	StartContainer(ctx context.Context, containerID string) error

	// Exec는 실행 중인 컨테이너 내부에서 명령어를 동기로 실행하고 표준 출력/표준 에러 결과 문자열을 반환합니다.
	Exec(ctx context.Context, containerID string, cmd []string) (string, error)

	// Logs는 컨테이너의 stdout/stderr 실행 로그 전체를 가져옵니다.
	Logs(ctx context.Context, containerID string) (string, error)

	// StopContainer는 컨테이너 작동을 중지시킵니다.
	StopContainer(ctx context.Context, containerID string) error

	// RemoveContainer는 컨테이너를 완전히 삭제합니다.
	RemoveContainer(ctx context.Context, containerID string) error
}
