package main

import (
	"net"
	"net/rpc"

	"github.com/g-Halo/go-server/logger"
)

var logicRPC *rpc.Client
var authRPC *rpc.Client

func init() {
	client, err := rpc.Dial("tcp", ":7302")
	if err != nil {
		logger.Fatal("7302 无效的地址")
	}
	logicRPC = client

	authRPC, err = rpc.Dial("tcp", ":7301")
	if err != nil {
		logger.Fatal("7301 无效的地址")
	}
}

func main() {
	httpListener, err := net.Listen("tcp", ":4072")
	if err != nil {
		logger.Fatal("error")
	}

	Serve(httpListener, StartServer(), "HTTP")
}
