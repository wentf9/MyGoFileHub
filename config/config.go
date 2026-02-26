package config

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	ServerPort      string `env:"MY_GO_FILE_HUB_SERVER_PORT" envDefault:"3939"`
	Listen          string `env:"MY_GO_FILE_HUB_LISTEN" envDefault:"localhost"`
	DataDir         string `env:"MY_GO_FILE_HUB_DATA_DIR" envDefault:"./data"`
	WhiteListStr    string `env:"MY_GO_FILE_HUB_WHITE_LIST" envDefault:"127.0.0.1"`
	SecretKey       string `env:"MY_GO_FILE_HUB_SECRET_KEY"`                             // 用于加密存储源配置的密钥，32 字节
	Mode            string `env:"MY_GO_FILE_HUB_MODE" envDefault:"prod"`                 // 运行模式：dev(开发) / prod(生产)
	FrontendDir     string `env:"MY_GO_FILE_HUB_FRONTEND_DIR" envDefault:"./frontend"`   // 开发模式下前端目录路径
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
				return nil, nil, fmt.Errorf("无效的子网格式：%s", trimmed)
			}
			subnets = append(subnets, ipNet)
		} else {
			// 3. 尝试解析为纯 IP (如 192.168.1.1)
			ip := net.ParseIP(trimmed)
			if ip == nil {
				return nil, nil, fmt.Errorf("无效的 IP 格式：%s", trimmed)
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

	// 初始化 DataDir 目录
	if err := os.MkdirAll(AppConfig.DataDir, 0755); err != nil {
		panic("Failed to create data dir: " + err.Error())
	}

	// 初始化 SecretKey
	initSecretKey()
}

// initSecretKey 负责初始化或加载 AES 加密密钥
func initSecretKey() {
	// 如果环境变量中已经提供了，直接使用
	if AppConfig.SecretKey != "" {
		if len(AppConfig.SecretKey) != 32 {
			log.Printf("[WARNING] MY_GO_FILE_HUB_SECRET_KEY length is %d, expected 32. Truncating or padding...", len(AppConfig.SecretKey))
			// 简单的补齐或截断，确保长度为 32
			if len(AppConfig.SecretKey) > 32 {
				AppConfig.SecretKey = AppConfig.SecretKey[:32]
			} else {
				AppConfig.SecretKey = fmt.Sprintf("%-32s", AppConfig.SecretKey)
			}
		}
		return
	}

	secretKeyPath := filepath.Join(AppConfig.DataDir, ".secret_key")

	// 尝试从文件读取
	if content, err := os.ReadFile(secretKeyPath); err == nil {
		AppConfig.SecretKey = strings.TrimSpace(string(content))
		if len(AppConfig.SecretKey) == 32 {
			return // 成功读取到 32 字节的密钥
		}
		log.Printf("[WARNING] Loaded secret key from %s has invalid length (%d). A new one will be generated and overwrite it.", secretKeyPath, len(AppConfig.SecretKey))
	}

	// 文件不存在或者读取的内容格式不对，生成一个新的
	key := make([]byte, 16) // 16 bytes = 32 hex chars
	if _, err := rand.Read(key); err != nil {
		panic("Failed to generate random secret key: " + err.Error())
	}
	newSecretKey := hex.EncodeToString(key) // 长度为 32 的字符串

	if err := os.WriteFile(secretKeyPath, []byte(newSecretKey), 0600); err != nil {
		panic(fmt.Sprintf("Failed to write generated secret key to %s: %v", secretKeyPath, err))
	}

	log.Printf("[INFO] Automatically generated a new 32-byte secret key and saved it to %s", secretKeyPath)
	AppConfig.SecretKey = newSecretKey
}
