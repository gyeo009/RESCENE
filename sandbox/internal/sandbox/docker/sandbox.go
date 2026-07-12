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

type dockerSandbox struct {
	id            string
	status        string
	workspacePath string
	cli           *client.Client
	cfg           *config.Config
}

func (s *dockerSandbox) ID() string {
	return s.id
}

func (s *dockerSandbox) Status() string {
	return s.status
}

func (s *dockerSandbox) Exec(ctx context.Context, cmd []string) (string, error) {
	execConfig := container.ExecOptions{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          cmd,
	}

	execCreateResp, err := s.cli.ContainerExecCreate(ctx, s.id, execConfig)
	if err != nil {
		return "", fmt.Errorf("Docker exec 생성 실패: %w", err)
	}

	attachResp, err := s.cli.ContainerExecAttach(ctx, execCreateResp.ID, container.ExecStartOptions{})
	if err != nil {
		return "", fmt.Errorf("Docker exec 연결 실패: %w", err)
	}
	defer attachResp.Close()

	var stdout, stderr bytes.Buffer
	_, err = stdcopy.StdCopy(&stdout, &stderr, attachResp.Reader)
	if err != nil {
		data, readErr := io.ReadAll(attachResp.Reader)
		if readErr == nil {
			return string(data), nil
		}
		return "", fmt.Errorf("Docker exec 결과 읽기 실패: %w", err)
	}

	return stdout.String() + stderr.String(), nil
}

func (s *dockerSandbox) Copy(ctx context.Context, content []byte, destPath string) error {
	dir, file := filepath.Split(destPath)
	if dir == "" {
		dir = "/workspace"
	}

	tarStream, err := makeTarStream(file, content)
	if err != nil {
		return fmt.Errorf("타르 스트림 생성 실패: %w", err)
	}

	copyOpts := container.CopyToContainerOptions{
		AllowOverwriteDirWithFile: true,
	}

	err = s.cli.CopyToContainer(ctx, s.id, dir, tarStream, copyOpts)
	if err != nil {
		return fmt.Errorf("컨테이너 내부로 파일 복사 실패: %w", err)
	}

	return nil
}

func (s *dockerSandbox) Logs(ctx context.Context) (string, error) {
	options := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     false,
	}

	reader, err := s.cli.ContainerLogs(ctx, s.id, options)
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

func (s *dockerSandbox) Stop(ctx context.Context) error {
	err := s.cli.ContainerStop(ctx, s.id, container.StopOptions{})
	if err != nil {
		return fmt.Errorf("Docker 컨테이너 중지 실패: %w", err)
	}
	s.status = "stopped"
	return nil
}

func (s *dockerSandbox) Destroy(ctx context.Context) error {
	// 1. 컨테이너 중지 시도 (에러 무시 가능)
	_ = s.cli.ContainerStop(ctx, s.id, container.StopOptions{})

	// 2. 컨테이너 삭제
	options := container.RemoveOptions{
		Force:         true,
		RemoveVolumes: true,
	}
	if err := s.cli.ContainerRemove(ctx, s.id, options); err != nil {
		return fmt.Errorf("Docker 컨테이너 제거 실패: %w", err)
	}
	s.status = "destroyed"

	// 3. 호스트 파일 제거 (workspacePath가 지정된 경우)
	if s.workspacePath != "" {
		if err := os.RemoveAll(s.workspacePath); err != nil {
			fmt.Printf("경고: 워크스페이스 경로 %s 삭제 실패: %v\n", s.workspacePath, err)
		}
	}

	return nil
}

func (s *dockerSandbox) Wait(ctx context.Context) (int64, error) {
	statusCh, errCh := s.cli.ContainerWait(ctx, s.id, container.WaitConditionNotRunning)
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
