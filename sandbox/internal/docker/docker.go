package docker

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

// dockerClient는 공식 Docker Go SDK를 활용하여 Docker Engine API와 직접 통신하는 구현체입니다.
type dockerClient struct {
	cli *client.Client
}

// NewDockerClient는 환경변수 설정을 자동으로 읽어들여 Docker Engine API 클라이언트를 초기화합니다.
// API 버전은 클라이언트와 서버 간 자동 조율(API Version Negotiation) 방식을 사용합니다.
func NewDockerClient() (Client, error) {
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return nil, fmt.Errorf("Docker 클라이언트 초기화 실패: %w", err)
	}

	return &dockerClient{
		cli: cli,
	}, nil
}

// CreateContainer는 새 Docker 컨테이너를 생성하고 ID를 반환합니다.
func (c *dockerClient) CreateContainer(ctx context.Context, name string, image string, cmd []string, binds []string, env []string) (string, error) {
	config := &container.Config{
		Image:      image,
		Cmd:        cmd,
		Env:        env,
		WorkingDir: "/workspace",
	}

	hostConfig := &container.HostConfig{
		Binds: binds,
	}

	resp, err := c.cli.ContainerCreate(
		ctx,
		config,
		hostConfig,
		nil,  // 네트워크 설정
		nil,  // 플랫폼 설정
		name, // 고유한 컨테이너 이름 지정
	)
	if err != nil {
		return "", fmt.Errorf("Docker 컨테이너 생성 실패: %w", err)
	}

	return resp.ID, nil
}

// StartContainer는 생성된 컨테이너를 구동합니다.
func (c *dockerClient) StartContainer(ctx context.Context, containerID string) error {
	err := c.cli.ContainerStart(
		ctx,
		containerID,
		container.StartOptions{},
	)
	if err != nil {
		return fmt.Errorf("Docker 컨테이너 시작 실패: %w", err)
	}
	return nil
}

// Exec는 실행 중인 컨테이너 내부에서 명령어를 동기로 수행하고 combined output을 반환합니다.
func (c *dockerClient) Exec(ctx context.Context, containerID string, cmd []string) (string, error) {
	execConfig := container.ExecOptions{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          cmd,
	}

	execCreateResp, err := c.cli.ContainerExecCreate(ctx, containerID, execConfig)
	if err != nil {
		return "", fmt.Errorf("Docker exec 생성 실패: %w", err)
	}

	attachResp, err := c.cli.ContainerExecAttach(ctx, execCreateResp.ID, container.ExecStartOptions{})
	if err != nil {
		return "", fmt.Errorf("Docker exec 연결 실패: %w", err)
	}
	defer attachResp.Close()

	// stdout과 stderr 스트림을 demultiplex(역다중화)하여 읽습니다.
	var stdout, stderr bytes.Buffer
	_, err = stdcopy.StdCopy(&stdout, &stderr, attachResp.Reader)
	if err != nil {
		// 오류 발생 시 단순 스트림 읽기 폴백
		data, readErr := io.ReadAll(attachResp.Reader)
		if readErr == nil {
			return string(data), nil
		}
		return "", fmt.Errorf("Docker exec 결과 읽기 실패: %w", err)
	}

	return stdout.String() + stderr.String(), nil
}

// Logs는 컨테이너의 표준 출력/에러 로그를 가져옵니다.
func (c *dockerClient) Logs(ctx context.Context, containerID string) (string, error) {
	options := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     false,
	}

	reader, err := c.cli.ContainerLogs(ctx, containerID, options)
	if err != nil {
		return "", fmt.Errorf("Docker 컨테이너 로그 조회 실패: %w", err)
	}
	defer reader.Close()

	var stdout, stderr bytes.Buffer
	_, err = stdcopy.StdCopy(&stdout, &stderr, reader)
	if err != nil {
		data, readErr := io.ReadAll(reader)
		if readErr == nil {
			return string(data), nil
		}
		return "", fmt.Errorf("Docker 로그 스트림 읽기 실패: %w", err)
	}

	return stdout.String() + stderr.String(), nil
}

// StopContainer는 작동 중인 컨테이너를 중지시킵니다.
func (c *dockerClient) StopContainer(ctx context.Context, containerID string) error {
	err := c.cli.ContainerStop(ctx, containerID, container.StopOptions{})
	if err != nil {
		return fmt.Errorf("Docker 컨테이너 중지 실패: %w", err)
	}
	return nil
}

// RemoveContainer는 지정된 컨테이너를 강제로 제거하고 관련된 볼륨을 정리합니다.
func (c *dockerClient) RemoveContainer(ctx context.Context, containerID string) error {
	options := container.RemoveOptions{
		Force:         true,
		RemoveVolumes: true,
	}

	err := c.cli.ContainerRemove(ctx, containerID, options)
	if err != nil {
		return fmt.Errorf("Docker 컨테이너 제거 실패: %w", err)
	}
	return nil
}
