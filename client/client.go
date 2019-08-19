package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:5000")
	if err != nil {
		panic(err)
	}

	// 从键盘里面输入
	inputReader := bufio.NewReader(os.Stdin)
	// 不知道干什么
	clientName, _ := inputReader.ReadString('\n')
	// 过滤 \r\n 仅 Windows 平台需要
	trimmedClient := strings.Trim(clientName, "\r\n")

	for {
		fmt.Println("What to send to the server?")
		input, _ := inputReader.ReadString('\n')
		trimmedInput := strings.Trim(input, "\r\n")
		if trimmedInput == "Q" {
			return
		}
		_, err = conn.Write([]byte(trimmedClient + "say:" + trimmedInput))
	}

}
