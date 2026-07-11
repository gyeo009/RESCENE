package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"path/filepath"

	"sandbox/internal/docker"
	"sandbox/internal/dto"
	"sandbox/internal/response"
	"sandbox/internal/workspace"
)

// SandboxService는 호스트 파일 시스템(Workspace)과 Docker 런타임 클라이언트를 조합하여
// 샌드박스의 수명 주기와 파일 관리, 실행 환경을 제공하는 서비스 계층입니다.
type SandboxService interface {
	// Sandbox 수명 주기 제어
	CreateSandbox(ctx context.Context, req *dto.SandboxCreateRequest) (*response.SandboxCreateResponse, error)
	DeleteSandbox(ctx context.Context, id string) error
	StopSandbox(ctx context.Context, id string) error
	GetLogs(ctx context.Context, id string) (*response.SandboxRunResponse, error)
	RunCommand(ctx context.Context, id string, req *dto.SandboxRunRequest) (*response.SandboxRunResponse, error)

	// 파일 CRUD 기능 (호스트 단독 제어)
	ListFiles(id string) (*response.FileListResponse, error)
	ReadFile(id string, filename string) (*response.FileContentResponse, error)
	WriteFile(id string, req *dto.FileWriteRequest) error
	DeleteFile(id string, filename string) error
}

type sandboxService struct {
	dockerClient      docker.Client
	workspaceManager  *workspace.Manager
	hostWorkspacesDir string
}

// NewSandboxService는 의존성 주입을 통해 SandboxService 인스턴스를 생성합니다.
func NewSandboxService(dockerClient docker.Client, workspaceManager *workspace.Manager, hostWorkspacesDir string) SandboxService {
	return &sandboxService{
		dockerClient:      dockerClient,
		workspaceManager:  workspaceManager,
		hostWorkspacesDir: hostWorkspacesDir,
	}
}

// CreateSandbox는 호스트에 작업 디렉토리를 생성하고 기본 파일을 배포한 뒤,
// 해당 디렉토리를 마운트하여 백그라운드에서 동작하는 Docker 컨테이너를 구동합니다.
func (s *sandboxService) CreateSandbox(ctx context.Context, req *dto.SandboxCreateRequest) (*response.SandboxCreateResponse, error) {
	id := generateSandboxID()

	// 1. 호스트 Workspace 생성
	if err := s.workspaceManager.CreateWorkspace(id); err != nil {
		return nil, fmt.Errorf("failed to create host workspace: %w", err)
	}

	// 2. 기본 python 파일 (main.py) 생성
	defaultCode := "print('Hello, Web IDE Sandbox!')\n"
	if err := s.workspaceManager.WriteFile(id, "main.py", []byte(defaultCode)); err != nil {
		_ = s.workspaceManager.RemoveWorkspace(id)
		return nil, fmt.Errorf("failed to create initial workspace file: %w", err)
	}

	// 3. 컨테이너 생성 설정
	containerName := "sandbox-" + id
	// Docker 데몬이 참조해야 하는 호스트 파일시스템 기준의 절대 경로를 빌드 및 매핑
	hostPath := filepath.ToSlash(filepath.Join(s.hostWorkspacesDir, id))
	binds := []string{hostPath + ":/workspace"}

	// 무한 대기 명령어를 지정하여 컨테이너가 즉시 죽지 않고 켜져있게 유지
	cmd := []string{"tail", "-f", "/dev/null"}

	containerID, err := s.dockerClient.CreateContainer(ctx, containerName, req.Scenario, cmd, binds, nil)
	if err != nil {
		_ = s.workspaceManager.RemoveWorkspace(id)
		return nil, fmt.Errorf("failed to create runtime container: %w", err)
	}

	// 4. 컨테이너 시작
	if err := s.dockerClient.StartContainer(ctx, containerID); err != nil {
		_ = s.dockerClient.RemoveContainer(ctx, containerID)
		_ = s.workspaceManager.RemoveWorkspace(id)
		return nil, fmt.Errorf("failed to start runtime container: %w", err)
	}

	return &response.SandboxCreateResponse{
		SandboxID:   id,
		ContainerID: containerID,
	}, nil
}

// DeleteSandbox는 구동 중인 컨테이너를 정지/제거하고 호스트 내 작업 디렉토리도 완전 삭제합니다.
func (s *sandboxService) DeleteSandbox(ctx context.Context, id string) error {
	containerName := "sandbox-" + id

	// 컨테이너 정리 (에러 무시하고 최대한 다 정제)
	_ = s.dockerClient.StopContainer(ctx, containerName)
	_ = s.dockerClient.RemoveContainer(ctx, containerName)

	// 호스트 파일 삭제
	if err := s.workspaceManager.RemoveWorkspace(id); err != nil {
		return fmt.Errorf("failed to clean up host workspace: %w", err)
	}

	return nil
}

// StopSandbox는 구동 중인 컨테이너를 정지합니다.
func (s *sandboxService) StopSandbox(ctx context.Context, id string) error {
	containerName := "sandbox-" + id
	return s.dockerClient.StopContainer(ctx, containerName)
}

// GetLogs는 컨테이너에서 발생한 로그(stdout/stderr)를 조회합니다.
func (s *sandboxService) GetLogs(ctx context.Context, id string) (*response.SandboxRunResponse, error) {
	containerName := "sandbox-" + id
	logs, err := s.dockerClient.Logs(ctx, containerName)
	if err != nil {
		return nil, err
	}
	return &response.SandboxRunResponse{
		Output: logs,
	}, nil
}

// RunCommand는 실행 중인 컨테이너 내부에서 사용자 정의 명령어를 exec로 수행하고 결과를 반환합니다.
func (s *sandboxService) RunCommand(ctx context.Context, id string, req *dto.SandboxRunRequest) (*response.SandboxRunResponse, error) {
	containerName := "sandbox-" + id
	output, err := s.dockerClient.Exec(ctx, containerName, req.Cmd)
	if err != nil {
		return nil, err
	}
	return &response.SandboxRunResponse{
		Output: output,
	}, nil
}

// ListFiles는 특정 샌드박스의 호스트 내 모든 파일 경로를 반환합니다.
func (s *sandboxService) ListFiles(id string) (*response.FileListResponse, error) {
	files, err := s.workspaceManager.ListFiles(id)
	if err != nil {
		return nil, err
	}
	return &response.FileListResponse{
		Files: files,
	}, nil
}

// ReadFile은 특정 샌드박스 내부의 한 파일의 텍스트 내용을 반환합니다.
func (s *sandboxService) ReadFile(id string, filename string) (*response.FileContentResponse, error) {
	content, err := s.workspaceManager.ReadFile(id, filename)
	if err != nil {
		return nil, err
	}
	return &response.FileContentResponse{
		Filename: filename,
		Content:  string(content),
	}, nil
}

// WriteFile은 호스트 작업 공간 내 특정 파일을 수정하거나 신규 작성합니다.
func (s *sandboxService) WriteFile(id string, req *dto.FileWriteRequest) error {
	return s.workspaceManager.WriteFile(id, req.Filename, []byte(req.Content))
}

// DeleteFile은 호스트 작업 공간 내의 파일을 제거합니다.
func (s *sandboxService) DeleteFile(id string, filename string) error {
	return s.workspaceManager.DeleteFile(id, filename)
}

func generateSandboxID() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "sandbox-default-id"
	}
	return hex.EncodeToString(bytes)
}
