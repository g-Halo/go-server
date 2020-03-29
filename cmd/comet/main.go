package main

import (
	"github.com/g-Halo/go-server/conf"
	"github.com/g-Halo/go-server/internal/comet"
	"github.com/g-Halo/go-server/pkg/logger"
)

func init() {
	// 初始化日志文件
	logger.InitLogger("./comet.log", "debug")
	// 初始化配置文件
	conf.LoadConf()
}

func main() {
	comet.Init(conf.Conf.CommetAddress)
}
