package main

import (
	"fmt"
	"github.com/yigger/go-server/server"
	"os"
)

func main() {
	chatServer, err := server.NewChatS()
	if err != nil {
		panic(err)
	}

	fmt.Println("Start the server and listening localhost:5000")
	err = chatServer.Main()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
