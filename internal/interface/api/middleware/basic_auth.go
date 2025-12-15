package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wentf9/MyGoFileHub/internal/application"
)

func BasicAuth(authService *application.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// -------------------------------------------------------------
		// 1. OPTIONS 请求免登处理 (兼容 Windows)
		// -------------------------------------------------------------
		if c.Request.Method == "OPTIONS" {
			if _, _, ok := c.Request.BasicAuth(); !ok {
				fmt.Println("[Debug] Handling OPTIONS bypass for Windows")
				c.Header("DAV", "1, 2")
				c.Header("Allow", "OPTIONS, PROPFIND, PUT, MKCOL, GET, HEAD, DELETE, COPY, MOVE")
				c.Header("MS-Author-Via", "DAV")
				c.Status(http.StatusOK)
				return
			}
		}

		// -------------------------------------------------------------
		// 2. Basic Auth 鉴权
		// -------------------------------------------------------------
		username, password, ok := c.Request.BasicAuth()
		if !ok {
			c.Header("WWW-Authenticate", `Basic realm="GoFile WebDAV"`)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 调用 AuthService
		user, err := authService.LoginBasic(c.Request.Context(), username, password)
		if err != nil {
			fmt.Printf("[Debug] Auth failed for user: %s\n", username)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// 鉴权成功，注入用户信息到 Context
		ctx := c.Request.Context()
		ctx = context.WithValue(c.Request.Context(), "userID", user.ID)
		ctx = context.WithValue(ctx, "username", user.Username)
		ctx = context.WithValue(ctx, "role", user.Role)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
