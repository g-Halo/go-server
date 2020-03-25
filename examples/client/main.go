package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	LoginUrl       = "http://127.0.0.1:7834/v1/login"
	PushMessageUrl = "http://127.0.0.1:7834/v1/room/push"
)

var Token = ""

// var WebSocketUrl = flag.String("addr", "", "http service address")
var currentUser = map[string]string{
	"username": "test1",
	"password": "123",
}

type RespCommon struct {
	Data   map[string]interface{} `json:"data"`
	Status int                    `json:"status"`
}

type TextMessageReq struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

func main() {
	// 登录
	log.Print("Start Login...\n\n")
	login()
	log.Printf("token: %s \n\n", Token)
	// ws 连接
	senderWsConnect()
	acceptorWsConnect()

	// 获取当前用户的房间信息
	// 获取客户端的输入，往对方发送消息
	go func() {
		index := 0
		for {
			log.Print("Start Push Message...")
			pushTextMessage("test2", fmt.Sprintf("hi, %d", index))
			index += 1
			time.Sleep(time.Second * 1)
		}
	}()

	go func() {
		index := 20000
		for {
			log.Print("Start Push Message...")
			pushTextMessage("test2", fmt.Sprintf("hi, %d", index))
			index += 1
			time.Sleep(time.Second * 1)
		}
	}()

	select {}
}

// example: https://github.com/gorilla/websocket/blob/master/examples/echo/client.go
func senderWsConnect() {
	c, _, err := websocket.DefaultDialer.Dial("ws://127.0.0.1:7771/v1/ws?username=test1&token=test", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Fatal("die")
				return
			}
			log.Printf("test1 receive: %s", message)
		}
	}()
}

func acceptorWsConnect() {
	c, _, err := websocket.DefaultDialer.Dial("ws://127.0.0.1:7771/v1/ws?username=test2&token=test", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Fatal("die")
				return
			}
			log.Printf("test2 receive: %s", message)
		}
	}()
}

func login() {
	jsonStr, _ := json.Marshal(currentUser)
	req, err := http.NewRequest("POST", LoginUrl, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	response := &RespCommon{}
	err = json.NewDecoder(resp.Body).Decode(response)
	if err != nil {
		panic(err)
	}

	if response.Status == 200 {
		Token = response.Data["token"].(string)
	} else {
		panic("token is nil, login fail")
	}
}

func pushTextMessage(acceptorUsername string, message string) {
	reqParams := &TextMessageReq{Username: acceptorUsername, Message: message}
	jsonStr, _ := json.Marshal(reqParams)
	req, err := http.NewRequest("POST", PushMessageUrl, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("token", Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// response := &RespCommon{}
	// err = json.NewDecoder(resp.Body).Decode(response)
	str, _ := ioutil.ReadAll(resp.Body)
	log.Print(string(str))
}
