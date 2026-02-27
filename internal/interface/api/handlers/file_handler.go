package handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/wentf9/MyGoFileHub/internal/application"
	"github.com/wentf9/MyGoFileHub/internal/domain/model"
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

	// 检查是否是下载请求
	if c.Query("download") == "true" {
		h.downloadStream(c, sourceKey, path)
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

func (h *FileHandler) downloadStream(c *gin.Context, sourceKey, path string) {
	stat, err := h.service.Stat(c.Request.Context(), sourceKey, path)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}
	if stat.IsDir {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot download a directory"})
		return
	}

	stream, err := h.service.GetFileStream(c.Request.Context(), sourceKey, path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer stream.Close()

	fileName := url.PathEscape(stat.Name)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename*=UTF-8''%s", fileName))
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Length", fmt.Sprintf("%d", stat.Size))

	c.DataFromReader(http.StatusOK, stat.Size, "application/octet-stream", stream, nil)
}

// 修改 GetHandler 使用 downloadStream
func (h *FileHandler) PostHandler(c *gin.Context) {
	sourceKey := c.Param("source_key")
	path := c.Param("path")
	if sourceKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "source_key is required"})
		return
	}

	// 检查是否是创建目录
	isDir := c.Query("type") == "dir"
	if isDir {
		err := h.service.Mkdir(c.Request.Context(), sourceKey, path)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "目录创建成功"})
		return
	}

	// 上传文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required in multipart form"})
		return
	}

	// 验证文件名
	if err := model.ValidateFileName(file.Filename); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	openedFile, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer openedFile.Close()

	// 如果 path 以 / 结尾，则自动拼接文件名
	if strings.HasSuffix(path, "/") || path == "" {
		path = path + file.Filename
	}

	err = h.service.CreateFile(c.Request.Context(), sourceKey, path, openedFile, file.Size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "文件上传成功"})
}

func (h *FileHandler) PutHandler(c *gin.Context) {
	sourceKey := c.Param("source_key")
	path := c.Param("path")
	if sourceKey == "" || path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "source_key and path are required"})
		return
	}

	var req struct {
		NewPath string `json:"new_path"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if req.NewPath == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "new_path is required"})
		return
	}

	// 规范化 NewPath: 去除前端可能带来的 /source_key 前缀
	newPath := req.NewPath
	prefix := "/" + sourceKey + "/"
	if after, ok := strings.CutPrefix(newPath, prefix); ok {
		newPath = "/" + after
	} else if newPath == "/"+sourceKey {
		newPath = "/"
	}

	err := h.service.Rename(c.Request.Context(), sourceKey, path, newPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "重命名成功"})
}

func (h *FileHandler) ActionHandler(c *gin.Context, action string) {
	sourcePath := c.Param("path")
	sourceKey := c.Param("source_key")
	dest := c.Query("dest")

	if sourceKey == "" || sourcePath == "" || dest == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "source_key, path and dest are required"})
		return
	}

	// 规范化 dest: 如果 dest 是 /source_key/path 格式，去除前面的 /source_key
	prefix := "/" + sourceKey + "/"
	if strings.HasPrefix(dest, prefix) {
		dest = "/" + strings.TrimPrefix(dest, prefix)
	} else if dest == "/"+sourceKey {
		dest = "/"
	}

	var err error
	switch action {
	case "cp":
		err = h.service.Copy(c.Request.Context(), sourceKey, sourcePath, dest)
	case "mv":
		err = h.service.Rename(c.Request.Context(), sourceKey, sourcePath, dest)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported action"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "操作成功"})
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
