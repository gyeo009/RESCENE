package config

import "os"

// Config는 애플리케이션의 설정 정보를 담고 있는 구조체입니다.
type Config struct {
	Port string // 서버가 수신 대기할 포트 번호
}

// Load는 환경 변수 또는 기본값을 기반으로 설정을 로드하여 반환합니다.
func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // 기본 포트는 8080으로 지정
	}
	return &Config{
		Port: port,
	}
}
