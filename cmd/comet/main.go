package main

import (
	"flag"
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
	flag.Parse()

	comet.Init(conf.Conf.CommetAddress)
}
