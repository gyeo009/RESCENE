package main

import (
	"log"
	"os"

	"sandbox/internal/config"
	"sandbox/internal/handler"
	"sandbox/internal/router"
	"sandbox/internal/sandbox"
	"sandbox/internal/sandbox/docker"
	"sandbox/internal/sandbox/mock"
	"sandbox/internal/service"
	"sandbox/internal/workspace"
)

// @title Sandbox API
// @version 1.0
// @description Docker 샌드박스 컨테이너 관리를 위한 확장 가능한 API 서버입니다.
// @host localhost:8080
// @BasePath /
func main() {
	// 1. 설정 정보 로드
	cfg := config.Load()

	// 2. 의존성 주입(Dependency Injection) 초기화
	// 2.1 Workspace Manager 초기화
	workspaceManager := workspace.NewManager(cfg.WorkspacesDir)

	// 2.2 Sandbox Manager 초기화 (Mock 지원)
	var sandboxManager sandbox.Manager
	var err error

	if os.Getenv("USE_MOCK") == "true" {
		log.Println("모의(Mock) Sandbox 매니저가 기동되었습니다 (Docker 데몬에 통신하지 않음).")
		sandboxManager = mock.NewManager()
	} else {
		log.Println("실제 Docker API 기반 Sandbox 매니저가 기동되었습니다.")
		sandboxManager, err = docker.NewManager(cfg)
		if err != nil {
			log.Fatalf("Docker Sandbox 매니저 초기화 실패: %v", err)
		}
	}

	// 2.3 서비스 및 핸들러 초기화
	sandboxService := service.NewSandboxService(sandboxManager, workspaceManager, cfg.HostWorkspacesDir)
	sandboxHandler := handler.NewSandboxHandler(sandboxService)

	// 3. 라우터 설정 초기화
	r := router.NewRouter(sandboxHandler)

	// 4. 웹 서버 시작
	log.Printf("서버를 포트 %s에서 시작하는 중...", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("서버 실행 실패: %v", err)
	}
}
