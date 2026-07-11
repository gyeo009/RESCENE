package docker

// SandboxInfo는 생성 및 실행된 샌드박스 컨테이너의 핵심 메타데이터를 담고 있습니다.
type SandboxInfo struct {
	ContainerID string // 생성된 Docker 컨테이너의 고유 식별자 ID
	Status      string // 현재 컨테이너 상태 (예: "running")
}

// CreateSandboxRequest는 새로운 샌드박스를 시작할 때 필요한 사용자 및 시나리오 메타데이터입니다.
type CreateSandboxRequest struct {
	User     string // 샌드박스를 요청하는 사용자 식별자
	Scenario string // 실행할 시나리오 및 컨테이너 이미지명
	Runtime  string // 요청하는 컨테이너 런타임 (선택)
}
