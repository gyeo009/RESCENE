package docker

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"

	"sandbox/internal/config"
)

// dockerClient는 공식 Docker Go SDK를 활용하여 Docker Engine API와 직접 통신하는 구현체입니다.
type dockerClient struct {
	cli *client.Client
	cfg *config.Config
}

// NewDockerClient는 설정을 받아 Docker Engine API 클라이언트를 초기화합니다.
// API 버전은 클라이언트와 서버 간 자동 조율(API Version Negotiation) 방식을 사용합니다.
func NewDockerClient(cfg *config.Config) (Client, error) {
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return nil, fmt.Errorf("Docker 클라이언트 초기화 실패: %w", err)
	}

	return &dockerClient{
		cli: cli,
		cfg: cfg,
	}, nil
}

// CreateSandbox는 새로운 Docker 컨테이너를 생성(ContainerCreate)하고 시작(ContainerStart)한 후 메타데이터를 반환합니다.
func (c *dockerClient) CreateSandbox(ctx context.Context, req CreateSandboxRequest) (*SandboxInfo, error) {
	// 런타임 결정: 요청에 명시되어 있지 않으면 설정 기본값 사용
	rtName := req.Runtime
	if rtName == "" {
		rtName = c.cfg.Runtime
	}

	rt, err := GetRuntime(rtName)
	if err != nil {
		return nil, err
	}

	// 1. 컨테이너 자체 설정 (이미지명 등 지정)
	config := &container.Config{
		Image:           req.Scenario,
		User:            "1000:1000", // 비루트(non-root) 사용자 실행
		NetworkDisabled: c.cfg.NetworkDisabled,
		WorkingDir:      "/workspace",
	}

	// 2. 호스트 레벨 설정 (리소스 제약사항 및 보안 옵션 정의)
	hostConfig := &container.HostConfig{
		ReadonlyRootfs: c.cfg.ReadOnlyRootfs,
		Tmpfs: map[string]string{
			"/tmp": "rw,noexec,nosuid,size=" + c.cfg.TmpfsSize,
		},
		Privileged: false,
	}

	// Capabilities 격리 설정
	if c.cfg.DropCapabilities {
		hostConfig.CapDrop = []string{"ALL"}
	}

	// 네트워크 격리 모드 설정
	if c.cfg.NetworkDisabled {
		hostConfig.NetworkMode = "none"
	}

	// 리소스(CPU, Memory, PIDs) 제한 설정
	hostConfig.Resources = container.Resources{
		NanoCPUs:  int64(c.cfg.CPULimit * 1e9),
		Memory:    c.cfg.MemoryLimit,
		PidsLimit: &c.cfg.PidsLimit,
	}

	// 권한 상승 방지 (no-new-privileges)
	if c.cfg.NoNewPrivileges {
		hostConfig.SecurityOpt = []string{"no-new-privileges:true"}
	}

	// 워크스페이스 호스트 볼륨 마운트 설정
	hostWorkspacePath := filepath.Join(c.cfg.WorkspaceBaseDir, req.User)
	// 디렉토리가 호스트/컨테이너 공유 디렉토리에 생성되도록 처리
	if err := os.MkdirAll(hostWorkspacePath, 0755); err != nil {
		return nil, fmt.Errorf("호스트 워크스페이스 디렉토리 생성 실패: %w", err)
	}
	hostConfig.Binds = []string{
		hostWorkspacePath + ":/workspace:rw",
	}

	// 3. 런타임 추상화에 의한 컨테이너 옵션 적용 (gVisor 등)
	if err := rt.Apply(config, hostConfig); err != nil {
		return nil, fmt.Errorf("컨테이너 런타임 옵션 적용 실패: %w", err)
	}

	// 4. 컨테이너 생성
	resp, err := c.cli.ContainerCreate(
		ctx,
		config,
		hostConfig,
		nil,                 // 네트워크 설정
		nil,                 // 플랫폼 설정
		"sandbox-"+req.User, // 고유한 컨테이너 명칭 부여
	)
	if err != nil {
		return nil, fmt.Errorf("Docker 컨테이너 생성 실패: %w", err)
	}

	// 5. 컨테이너 시작
	err = c.cli.ContainerStart(
		ctx,
		resp.ID,
		container.StartOptions{},
	)
	if err != nil {
		// 생성된 컨테이너가 에러로 남아있지 않도록 제거 시도
		_ = c.cli.ContainerRemove(ctx, resp.ID, container.RemoveOptions{Force: true})
		return nil, fmt.Errorf("Docker 컨테이너 시작 실패: %w", err)
	}

	// 6. 컨테이너 정보 반환
	return &SandboxInfo{
		ContainerID: resp.ID,
		Status:      "running",
	}, nil
}

