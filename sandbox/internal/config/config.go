package config

import "os"

// Config는 애플리케이션의 설정 정보를 담고 있는 구조체입니다.
type Config struct {
	Port              string // 서버가 수신 대기할 포트 번호
	WorkspacesDir     string // 샌드박스 작업 공간이 저장될 로컬(컨테이너 내부) 디렉토리
	HostWorkspacesDir string // Docker 데몬이 참조할 호스트 기준의 절대 디렉토리 경로
}

// Load는 환경 변수 또는 기본값을 기반으로 설정을 로드하여 반환합니다.
func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // 기본 포트는 8080으로 지정
	}

	workspacesDir := os.Getenv("WORKSPACES_DIR")
	if workspacesDir == "" {
		workspacesDir = "workspaces" // 기본값은 relative workspaces 디렉토리
	}

	hostWorkspacesDir := os.Getenv("HOST_WORKSPACES_DIR")
	if hostWorkspacesDir == "" {
		// 지정되지 않으면 로컬 호스트 단독 실행 상태로 보고 동일하게 매핑
		hostWorkspacesDir = workspacesDir
	}

	return &Config{
		Port:              port,
		WorkspacesDir:     workspacesDir,
		HostWorkspacesDir: hostWorkspacesDir,
	}
}
