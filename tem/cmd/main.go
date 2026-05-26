package main

import (
	"flag"
	"tem/internal/config"
	"tem/internal/srv"
)
import "github.com/baibanzz/jdk/core"

var configDir = flag.String("config", "./cmd/market/conf.yaml", "配置文件路径") // 定义配置文件路径命令行参数

func init() {
	flag.Parse()
}
func main() {
	var config config.Config
	if _, err := core.NewConfigListen(*configDir, &config, nil); err != nil {
		panic(err)
	}

	if srv, err := srv.NewSrv(config); err != nil {
		panic(err)
	} else {

	}
}
