package docker

import "context"

// Client 인터페이스는 Docker 엔진과 상호작용하기 위한 공통 계약을 정의합니다.
// 샌드박스 라이프사이클 관리를 담당하는 기능들이 선언됩니다.
type Client interface {
	// CreateSandbox는 지정된 요청 파라미터를 기반으로 새로운 샌드박스 컨테이너를 생성하고 시작합니다.
	CreateSandbox(ctx context.Context, req CreateSandboxRequest) (*SandboxInfo, error)
}
