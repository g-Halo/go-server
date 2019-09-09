package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func receive(conn net.Conn) {
	for {
		buf := make([]byte, 1024)
		len, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading", err.Error())
			return //终止程序
		}
		switch string(buf[:len]) {
		case "__heartbeat__":
			fmt.Println("heartbeat ++")
		}
		fmt.Println("Received data: ", string(buf[:len]))
	}
}

func main() {
	conn, err := net.Dial("tcp", "localhost:5000")
	if err != nil {
		panic(err)
	}

	// 发送特定的协议内容开启通讯
	_, err = conn.Write([]byte("  CHAT"))

	// 接受服务端发送来的消息
	go receive(conn)

	for {
		fmt.Print("please input：")
		inputReader := bufio.NewReader(os.Stdin)

		clientName, _ := inputReader.ReadString('\n')
		message := fmt.Sprintf("SEND test %s\n", clientName)
		_, err = conn.Write([]byte(message))
	}
}
