package main

import (
	"fmt"
	"net"
	"gopkg.in/mgo.v2"
	"github.com/yigger/go-server/server"
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
	fmt.Println("start a server")
	listener, err := net.Listen("tcp", "localhost:5000")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("无效的请求链接")
		}

		srv := &server.Server{conn}
		go doServer(srv)
	}
}

func doServer(server *server.Server) {
	for {
		server.Main()
	}
}