package api

import (
	"github.com/wentf9/MyGoFileHub/internal/application"
	"github.com/wentf9/MyGoFileHub/internal/interface/api/handlers"
	"github.com/wentf9/MyGoFileHub/internal/interface/api/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter(fileService *application.FileService, authService *application.AuthService) *gin.Engine {
	r := gin.Default()
	// ---------------------------------------------------------
	// 关闭 Gin 的自动重定向
	// Windows WebDAV 客户端不支持在 OPTIONS 请求中遇到 301/307 跳转
	// ---------------------------------------------------------
	r.RedirectTrailingSlash = false
	r.RedirectFixedPath = false
	// 简单的 CORS 中间件（允许前端跨域调试）
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Next()
	})

	// 依赖注入 Handler
	fileHandler := handlers.NewFileHandler(fileService)
	authHandler := handlers.NewAuthHandler(authService)
	webDAVHandler := handlers.NewWebDAVHandler(fileService, authService)

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

	// -------------------------------------------------------------
	// 注册所有 WebDAV 方法
	// -------------------------------------------------------------
	// WebDAV 协议包含许多非标准 HTTP 动词，Gin 的 Any() 不支持它们
	webdavMethods := []string{
		"OPTIONS", "HEAD", "GET", "PUT", "POST", "DELETE",
		"PROPFIND", "PROPPATCH", "MKCOL", "COPY", "MOVE", "LOCK", "UNLOCK",
	}
	// 循环注册所有方法
	for _, method := range webdavMethods {
		// 路由 1: 匹配 /webdav/1/foo
		r.Handle(method, "/webdav/:source_id/*path", webDAVHandler.Handler)

		// 路由 2: 匹配 /webdav/1 (必须单独注册，否则不带斜杠时会 404)
		r.Handle(method, "/webdav/:source_id", webDAVHandler.Handler)
	}
	return r
}
