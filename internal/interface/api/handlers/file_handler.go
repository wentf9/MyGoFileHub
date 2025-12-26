package handlers

import (
	"net/http"

	"github.com/wentf9/MyGoFileHub/internal/application"

	"github.com/gin-gonic/gin"
)

type FileHandler struct {
	service *application.FileService
}

func NewFileHandler(s *application.FileService) *FileHandler {
	return &FileHandler{service: s}
}

func (h *FileHandler) GetHandler(c *gin.Context) {
	sourceKey := c.Param("source_key")
	path := c.Param("path")
	if path == "" {
		path = "/"
	}
	if sourceKey == "" {
		sources, err := h.service.GetAllSource(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "success",
			"data": gin.H{
				"files": sources,
			},
		})
		return
	}
	pathStat, err := h.service.Stat(c.Request.Context(), sourceKey, path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if pathStat.IsDir {
		h.List(c, sourceKey, path)
	} else {
		// 返回标准化 JSON
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "success",
			"data": gin.H{
				"files": pathStat,
			},
		})
	}
}

func (h *FileHandler) DeleteHandler(c *gin.Context) {
	sourceKey := c.Param("source_key")
	path := c.Param("path")
	if sourceKey == "" || path == "" || path == "/" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不可直接删除根路径"})
		return
	}
	err := h.service.Delete(c.Request.Context(), sourceKey, path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "删除成功",
	})
}

// List 处理 /files/list 请求
// Query Param: source_key, path
func (h *FileHandler) List(c *gin.Context, sourceKey, path string) {
	if sourceKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "source_key is required"})
		return
	}

	// 调用 Application 层
	files, err := h.service.ListFiles(c.Request.Context(), sourceKey, path)
	if err != nil {
		// 实际项目中应根据 error 类型返回 403, 404 或 500
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回标准化 JSON
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{
			"files": files,
		},
	})
}
