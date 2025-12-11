package api

import (
	"github.com/wentf9/MyGoFileHub/internal/application"
	"github.com/wentf9/MyGoFileHub/internal/interface/api/handlers"
	"github.com/wentf9/MyGoFileHub/internal/interface/api/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter(fileService *application.FileService, authService *application.AuthService) *gin.Engine {
	r := gin.Default()

	// 简单的 CORS 中间件（允许前端跨域调试）
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Next()
	})

	// 依赖注入 Handler
	fileHandler := handlers.NewFileHandler(fileService)
	authHandler := handlers.NewAuthHandler(authService)

	// API 版本控制
	v1 := r.Group("/api/v1")
	{
		// 公开接口
		v1.POST("/login", authHandler.Login)
		// 保护接口 (使用 JWTAuth 中间件)
		protected := v1.Group("/")
		protected.Use(middleware.JWTAuth())
		{
			files := protected.Group("/files")
			files.GET("/list", fileHandler.List)
		}
	}

	return r
}
