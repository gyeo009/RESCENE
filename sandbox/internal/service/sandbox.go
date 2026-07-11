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
	// DestroySandbox는 실행 중인 샌드박스 컨테이너를 종료하고 자원을 제거합니다.
	DestroySandbox(ctx context.Context, req *dto.SandboxDestroyRequest) error
	// ExecuteCommand는 실행 중인 컨테이너 내부에서 명령을 실행하고 그 결과를 수집합니다.
	ExecuteCommand(ctx context.Context, req *dto.SandboxExecRequest) (*response.SandboxExecResponse, error)
	// CopyFile은 컨테이너 내부 특정 경로에 파일을 복사합니다.
	CopyFile(ctx context.Context, req *dto.SandboxCopyRequest) error
	// GetLogs는 컨테이너의 stdout/stderr 실행 로그 전체를 조회합니다.
	GetLogs(ctx context.Context, containerID string) (*response.SandboxLogsResponse, error)
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
	info, err := s.dockerClient.CreateSandbox(ctx, docker.CreateSandboxRequest{
		User:     req.User,
		Scenario: req.Scenario,
		Runtime:  req.Runtime,
	})
	if err != nil {
		return nil, err
	}

	return &response.SandboxRunResponse{
		ContainerID: info.ContainerID,
		Status:      info.Status,
	}, nil
}

// DestroySandbox는 컨테이너를 종료하고 삭제합니다.
func (s *sandboxService) DestroySandbox(ctx context.Context, req *dto.SandboxDestroyRequest) error {
	return s.dockerClient.DestroySandbox(ctx, req.ContainerID, req.User)
}

// ExecuteCommand는 컨테이너 내부에서 명령어를 수행하고 출력을 반환합니다.
func (s *sandboxService) ExecuteCommand(ctx context.Context, req *dto.SandboxExecRequest) (*response.SandboxExecResponse, error) {
	out, err := s.dockerClient.ExecuteCommand(ctx, req.ContainerID, req.Cmd)
	if err != nil {
		return nil, err
	}
	return &response.SandboxExecResponse{
		Output: out,
	}, nil
}

// CopyFile은 지정한 텍스트 파일 내용을 컨테이너 내부에 기록합니다.
func (s *sandboxService) CopyFile(ctx context.Context, req *dto.SandboxCopyRequest) error {
	return s.dockerClient.CopyFileToSandbox(ctx, req.ContainerID, []byte(req.Content), req.DestPath)
}

// GetLogs는 컨테이너의 표준 로그 기록을 조회합니다.
func (s *sandboxService) GetLogs(ctx context.Context, containerID string) (*response.SandboxLogsResponse, error) {
	logs, err := s.dockerClient.GetSandboxLogs(ctx, containerID)
	if err != nil {
		return nil, err
	}
	return &response.SandboxLogsResponse{
		Logs: logs,
	}, nil
}
