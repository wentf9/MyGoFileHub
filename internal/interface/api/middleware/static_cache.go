package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// StaticCache 为静态资源设置缓存头
// - JS/CSS/图片/字体：1 年缓存（内容哈希文件名可安全缓存）
// - HTML/JSON：不缓存（确保获取最新版本）
// - 其他：1 天默认缓存
func StaticCache() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		// 根据文件扩展名设置缓存策略
		switch {
		case hasExtension(path, ".js", ".css"):
			c.Header("Cache-Control", "public, max-age=31536000, immutable")
		case hasExtension(path, ".png", ".jpg", ".jpeg", ".gif", ".webp", ".svg", ".ico"):
			c.Header("Cache-Control", "public, max-age=31536000, immutable")
		case hasExtension(path, ".woff", ".woff2", ".ttf", ".eot"):
			c.Header("Cache-Control", "public, max-age=31536000, immutable")
		case hasExtension(path, ".html", ".htm", ".json"):
			c.Header("Cache-Control", "no-cache")
		default:
			c.Header("Cache-Control", "public, max-age=86400")
		}

		c.Next()
	}
}

// hasExtension 检查路径是否有指定的任一扩展名
func hasExtension(path string, extensions ...string) bool {
	lowerPath := strings.ToLower(path)
	for _, ext := range extensions {
		if strings.HasSuffix(lowerPath, ext) {
			return true
		}
	}
	return false
}
