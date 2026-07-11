package docker

import (
	"fmt"

	"github.com/docker/docker/api/types/container"
)

// SandboxRuntime은 Docker 샌드박스의 컨테이너 런타임 설정을 추상화합니다.
type SandboxRuntime interface {
	Name() string
	Apply(config *container.Config, hostConfig *container.HostConfig) error
}

// RuncRuntime은 기본 OCI 런타임(runc)입니다.
type RuncRuntime struct{}

func (r *RuncRuntime) Name() string {
	return "runc"
}

func (r *RuncRuntime) Apply(config *container.Config, hostConfig *container.HostConfig) error {
	hostConfig.Runtime = "" // 기본 런타임(공백 또는 runc) 사용
	return nil
}

// RunscRuntime은 gVisor 보안 샌드박스 런타임(runsc)입니다.
type RunscRuntime struct{}

func (r *RunscRuntime) Name() string {
	return "runsc"
}

func (r *RunscRuntime) Apply(config *container.Config, hostConfig *container.HostConfig) error {
	hostConfig.Runtime = "runsc"
	return nil
}

// KataRuntime은 Kata Containers 가상 머신 기반 런타임입니다.
type KataRuntime struct{}

func (r *KataRuntime) Name() string {
	return "kata"
}

func (r *KataRuntime) Apply(config *container.Config, hostConfig *container.HostConfig) error {
	hostConfig.Runtime = "kata"
	return nil
}

// GetRuntime은 이름에 매칭되는 SandboxRuntime 인스턴스를 반환합니다.
func GetRuntime(name string) (SandboxRuntime, error) {
	switch name {
	case "runc":
		return &RuncRuntime{}, nil
	case "runsc":
		return &RunscRuntime{}, nil
	case "kata":
		return &KataRuntime{}, nil
	default:
		return nil, fmt.Errorf("지원되지 않는 샌드박스 런타임입니다: %s", name)
	}
}
