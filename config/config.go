package config

import (
	"net"
	"strconv"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	ServerPort string `env:"MY_GO_FILE_HUB_SERVER_PORT" envDefault:"3939"`
	Listen     string `env:"MY_GO_FILE_HUB_LISTEN" envDefault:"localhost"`
	DataDir    string `env:"MY_GO_FILE_HUB_DATA_DIR" envDefault:"./data"`
	LanOnly    string `env:"MY_GO_FILE_HUB_LAN_ONLYDB_USER" envDefault:"false"`
}

var AppConfig Config

func init() {
	if err := env.Parse(&AppConfig); err != nil {
		panic(err)
	}
	if AppConfig.LanOnly != "true" && AppConfig.LanOnly != "false" {
		panic("Invalid value for MY_GO_FILE_HUB_LAN_ONLY, must be 'true' or 'false'")
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
