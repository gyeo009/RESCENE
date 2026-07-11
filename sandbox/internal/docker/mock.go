package docker

import (
	"context"
	"fmt"
)

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

// DestroySandbox는 모의(Mock) 샌드박스 파괴 요청을 처리합니다.
func (c *mockClient) DestroySandbox(ctx context.Context, containerID string, user string) error {
	return nil
}

// ExecuteCommand는 모의(Mock) 컨테이너 명령어 실행을 처리합니다.
func (c *mockClient) ExecuteCommand(ctx context.Context, containerID string, cmd []string) (string, error) {
	return fmt.Sprintf("mock execution output for cmd: %v", cmd), nil
}

// CopyFileToSandbox는 모의(Mock) 파일 전송을 처리합니다.
func (c *mockClient) CopyFileToSandbox(ctx context.Context, containerID string, content []byte, destPath string) error {
	return nil
}

// GetSandboxLogs는 모의(Mock) 로그 수집을 처리합니다.
func (c *mockClient) GetSandboxLogs(ctx context.Context, containerID string) (string, error) {
	return fmt.Sprintf("mock logs for container %s", containerID), nil
}

// WaitSandbox는 모의(Mock) 종료 완료 대기를 처리합니다.
func (c *mockClient) WaitSandbox(ctx context.Context, containerID string) (int64, error) {
	return 0, nil
}
