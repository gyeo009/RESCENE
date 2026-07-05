package docker

import "context"

// mockClient는 실제 Docker Engine과의 통신 없이 가상 데이터를 반환하는 Mock 구현체입니다.
// 개발 및 API 테스트 환경에서 활용됩니다.
type mockClient struct{}

// NewMockClient는 mockClient의 새 인스턴스를 반환합니다.
func NewMockClient() Client {
	return &mockClient{}
}

// CreateSandbox는 모의(Mock) 샌드박스 컨테이너 생성 결과를 반환합니다.
func (c *mockClient) CreateSandbox(ctx context.Context, req CreateSandboxRequest) (*SandboxInfo, error) {
	return &SandboxInfo{
		ContainerID: "sandbox-" + req.User,
		Status:      "running",
	}, nil
}
