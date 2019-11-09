package main

import (
	"github.com/g-Halo/go-server/logger"
	"net/rpc"
	"net"
)

func StartRpc() {
	token := new(Token)
	rpc.Register(token)
	tcpAddr, err := net.ResolveTCPAddr("tcp",":7301")
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