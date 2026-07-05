package docker

import "context"

// SandboxInfo contains details of a managed sandbox container
type SandboxInfo struct {
	ContainerID string
	Status      string
}

// Client defines the interface for interacting with Docker (or mock docker)
type Client interface {
	CreateSandbox(ctx context.Context, scenario string, user string) (*SandboxInfo, error)
}

// mockClient implements the Client interface for development and testing without real Docker
type mockClient struct{}

// NewMockClient creates a new mock Client instance
func NewMockClient() Client {
	return &mockClient{}
}

// CreateSandbox simulates container creation and returns mock container information
func (c *mockClient) CreateSandbox(ctx context.Context, scenario string, user string) (*SandboxInfo, error) {
	// In the future, this will use the official Docker SDK.
	// For now, it returns mock data to satisfy requirements.
	return &SandboxInfo{
		ContainerID: "sandbox-" + user,
		Status:      "running",
	}, nil
}
