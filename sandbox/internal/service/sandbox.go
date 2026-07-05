package service

import (
	"context"

	"sandbox/internal/docker"
	"sandbox/internal/dto"
	"sandbox/internal/response"
)

// SandboxService는 샌드박스 비즈니스 로직을 처리하는 인터페이스입니다.
type SandboxService interface {
	// RunSandbox는 요청받은 시나리오와 사용자 정보를 바탕으로 샌드박스를 실행합니다.
	RunSandbox(ctx context.Context, req *dto.SandboxRunRequest) (*response.SandboxRunResponse, error)
}

type sandboxService struct {
	dockerClient docker.Client
}

// NewSandboxService는 의존성 주입(DI)을 통해 SandboxService 인스턴스를 생성합니다.
func NewSandboxService(dockerClient docker.Client) SandboxService {
	return &sandboxService{
		dockerClient: dockerClient,
	}
}

// RunSandbox는 Docker 클라이언트를 호출하여 컨테이너를 구동하고 실행 결과를 반환합니다.
func (s *sandboxService) RunSandbox(ctx context.Context, req *dto.SandboxRunRequest) (*response.SandboxRunResponse, error) {
	// 향후 검증(Validation), 자원 한도 제한, 데이터베이스 연동 등의 비즈니스 로직이 이곳에 위치합니다.
	info, err := s.dockerClient.CreateSandbox(ctx, docker.CreateSandboxRequest{
		User:     req.User,
		Scenario: req.Scenario,
	})
	if err != nil {
		return nil, err
	}

	return &response.SandboxRunResponse{
		ContainerID: info.ContainerID,
		Status:      info.Status,
	}, nil
}
