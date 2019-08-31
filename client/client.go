package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/yigger/go-server/model"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net"
	"os"
	"strings"
)

func getCommand() string {
	// 从键盘里面输入
	inputReader := bufio.NewReader(os.Stdin)
	// 不知道干什么
	clientName, _ := inputReader.ReadString('\n')
	// 过滤 \r\n 仅 Windows 平台需要
	return strings.Trim(clientName, "\r\n")
}

func init() {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	model.DB = client
}

func main() {
	conn, err := net.Dial("tcp", "localhost:5000")
	if err != nil {
		panic(err)
	}
	// initialize tcp protocol
	_, err = conn.Write([]byte("  CHAT"))
	for {
		fmt.Print("please input：")
		inputReader := bufio.NewReader(os.Stdin)
		clientName, _ := inputReader.ReadString('\n')
		_, err = conn.Write([]byte(clientName))
	}

	var User model.User
	for {
		//fmt.Println("---------------------")
		fmt.Println("1. Login")
		fmt.Println("2. Sign")
		str := getCommand()
		switch str {
		case "1":
			//fmt.Println("1. back")
			fmt.Print("请输入用户名：")
			username := getCommand()
			//fmt.Println("1. back")
			fmt.Print("请输入密码：")
			password := getCommand()
			if User.Login(username, password) {
				fmt.Println("登录成功")
				goto LOGINED
			}
		case "2":
			fmt.Print("请输入用户昵称：")
			nickname := getCommand()
			fmt.Print("请输入用户名：")
			username := getCommand()
			//fmt.Println("1. back")
			fmt.Print("请输入密码：")
			password := getCommand()
			User.SignUp(nickname, username, password)
		}
	}



LOGINED:
	fmt.Println("logined")
	//for {
	//	fmt.Println("What to send to the server?")
	//	input, _ := inputReader.ReadString('\n')
	//	trimmedInput := strings.Trim(input, "\r\n")
	//	if trimmedInput == "Q" {
	//		return
	//	}
	//	_, err = conn.Write([]byte(trimmedClient + "say:" + trimmedInput))
	//}
}
