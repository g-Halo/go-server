package main

import (
	"github.com/g-Halo/go-server/logger"
	"net"
	"net/rpc"
)

type Logic struct {
	userLogic
}

func StartRpc() {
	logic := new(Logic)
	rpc.Register(logic)

	tcpAddr, err := net.ResolveTCPAddr("tcp",":7302")
	if err != nil {
		logger.Fatal("Bind Rpc Address Err", err.Error())
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		logger.Fatal("Listen Tcp Err")
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Error("accept error")
			continue
		}
		rpc.ServeConn(conn)
	}
}