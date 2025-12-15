package middleware

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/wentf9/MyGoFileHub/config"
)

var localSubnets []net.Addr

func ClientCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		if config.AppConfig.LanOnly == "false" {
			c.Next()
			return
		}
		clientIp := getClientIp(c)

		if ok, err := isSameSubnet(clientIp); err != nil || !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: unable to determine client IP"})
			c.Abort()
			return
		}
		// 继续处理请求
		c.Next()
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

// 判断 targetIP 是否与本机在同一个子网
func isSameSubnet(targetIPStr string) (bool, error) {
	targetIP := net.ParseIP(targetIPStr)
	if targetIP == nil {
		return false, fmt.Errorf("invalid target IP: %s", targetIPStr)
	}
	for _, addr := range localSubnets {
		switch v := addr.(type) {
		case *net.IPNet:
			// 只处理 IPv4
			localIP := v.IP.To4()
			if localIP == nil {
				continue // 不是 IPv4
			}

			// 检查目标 IP 是否属于当前接口的子网
			if v.Contains(targetIP) {
				return true, nil
			}
		}
	}

	return false, nil
}

func init() {
	// 获取所有网络接口
	interfaces, err := net.Interfaces()
	if err != nil {
		panic("获取本机网络接口失败!")
	}

	for _, iface := range interfaces {
		// 跳过非活动或回环接口
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			localSubnets = append(localSubnets, addr)
		}
	}
}
