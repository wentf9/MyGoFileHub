package config

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	ServerPort      string `env:"MY_GO_FILE_HUB_SERVER_PORT" envDefault:"3939"`
	Listen          string `env:"MY_GO_FILE_HUB_LISTEN" envDefault:"localhost"`
	DataDir         string `env:"MY_GO_FILE_HUB_DATA_DIR" envDefault:"./data"`
	WhiteListStr    string `env:"MY_GO_FILE_HUB_WHITE_LIST" envDefault:"127.0.0.1"`
	WhiteListIP     []net.IP
	WhiteListSubnet []*net.IPNet
}

var AppConfig Config

// ParseIPOrSubnet 解析逗号分隔的字符串，返回有效的 IP 和 IPNet(子网)
func ParseIPOrSubnet(input string) ([]net.IP, []*net.IPNet, error) {
	var ips []net.IP
	var subnets []*net.IPNet

	// 1. 按逗号分割并清理空格
	parts := strings.SplitSeq(input, ",")

	for part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}

		// 2. 尝试解析为 CIDR (如 10.0.0.1/24)
		if strings.Contains(trimmed, "/") {
			_, ipNet, err := net.ParseCIDR(trimmed)
			if err != nil {
				return nil, nil, fmt.Errorf("无效的子网格式: %s", trimmed)
			}
			subnets = append(subnets, ipNet)
		} else {
			// 3. 尝试解析为纯 IP (如 192.168.1.1)
			ip := net.ParseIP(trimmed)
			if ip == nil {
				return nil, nil, fmt.Errorf("无效的 IP 格式: %s", trimmed)
			}
			ips = append(ips, ip)
		}
	}

	return ips, subnets, nil
}

func init() {
	if err := env.Parse(&AppConfig); err != nil {
		panic(err)
	}
	// 解析白名单
	if AppConfig.WhiteListStr == "" {
		AppConfig.WhiteListStr = "127.0.0.1"
	}
	if AppConfig.WhiteListStr != "*" {
		_, _, err := ParseIPOrSubnet(AppConfig.WhiteListStr)
		if err != nil {
			panic(err)
		}
	}
	if AppConfig.Listen == "localhost" {
		AppConfig.Listen = "127.0.0.1"
	}
	if net.ParseIP(AppConfig.Listen) == nil {
		panic("Invalid listen address: " + AppConfig.Listen)
	}
	if port, err := strconv.ParseUint(AppConfig.ServerPort, 10, 16); err != nil || port == 0 {
		panic("Invalid server port: " + AppConfig.ServerPort)
	}
}
