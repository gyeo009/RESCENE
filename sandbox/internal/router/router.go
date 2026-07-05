package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "sandbox/docs" // Swagger로 생성된 docs 패키지를 등록하기 위해 빈 임포트(blank import)를 수행합니다.
	"sandbox/internal/handler"
)

// NewRouter는 Gin 엔진을 초기화하고 라우트 설정을 수행합니다.
func NewRouter(sandboxHandler *handler.SandboxHandler) *gin.Engine {
	r := gin.Default()

	// Swagger UI 라우트 설정 (/swagger/index.html 경로로 접속 가능)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 샌드박스 관련 라우트 그룹 구성 (추후 추가 그룹 설정이 용이한 구조)
	sandboxGroup := r.Group("/sandbox")
	{
		sandboxGroup.POST("/run", sandboxHandler.RunSandbox)
	}

	return r
}
