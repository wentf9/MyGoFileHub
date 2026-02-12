package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AdminOnly 限制仅管理员访问
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.Request.Context().Value("role")
		if role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"code": 403,
				"msg":  "权限不足，仅限管理员访问",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
