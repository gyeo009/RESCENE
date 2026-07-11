package sandbox

import "context"

// Config는 런타임에 종속되지 않는 샌드박스 생성 설정입니다.
type Config struct {
	ID        string // 샌드박스 식별 ID
	Scenario  string // 사용할 이미지/시나리오명
	Runtime   string // 컨테이너 런타임 (예: runc, runsc, kata 등)
	CPU       int    // CPU 제한 개수 (0이면 기본값/제한 없음)
	Memory    string // 메모리 제한 문자열 (예: "512m", "1g")
	Workspace string // 컨테이너에 마운트될 호스트 상의 작업 디렉토리 절대 경로
	Network   bool   // 네트워크 사용 가능 여부
}

// Sandbox는 가상의 격리 구동 환경을 나타내는 인터페이스입니다.
type Sandbox interface {
	ID() string
	Status() string
	Exec(ctx context.Context, cmd []string) (string, error)
	Copy(ctx context.Context, content []byte, destPath string) error
	Logs(ctx context.Context) (string, error)
	Stop(ctx context.Context) error
	Destroy(ctx context.Context) error
	Wait(ctx context.Context) (int64, error)
}

// Manager는 샌드박스의 수명 주기를 제어하고 검색하는 인터페이스입니다.
type Manager interface {
	Create(ctx context.Context, cfg Config) (Sandbox, error)
	Get(ctx context.Context, id string) (Sandbox, error)
}
