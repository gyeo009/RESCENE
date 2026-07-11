package docker

import (
	"context"
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"

	"sandbox/internal/config"
	"sandbox/internal/sandbox"
)

type dockerManager struct {
	cli *client.Client
	cfg *config.Config
}

// NewManager는 Docker API를 사용하는 Sandbox Manager를 생성합니다.
func NewManager(cfg *config.Config) (sandbox.Manager, error) {
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return nil, fmt.Errorf("Docker 클라이언트 초기화 실패: %w", err)
	}

	return &dockerManager{
		cli: cli,
		cfg: cfg,
	}, nil
}

func (m *dockerManager) Create(ctx context.Context, cfg sandbox.Config) (sandbox.Sandbox, error) {
	// 런타임 결정: 설정 기본값 우선 사용 및 config override
	rtName := cfg.Runtime
	if rtName == "" {
		rtName = m.cfg.Runtime
	}

	rt, err := GetRuntime(rtName)
	if err != nil {
		return nil, err
	}

	// 1. 컨테이너 자체 설정
	cmd := []string{"tail", "-f", "/dev/null"}
	containerConfig := &container.Config{
		Image:           cfg.Scenario,
		Cmd:             cmd,
		User:            "1000:1000",
		NetworkDisabled: !cfg.Network,
		WorkingDir:      "/workspace",
	}

	// 2. 호스트 레벨 설정
	hostConfig := &container.HostConfig{
		ReadonlyRootfs: m.cfg.ReadOnlyRootfs,
		Tmpfs: map[string]string{
			"/tmp": "rw,noexec,nosuid,size=" + m.cfg.TmpfsSize,
		},
		Privileged: false,
	}

	if m.cfg.DropCapabilities {
		hostConfig.CapDrop = []string{"ALL"}
	}

	if !cfg.Network {
		hostConfig.NetworkMode = "none"
	}

	// 리소스 제약사항 적용
	var memoryLimit int64
	if cfg.Memory != "" {
		memoryLimit = parseMemoryLimit(cfg.Memory)
	} else {
		memoryLimit = m.cfg.MemoryLimit
	}

	var cpuLimit float64
	if cfg.CPU > 0 {
		cpuLimit = float64(cfg.CPU)
	} else {
		cpuLimit = m.cfg.CPULimit
	}

	hostConfig.Resources = container.Resources{
		NanoCPUs:  int64(cpuLimit * 1e9),
		Memory:    memoryLimit,
		PidsLimit: &m.cfg.PidsLimit,
	}

	if m.cfg.NoNewPrivileges {
		hostConfig.SecurityOpt = []string{"no-new-privileges:true"}
	}

	// 워크스페이스 바인드 설정
	workspacePath := cfg.Workspace
	if workspacePath == "" {
		workspacePath = filepath.ToSlash(filepath.Join(m.cfg.HostWorkspacesDir, cfg.ID))
	}
	hostConfig.Binds = []string{workspacePath + ":/workspace"}

	// 런타임 옵션 적용 (gVisor 등)
	if err := rt.Apply(containerConfig, hostConfig); err != nil {
		return nil, fmt.Errorf("컨테이너 런타임 옵션 적용 실패: %w", err)
	}

	containerName := "sandbox-" + cfg.ID

	// 3. 컨테이너 생성
	resp, err := m.cli.ContainerCreate(
		ctx,
		containerConfig,
		hostConfig,
		nil,
		nil,
		containerName,
	)
	if err != nil {
		return nil, fmt.Errorf("Docker 컨테이너 생성 실패: %w", err)
	}

	// 4. 컨테이너 시작
	err = m.cli.ContainerStart(
		ctx,
		resp.ID,
		container.StartOptions{},
	)
	if err != nil {
		_ = m.cli.ContainerRemove(ctx, resp.ID, container.RemoveOptions{Force: true})
		return nil, fmt.Errorf("Docker 컨테이너 시작 실패: %w", err)
	}

	return &dockerSandbox{
		id:            containerName, // API 호출 시 컨테이너명을 ID로 활용
		status:        "running",
		workspacePath: filepath.Join(m.cfg.WorkspacesDir, cfg.ID),
		cli:           m.cli,
		cfg:           m.cfg,
	}, nil
}

func (m *dockerManager) Get(ctx context.Context, id string) (sandbox.Sandbox, error) {
	containerName := "sandbox-" + id
	return &dockerSandbox{
		id:            containerName,
		status:        "running",
		workspacePath: filepath.Join(m.cfg.WorkspacesDir, id),
		cli:           m.cli,
		cfg:           m.cfg,
	}, nil
}

func parseMemoryLimit(memStr string) int64 {
	if memStr == "" {
		return 0
	}
	memStr = strings.ToLower(strings.TrimSpace(memStr))
	re := regexp.MustCompile(`^(\d+)([kmg]?)$`)
	matches := re.FindStringSubmatch(memStr)
	if len(matches) != 3 {
		return 0
	}
	val, err := strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return 0
	}
	switch matches[2] {
	case "k":
		return val * 1024
	case "m":
		return val * 1024 * 1024
	case "g":
		return val * 1024 * 1024 * 1024
	default:
		return val
	}
}
