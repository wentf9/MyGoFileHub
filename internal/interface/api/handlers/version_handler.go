package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wentf9/MyGoFileHub/internal/application"
)

// VersionHandler 返回版本信息
// @Summary 获取版本信息
// @Description 获取当前构建的版本信息，包括版本号、Git commit、构建时间等
// @Tags version
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /@api/v1/version [GET]
func VersionHandler(c *gin.Context) {
	buildInfo := application.GetBuildInfo()
	c.JSON(http.StatusOK, gin.H{
		"data": buildInfo,
	})
}
