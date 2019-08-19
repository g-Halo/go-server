package main

import (
	"fmt"
	"net"
)

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

		go doServer(conn)
	}
}

func doServer(conn net.Conn) {
	for {
		buf := make([]byte, 512)
		len, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading")
			return
		}

		fmt.Println("Receive data is: %v", string(buf[:len]))
	}
}