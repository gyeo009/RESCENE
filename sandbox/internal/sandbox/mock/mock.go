package mock

import (
	"context"
	"fmt"

	"sandbox/internal/sandbox"
)

type mockSandbox struct {
	id     string
	status string
}

func (s *mockSandbox) ID() string {
	return s.id
}

func (s *mockSandbox) Status() string {
	return s.status
}

func (s *mockSandbox) Exec(ctx context.Context, cmd []string) (string, error) {
	return fmt.Sprintf("mock execution output for cmd: %v", cmd), nil
}

func (s *mockSandbox) Copy(ctx context.Context, content []byte, destPath string) error {
	return nil
}

func (s *mockSandbox) Logs(ctx context.Context) (string, error) {
	return fmt.Sprintf("mock logs for sandbox %s", s.id), nil
}

func (s *mockSandbox) Stop(ctx context.Context) error {
	s.status = "stopped"
	return nil
}

func (s *mockSandbox) Destroy(ctx context.Context) error {
	s.status = "destroyed"
	return nil
}

func (s *mockSandbox) Wait(ctx context.Context) (int64, error) {
	return 0, nil
}

type mockManager struct{}

// NewManager는 모의(Mock) 샌드박스 매니저를 생성합니다.
func NewManager() sandbox.Manager {
	return &mockManager{}
}

func (m *mockManager) Create(ctx context.Context, cfg sandbox.Config) (sandbox.Sandbox, error) {
	return &mockSandbox{
		id:     cfg.ID,
		status: "running",
	}, nil
}

func (m *mockManager) Get(ctx context.Context, id string) (sandbox.Sandbox, error) {
	return &mockSandbox{
		id:     id,
		status: "running",
	}, nil
}
