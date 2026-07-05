package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
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

// CreateSandbox는 새로운 Docker 컨테이너를 생성(ContainerCreate)하고 시작(ContainerStart)한 후 메타데이터를 반환합니다.
func (c *dockerClient) CreateSandbox(ctx context.Context, req CreateSandboxRequest) (*SandboxInfo, error) {
	// 컨테이너 자체 설정 (이미지명 등 지정)
	config := &container.Config{
		Image: req.Scenario,
	}

	// 호스트 레벨 설정 (추후 리소스 제약사항 등을 정의하게 됩니다)
	hostConfig := &container.HostConfig{}

	// 1. 컨테이너 생성
	resp, err := c.cli.ContainerCreate(
		ctx,
		config,
		hostConfig,
		nil, // 네트워크 설정
		nil, // 플랫폼 설정
		"sandbox-"+req.User, // 고유한 컨테이너 명칭 부여
	)
	if err != nil {
		return nil, fmt.Errorf("Docker 컨테이너 생성 실패: %w", err)
	}

	// 2. 컨테이너 시작
	err = c.cli.ContainerStart(
		ctx,
		resp.ID,
		container.StartOptions{},
	)
	if err != nil {
		return nil, fmt.Errorf("Docker 컨테이너 시작 실패: %w", err)
	}

	// 3. 컨테이너 정보 반환
	return &SandboxInfo{
		ContainerID: resp.ID,
		Status:      "running",
	}, nil
}
