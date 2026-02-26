package api

import (
	"embed"
	"io"
	"io/fs"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/wentf9/MyGoFileHub/config"
	"github.com/wentf9/MyGoFileHub/internal/application"
	"github.com/wentf9/MyGoFileHub/internal/interface/api/handlers"
	"github.com/wentf9/MyGoFileHub/internal/interface/api/middleware"
)

// InitRouter 初始化路由
// staticFS: 嵌入式静态文件系统，为 nil 时表示开发模式 (不服务静态文件)
func InitRouter(
	fileService *application.FileService,
	authService *application.AuthService,
	userService *application.UserService,
	sourceService *application.SourceService,
	staticFS fs.FS,
) *gin.Engine {
	r := gin.Default()

	// ---------------------------------------------------------
	// 关闭 Gin 的自动重定向
	// Windows WebDAV 客户端不支持在 OPTIONS 请求中遇到 301/307 跳转
	// ---------------------------------------------------------
	// r.RedirectTrailingSlash = false
	// r.RedirectFixedPath = false

	// 简单的 CORS 中间件（允许前端跨域调试）
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Next()
	})

	// 依赖注入 Handler
	fileHandler := handlers.NewFileHandler(fileService)
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	sourceHandler := handlers.NewSourceHandler(sourceService)
	webDAVHandler := handlers.NewWebDAVHandler(fileService, authService)
	versionHandler := handlers.VersionHandler

	// ---------------------------------------------------------
	// 服务嵌入式静态文件 (生产模式)
	// 前端应用通过 /ui 路径访问
	// ---------------------------------------------------------
	if staticFS != nil {
		// 服务前端静态文件到 /ui 路径下
		r.StaticFS("/ui", http.FS(staticFS))
	}

	// ---------------------------------------------------------
	// 文件操作路由 (根路径)
	// ---------------------------------------------------------
	file := r.Group("/")
	file.Use(middleware.ClientCheck(), middleware.JWTAuth())
	{
		file.GET("/", fileHandler.GetHandler)
		file.GET("/:source_key/*path", fileHandler.GetHandler)
		file.POST("/:source_key/*path", fileHandler.PostHandler)
		file.PUT("/:source_key/*path", fileHandler.PutHandler)
		file.DELETE("/:source_key/*path", fileHandler.DeleteHandler)

		// 复杂操作：/@cp/source_key/path?dest=...
		file.POST("/@cp/:source_key/*path", func(c *gin.Context) {
			fileHandler.ActionHandler(c, "cp")
		})
		file.POST("/@mv/:source_key/*path", func(c *gin.Context) {
			fileHandler.ActionHandler(c, "mv")
		})
	}

	// API 版本控制
	v1 := r.Group("/@api/v1")
	{
		// 公开接口
		v1.POST("/login", authHandler.Login).Use(middleware.ClientCheck())

		// 版本信息 (公开接口)
		v1.GET("/version", versionHandler)

		// 用户管理 (通常仅限 admin)
		users := v1.Group("/users")
		users.Use(middleware.ClientCheck(), middleware.JWTAuth(), middleware.AdminOnly())
		{
			users.GET("/", userHandler.List)
			users.POST("/", userHandler.Create)
			users.PUT("/:id", userHandler.Update)
			users.DELETE("/:id", userHandler.Delete)
		}

		// 存储源管理 (仅限 admin)
		sources := v1.Group("/sources")
		sources.Use(middleware.ClientCheck(), middleware.JWTAuth(), middleware.AdminOnly())
		{
			sources.GET("/schema", sourceHandler.GetSchema)
			sources.GET("/", sourceHandler.List)
			sources.GET("/:id", sourceHandler.Get)
			sources.POST("/", sourceHandler.Create)
			sources.PUT("/:id", sourceHandler.Update)
			sources.DELETE("/:id", sourceHandler.Delete)
		}
	}

	// -------------------------------------------------------------
	// 注册所有 WebDAV 方法
	// -------------------------------------------------------------
	webdavMethods := []string{
		"OPTIONS", "HEAD", "GET", "PUT", "POST", "DELETE",
		"PROPFIND", "PROPPATCH", "MKCOL", "COPY", "MOVE", "LOCK", "UNLOCK",
	}
	webdav := r.Group("/webdav")
	{
		webdav.Use(middleware.ClientCheck(), middleware.BasicAuth(authService))
		for _, method := range webdavMethods {
			webdav.Handle(method, "/:source_key/*path", webDAVHandler.Handler)
			webdav.Handle(method, "/:source_key", webDAVHandler.Handler)
		}
	}

	// -------------------------------------------------------------
	// HTML5 History API 支持 - 将 /ui 的 SPA 路由回退到 /ui/index.html
	// -------------------------------------------------------------
	if staticFS != nil {
		// 处理 /ui 前端的 SPA 路由
		r.NoRoute(func(c *gin.Context) {
			path := c.Request.URL.Path
			method := c.Request.Method

			// 只处理 GET 请求
			if method != http.MethodGet {
				c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
				return
			}

			// API 和 WebDAV 请求不处理
			if strings.HasPrefix(path, "/@api/") ||
				strings.HasPrefix(path, "/webdav/") ||
				strings.HasPrefix(path, "/@") {
				c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
				return
			}

			// 对于 /ui 路径的 SPA 路由，返回 index.html
			if strings.HasPrefix(path, "/ui") || path == "/" {
				var indexFile []byte
				readFS, ok := staticFS.(interface{ ReadFile(string) ([]byte, error) })
				if ok {
					indexFile, _ = readFS.ReadFile("index.html")
				} else {
					f, err := staticFS.Open("index.html")
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load index.html"})
						return
					}
					defer f.Close()
					indexFile, _ = io.ReadAll(f)
				}

				if len(indexFile) == 0 {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load index.html"})
					return
				}
				c.Data(http.StatusOK, "text/html; charset=utf-8", indexFile)
				return
			}

			c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		})
	}

	// -------------------------------------------------------------
	// 配置信任代理
	// -------------------------------------------------------------
	ips, subnets, err := config.ParseIPOrSubnet(config.AppConfig.WhiteListStr)
	if err != nil {
		r.SetTrustedProxies(nil)
		return r
	}

	var trustedList []string
	for _, ip := range ips {
		trustedList = append(trustedList, ip.String())
	}
	for _, sn := range subnets {
		trustedList = append(trustedList, sn.String())
	}

	if len(trustedList) > 0 {
		r.SetTrustedProxies(trustedList)
	} else {
		r.SetTrustedProxies(nil)
	}

	return r
}

// ensure embed.FS implements fs.FS
var _ fs.FS = (*embed.FS)(nil)
