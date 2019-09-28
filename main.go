package main

import (
	"github.com/yigger/go-server/conf"
	"github.com/yigger/go-server/logger"
	"github.com/yigger/go-server/server"
	"log"
	"os"
)

func init() {
	logger.InitLogger("./application.log", "debug")
}

func main() {
	config := conf.LoadConf()
	chatServer, err := server.NewChatS(config)
	if err != nil {
		panic(err)
	}

	err = chatServer.Main()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
