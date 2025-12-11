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

// List 处理 /files/list 请求
// Query Param: source_id, path
func (h *FileHandler) List(c *gin.Context) {
	sourceID := c.Query("source_id")
	path := c.DefaultQuery("path", "/") // 默认为根目录

	if sourceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "source_id is required"})
		return
	}

	// 调用 Application 层
	files, err := h.service.ListFiles(c.Request.Context(), sourceID, path)
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
