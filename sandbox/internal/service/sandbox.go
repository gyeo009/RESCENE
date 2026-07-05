package service

import (
	"context"

	"sandbox/internal/docker"
	"sandbox/internal/dto"
	"sandbox/internal/response"
)

// SandboxService defines the business logic contract for sandbox operations
type SandboxService interface {
	RunSandbox(ctx context.Context, req *dto.SandboxRunRequest) (*response.SandboxRunResponse, error)
}

type sandboxService struct {
	dockerClient docker.Client
}

// NewSandboxService constructs a new SandboxService with dependency injection
func NewSandboxService(dockerClient docker.Client) SandboxService {
	return &sandboxService{
		dockerClient: dockerClient,
	}
}

// RunSandbox runs the sandbox using the underlying docker client
func (s *sandboxService) RunSandbox(ctx context.Context, req *dto.SandboxRunRequest) (*response.SandboxRunResponse, error) {
	// Future business logic like validation, checking limits, logging, etc. will go here.
	info, err := s.dockerClient.CreateSandbox(ctx, req.Scenario, req.User)
	if err != nil {
		return nil, err
	}

	return &response.SandboxRunResponse{
		ContainerID: info.ContainerID,
		Status:      info.Status,
	}, nil
}
