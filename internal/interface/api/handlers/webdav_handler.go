package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/wentf9/MyGoFileHub/internal/application"
	webdav_adapter "github.com/wentf9/MyGoFileHub/internal/interface/webdav"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/webdav"
)

type WebDAVHandler struct {
	fileService *application.FileService
	authService *application.AuthService
	// 用于存储不同 SourceID 的锁系统
	// Key: sourceID (string), Value: webdav.LockSystem
	lockSystems sync.Map
}

func NewWebDAVHandler(fs *application.FileService, as *application.AuthService) *WebDAVHandler {
	return &WebDAVHandler{fileService: fs, authService: as}
}

// ServeHTTP 处理 WebDAV 请求
// 路由规则: /webdav/:source_id/*path
func (h *WebDAVHandler) Handler(c *gin.Context) {
	// -------------------------------------------------------------
	// 处理尾部斜杠 (Slash Hack)
	// Windows 请求 /webdav/1 -> 我们强制改为 /webdav/1/
	// -------------------------------------------------------------
	sourceID := c.Param("source_id")

	// 获取原始请求路径
	reqPath := c.Request.URL.Path

	// 如果请求路径以 sourceID 结尾（例如 /webdav/1），说明缺少斜杠
	// 我们手动补全它，这样 webdav 库就会认为我们在访问根目录 "/"，而不是空 ""
	if strings.HasSuffix(reqPath, sourceID) {
		c.Request.URL.Path += "/"
	}

	// -------------------------------------------------------------
	// 获取驱动
	// -------------------------------------------------------------
	fmt.Printf("[Debug] Getting driver for sourceID: %s\n", sourceID)
	driver, err := h.fileService.GetDriver(c.Request.Context(), sourceID)
	if err != nil {
		fmt.Printf("[Debug] Driver not found: %v\n", err)
		c.AbortWithStatus(http.StatusNotFound) // <--- 如果是这里报404，说明数据库没查到ID
		return
	}

	// -------------------------------------------------------------
	// 获取或创建单例 LockSystem
	// -------------------------------------------------------------
	// 尝试从 map 中加载这个 sourceID 对应的锁系统
	var lockSystem webdav.LockSystem

	if val, ok := h.lockSystems.Load(sourceID); ok {
		lockSystem = val.(webdav.LockSystem)
	} else {
		// 如果不存在，创建一个新的并存入
		fmt.Printf("[Debug] Creating new LockSystem for SourceID: %s\n", sourceID)
		newLS := webdav.NewMemLS()
		// LoadOrStore 防止并发时的竞争条件
		actual, _ := h.lockSystems.LoadOrStore(sourceID, newLS)
		lockSystem = actual.(webdav.LockSystem)
	}

	// -------------------------------------------------------------
	// 构造 WebDAV 处理逻辑
	// -------------------------------------------------------------
	// 构造前缀：必须以 "/" 结尾，与上面修改过的 Request Path 保持一致
	prefix := "/webdav/" + sourceID + "/"

	handler := &webdav.Handler{
		FileSystem: &webdav_adapter.DriverFileSystem{Driver: driver},
		LockSystem: lockSystem,
		Logger: func(r *http.Request, err error) {
			// 只打印错误，或者打印所有请求以便调试
			if err != nil {
				fmt.Printf("[WebDAV Lib Error] %s %s: %v\n", r.Method, r.URL.Path, err)
			}
		},
	}

	handler.Prefix = prefix

	// 转交处理
	handler.ServeHTTP(c.Writer, c.Request)
}
