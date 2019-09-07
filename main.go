package main

import (
	"fmt"
	"github.com/yigger/go-server/conf"
	"github.com/yigger/go-server/server"
	"os"
)

func main() {
	config := conf.LoadConf()
	chatServer, err := server.NewChatS(config)
	if err != nil {
		panic(err)
	}

	err = chatServer.Main()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
