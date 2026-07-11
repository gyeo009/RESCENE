package config

import (
	"os"
	"strconv"
)

// Config는 애플리케이션의 설정 정보를 담고 있는 구조체입니다.
type Config struct {
	Port              string  // 서버가 수신 대기할 포트 번호
	WorkspacesDir     string  // 샌드박스 작업 공간이 저장될 로컬(컨테이너 내부) 디렉토리
	HostWorkspacesDir string  // Docker 데몬이 참조할 호스트 기준의 절대 디렉토리 경로
	Runtime           string  // 기본 컨테이너 런타임 (예: runc, runsc, kata)
	CleanupWorkspace  bool    // 컨테이너 종료 시 워크스페이스 디렉토리 삭제 여부
	ReadOnlyRootfs    bool    // 컨테이너 루트 파일시스템 읽기 전용 적용 여부
	TmpfsSize         string  // /tmp용 tmpfs 용량 설정 (예: 64m)
	DropCapabilities  bool    // 모든 리눅스 Capability 드롭 여부
	NetworkDisabled   bool    // 네트워크 격리 적용 여부
	CPULimit          float64 // CPU 제한 (예: 1.0 = 1 CPU)
	MemoryLimit       int64   // 메모리 제한 바이트 (예: 536870912 = 512MB)
	PidsLimit         int64   // 최대 프로세스 개수 제한
	NoNewPrivileges   bool    // no-new-privileges 권한 상승 제한 적용 여부
}

// Load는 환경 변수 또는 기본값을 기반으로 설정을 로드하여 반환합니다.
func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	workspacesDir := os.Getenv("WORKSPACES_DIR")
	if workspacesDir == "" {
		workspacesDir = "workspaces"
	}

	hostWorkspacesDir := os.Getenv("HOST_WORKSPACES_DIR")
	if hostWorkspacesDir == "" {
		hostWorkspacesDir = workspacesDir
	}

	runtime := os.Getenv("SANDBOX_RUNTIME")
	if runtime == "" {
		runtime = "runc"
	}

	cleanupWorkspace := true
	if os.Getenv("SANDBOX_CLEANUP_WORKSPACE") == "false" {
		cleanupWorkspace = false
	}

	readOnlyRootfs := true
	if os.Getenv("SANDBOX_READ_ONLY_ROOTFS") == "false" {
		readOnlyRootfs = false
	}

	tmpfsSize := os.Getenv("SANDBOX_TMPFS_SIZE")
	if tmpfsSize == "" {
		tmpfsSize = "64m"
	}

	dropCapabilities := true
	if os.Getenv("SANDBOX_DROP_CAPABILITIES") == "false" {
		dropCapabilities = false
	}

	networkDisabled := true
	if os.Getenv("SANDBOX_NETWORK_DISABLED") == "false" {
		networkDisabled = false
	}

	cpuLimit := 1.0
	if val := os.Getenv("SANDBOX_CPU_LIMIT"); val != "" {
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			cpuLimit = f
		}
	}

	var memoryLimit int64 = 512 * 1024 * 1024 // 512MB
	if val := os.Getenv("SANDBOX_MEMORY_LIMIT"); val != "" {
		if i, err := strconv.ParseInt(val, 10, 64); err == nil {
			memoryLimit = i
		}
	}

	var pidsLimit int64 = 100
	if val := os.Getenv("SANDBOX_PIDS_LIMIT"); val != "" {
		if i, err := strconv.ParseInt(val, 10, 64); err == nil {
			pidsLimit = i
		}
	}

	noNewPrivileges := true
	if os.Getenv("SANDBOX_NO_NEW_PRIVILEGES") == "false" {
		noNewPrivileges = false
	}

	return &Config{
		Port:              port,
		WorkspacesDir:     workspacesDir,
		HostWorkspacesDir: hostWorkspacesDir,
		Runtime:           runtime,
		CleanupWorkspace:  cleanupWorkspace,
		ReadOnlyRootfs:    readOnlyRootfs,
		TmpfsSize:         tmpfsSize,
		DropCapabilities:  dropCapabilities,
		NetworkDisabled:   networkDisabled,
		CPULimit:          cpuLimit,
		MemoryLimit:       memoryLimit,
		PidsLimit:         pidsLimit,
		NoNewPrivileges:   noNewPrivileges,
	}
}
