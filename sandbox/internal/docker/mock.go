package docker

import (
	"context"
	"fmt"
	"strings"
)

// mockClient는 실제 Docker Engine과의 통신 없이 가상 데이터를 반환하는 Mock 구현체입니다.
// 개발 및 API 테스트 환경에서 활용됩니다.
type mockClient struct{}

// NewMockClient는 mockClient의 새 인스턴스를 반환합니다.
func NewMockClient() Client {
	return &mockClient{}
}

// CreateContainer는 모의 컨테이너 ID를 생성하여 반환합니다.
func (c *mockClient) CreateContainer(ctx context.Context, name string, image string, cmd []string, binds []string, env []string) (string, error) {
	return "mock-container-id-12345", nil
}

// StartContainer는 모의 동작을 수행합니다.
func (c *mockClient) StartContainer(ctx context.Context, containerID string) error {
	return nil
}

// Exec는 명령어 실행 시뮬레이션을 한 결과를 리턴합니다.
func (c *mockClient) Exec(ctx context.Context, containerID string, cmd []string) (string, error) {
	return fmt.Sprintf("[Mock Exec Output] ran command: %s\n", strings.Join(cmd, " ")), nil
}

// Logs는 모의 컨테이너 실행 로그 전체를 가져옵니다.
func (c *mockClient) Logs(ctx context.Context, containerID string) (string, error) {
	return "[Mock Container Logs] Container started successfully.\n", nil
}

// StopContainer는 모의 동작을 수행합니다.
func (c *mockClient) StopContainer(ctx context.Context, containerID string) error {
	return nil
}

// RemoveContainer는 모의 동작을 수행합니다.
func (c *mockClient) RemoveContainer(ctx context.Context, containerID string) error {
	return nil
}
