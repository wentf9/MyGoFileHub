package handlers

import (
	"net/http"
	"strconv"

	"github.com/wentf9/MyGoFileHub/internal/application"
	"github.com/wentf9/MyGoFileHub/internal/domain/model"
	"github.com/wentf9/MyGoFileHub/internal/infrastructure/drivers"

	"github.com/gin-gonic/gin"
)

type SourceHandler struct {
	service *application.SourceService
}

func NewSourceHandler(s *application.SourceService) *SourceHandler {
	return &SourceHandler{service: s}
}

func (h *SourceHandler) List(c *gin.Context) {
	sources, err := h.service.ListSources(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success", "data": sources})
}

func (h *SourceHandler) GetSchema(c *gin.Context) {
	schemas := drivers.GetRegisteredSchemas()
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success", "data": schemas})
}

func (h *SourceHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid source id"})
		return
	}

	source, err := h.service.GetSource(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success", "data": source})
}

func (h *SourceHandler) Create(c *gin.Context) {
	var source model.StorageSource
	if err := c.ShouldBindJSON(&source); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.CreateSource(c.Request.Context(), &source)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "存储源创建成功", "data": source})
}

func (h *SourceHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid source id"})
		return
	}

	var source model.StorageSource
	if err := c.ShouldBindJSON(&source); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	source.ID = uint(id)

	err = h.service.UpdateSource(c.Request.Context(), &source)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "存储源更新成功", "data": source})
}

func (h *SourceHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid source id"})
		return
	}

	err = h.service.DeleteSource(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "存储源删除成功"})
}
