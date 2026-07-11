package main

import (
	"log"

	"sandbox/internal/config"
	"sandbox/internal/docker"
	"sandbox/internal/handler"
	"sandbox/internal/router"
	"sandbox/internal/service"
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
	// 실제 Docker 환경과 연동하기 위해 NewDockerClient를 사용합니다.
	dockerClient, err := docker.NewDockerClient(cfg)
	if err != nil {
		log.Fatalf("Docker 클라이언트 초기화 실패: %v", err)
	}
	sandboxService := service.NewSandboxService(dockerClient)
	sandboxHandler := handler.NewSandboxHandler(sandboxService)

	// 3. 라우터 설정 초기화
	r := router.NewRouter(sandboxHandler)

	// 4. 웹 서버 시작
	log.Printf("서버를 포트 %s에서 시작하는 중...", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("서버 실행 실패: %v", err)
	}
}