// DestroySandbox는 실행 중인 샌드박스를 종료하고 삭제하며 필요시 워크스페이스를 제거합니다.
func (c *dockerClient) DestroySandbox(ctx context.Context, containerID string, user string) error {
	stopTimeout := 10 // 초
	stopOptions := container.StopOptions{
		Timeout: &stopTimeout,
	}

	// 1. 컨테이너 중지 시도 (이미 중지되었거나 없는 경우도 있으므로 에러는 로깅 후 무시 가능하나, 정밀 제어 수행)
	if err := c.cli.ContainerStop(ctx, containerID, stopOptions); err != nil {
		// 컨테이너가 존재하지 않거나 정지 불능이어도 삭제 단계로 강제 전환할 수 있도록 처리
	}

	// 2. 컨테이너 강제 삭제
	removeOptions := container.RemoveOptions{
		Force:         true,
		RemoveVolumes: true,
	}
	if err := c.cli.ContainerRemove(ctx, containerID, removeOptions); err != nil {
		return fmt.Errorf("컨테이너 삭제 실패: %w", err)
	}

	// 3. 워크스페이스 클린업
	if c.cfg.CleanupWorkspace {
		hostWorkspacePath := filepath.Join(c.cfg.WorkspaceBaseDir, user)
		if err := os.RemoveAll(hostWorkspacePath); err != nil {
			// 워크스페이스 디렉토리 삭제 오류는 경고 처리 (경고 수준 로그)
			fmt.Printf("경고: 워크스페이스 경로 %s 삭제 실패: %v\n", hostWorkspacePath, err)
		}
	}

	return nil
}

// ExecuteCommand는 실행 중인 샌드박스 내부에서 명령을 실행하고 그 결과를 반환합니다.
func (c *dockerClient) ExecuteCommand(ctx context.Context, containerID string, cmd []string) (string, error) {
	execCfg := container.ExecOptions{
		Cmd:          cmd,
		AttachStdout: true,
		AttachStderr: true,
	}

	// 1. Exec 인스턴스 생성
	execResp, err := c.cli.ContainerExecCreate(ctx, containerID, execCfg)
	if err != nil {
		return "", fmt.Errorf("Exec 인스턴스 생성 실패: %w", err)
	}

	// 2. Exec 시작 및 스트림 연결
	resp, err := c.cli.ContainerExecAttach(ctx, execResp.ID, container.ExecStartOptions{})
	if err != nil {
		return "", fmt.Errorf("Exec 실행 실패: %w", err)
	}
	defer resp.Close()

	// 3. 출력 결과 읽기 (stdcopy를 이용해 multiplexed stdout/stderr 분리 독출)
	var outBuf bytes.Buffer
	_, err = stdcopy.StdCopy(&outBuf, &outBuf, resp.Reader)
	if err != nil {
		return "", fmt.Errorf("출력 복사 중 에러 발생: %w", err)
	}

	return outBuf.String(), nil
}

// CopyFileToSandbox는 타르 포맷 스트림을 통해 컨테이너 내부 특정 경로에 파일을 복사합니다.
func (c *dockerClient) CopyFileToSandbox(ctx context.Context, containerID string, content []byte, destPath string) error {
	dir, file := filepath.Split(destPath)
	if dir == "" {
		dir = "/workspace" // 기본 경로는 작업 디렉토리
	}

	tarStream, err := makeTarStream(file, content)
	if err != nil {
		return fmt.Errorf("타르 스트림 생성 실패: %w", err)
	}

	copyOpts := container.CopyToContainerOptions{
		AllowOverwriteDirWithFile: true,
	}

	err = c.cli.CopyToContainer(ctx, containerID, dir, tarStream, copyOpts)
	if err != nil {
		return fmt.Errorf("컨테이너 내부로 파일 복사 실패: %w", err)
	}

	return nil
}

// GetSandboxLogs는 샌드박스의 stdout/stderr 스트림 로그를 전체 조회합니다.
func (c *dockerClient) GetSandboxLogs(ctx context.Context, containerID string) (string, error) {
	logOpts := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     false,
	}

	reader, err := c.cli.ContainerLogs(ctx, containerID, logOpts)
	if err != nil {
		return "", fmt.Errorf("로그 스트림 획득 실패: %w", err)
	}
	defer reader.Close()

	var buf bytes.Buffer
	_, err = stdcopy.StdCopy(&buf, &buf, reader)
	if err != nil {
		return "", fmt.Errorf("로그 출력 변환 실패: %w", err)
	}

	return buf.String(), nil
}

// WaitSandbox는 컨테이너가 중단/종료될 때까지 대기하고 종료 코드를 반환합니다.
func (c *dockerClient) WaitSandbox(ctx context.Context, containerID string) (int64, error) {
	statusCh, errCh := c.cli.ContainerWait(ctx, containerID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return 0, fmt.Errorf("컨테이너 대기 중 에러 발생: %w", err)
		}
	case result := <-statusCh:
		return result.StatusCode, nil
	}
	return 0, fmt.Errorf("대기 중 예기치 않은 상태에 도달했습니다")
}

// makeTarStream은 파일 이름과 바이너리 내용을 포함하는 메모리 상의 tar 스트림을 생성합니다.
func makeTarStream(filename string, content []byte) (io.Reader, error) {
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	hdr := &tar.Header{
		Name: filename,
		Mode: 0644,
		Size: int64(len(content)),
	}
	if err := tw.WriteHeader(hdr); err != nil {
		return nil, err
	}
	if _, err := tw.Write(content); err != nil {
		return nil, err
	}
	if err := tw.Close(); err != nil {
		return nil, err
	}
	return &buf, nil
}
