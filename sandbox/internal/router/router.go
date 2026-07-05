package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "sandbox/docs" // This import is required to register Swagger documentation files
	"sandbox/internal/handler"
)

// NewRouter initializes the Gin engine and configures routes
func NewRouter(sandboxHandler *handler.SandboxHandler) *gin.Engine {
	r := gin.Default()

	// Swagger UI route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Sandbox routes group (easy to add more route groups later)
	sandboxGroup := r.Group("/sandbox")
	{
		sandboxGroup.POST("/run", sandboxHandler.RunSandbox)
	}

	return r
}
