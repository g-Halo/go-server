package main

import (
	"github.com/yigger/go-server/conf"
	"github.com/yigger/go-server/logger"
	"github.com/yigger/go-server/server"
	"runtime"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	logger.InitLogger("./application.log", "debug")
}

func main() {
	config := conf.LoadConf()
	chatServer, err := server.NewChatS(config)
	if err != nil {
		logger.Fatal(err)
	}

	err = chatServer.Main()
	if err != nil {
		logger.Fatal(err)
	}
}
