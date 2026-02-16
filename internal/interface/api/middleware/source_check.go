package middleware

import (
	"net"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/wentf9/MyGoFileHub/config"
)

func ClientCheck() gin.HandlerFunc {
	// 如果白名单配置为 "*"，则允许所有 IP 访问，无需进一步检查
	if config.AppConfig.WhiteListStr == "*" {
		return func(c *gin.Context) {
			c.Next()
		}
	}
	// 预解析白名单，避免每次请求都重新解析，提高性能
	allowedIPs := make(map[string]bool)
	allowedSubnets := []*net.IPNet{}

	ip, subnets, _ := config.ParseIPOrSubnet(config.AppConfig.WhiteListStr)
	for _, ip := range ip {
		allowedIPs[ip.String()] = true
	}
	allowedSubnets = append(allowedSubnets, subnets...)

	return func(c *gin.Context) {
		// 1. 获取客户端真实 IP
		// 注意：在镜像模式/代理环境下，ClientIP() 会处理 X-Forwarded-For
		clientIPStr := getClientIp(c)
		clientIP := net.ParseIP(clientIPStr)

		// 2. 校验是否在单一 IP 白名单中
		if allowedIPs[clientIPStr] {
			c.Next()
			return
		}

		// 3. 校验是否在子网白名单中
		for _, subnet := range allowedSubnets {
			if subnet.Contains(clientIP) {
				c.Next()
				return
			}
		}

		// 4. 校验不通过，拦截请求
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "Access denied: your IP is not in the allowlist",
			"ip":    clientIPStr,
		})
	}

}

func getClientIp(c *gin.Context) string {
	xForwardedFor := c.Request.Header.Get("X-Forwarded-For")
	if xForwardedFor != "" {
		ips := strings.Split(xForwardedFor, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}
	xRealIP := c.Request.Header.Get("X-Real-IP")
	if xRealIP != "" {
		return strings.TrimSpace(xRealIP)
	}
	return c.ClientIP()
}
