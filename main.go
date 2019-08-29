package main

import (
	"fmt"
	"github.com/yigger/go-server/server"
	"gopkg.in/mgo.v2"
	"net"
	"sync"

	//"gopkg.in/mgo.v2/bson"
)

func init() {
	fmt.Println("init")
	session, err:=mgo.Dial("mongodb://localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()
}

func main() {
	// 监听端口号
	listener, err := net.Listen("tcp", "localhost:5000")
	if err != nil {
		panic(err)
	}

	// 定义全局的 Once
	var exitCh = make(chan error)
	var once sync.Once
	exitFunc := func(err error){
		once.Do(func() {
			exitCh <- err
		})
	}


	chatServer := server.NewChatS()
	exitFunc(chatServer.Main(listener))
	err = <- exitCh
	panic(err)
}
