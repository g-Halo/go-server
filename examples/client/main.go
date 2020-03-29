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

var senderConn *websocket.Conn

type client struct {
	username string
	wsUrl    string
	conn     *websocket.Conn
}

func main() {
	// 登录
	log.Print("Start Login...\n\n")
	login()
	log.Printf("token: %s \n\n", Token)

	senderClient := &client{
		username: "test1",
		wsUrl:    "ws://127.0.0.1:7771/v1/ws?username=test1&token=test",
	}
	accepClient := &client{
		username: "test2",
		wsUrl:    "ws://127.0.0.1:7771/v1/ws?username=test2&token=test",
	}
	// ws 连接
	senderClient.connetWs()
	accepClient.connetWs()

	// go senderClient.Ping()
	go senderClient.receiveMessage()
	// go accepClient.Ping()
	go accepClient.receiveMessage()

	// 获取客户端的输入，往对方发送消息
	for i := 1; i < 2; i += 1 {
		go func(i int) {
			index := 0
			fmt.Printf("协程 %d", i)
			for {
				log.Print("Start Push Message...")
				pushTextMessage("test2", fmt.Sprintf("hi, 协程 %d, %d", i, index))
				index += 1
				time.Sleep(time.Millisecond * 100)
			}
		}(i)
	}

	select {}
}

// example: https://github.com/gorilla/websocket/blob/master/examples/echo/client.go
func (client *client) connetWs() {
	url := client.wsUrl
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Print("connect fail..")
		for {
			t := time.NewTicker(3 * time.Second)
			select {
			case <-t.C:
				log.Println("retry to connect...")
				c, _, err := websocket.DefaultDialer.Dial(url, nil)
				if err == nil {
					t.Stop()
					client.conn = c
					return
				}
			}
		}
	} else {
		client.conn = c
	}
}

func (c *client) receiveMessage() {
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			c.connetWs()
		}
		log.Printf("%s receive: %s", c.username, message)
	}
}

// func (c *client) Ping() {
// 	for {
// 		err := c.conn.WriteMessage(websocket.PongMessage, []byte{})
// 		if err != nil {
// 			// 服务端连接已断开

// 		}
// 		time.Sleep(time.Second * 3)
// 	}
// }

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

	client := &http.Client{
		Timeout: time.Second * 3,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Print(err)
		return
	}
	defer resp.Body.Close()

	// response := &RespCommon{}
	// err = json.NewDecoder(resp.Body).Decode(response)
	str, _ := ioutil.ReadAll(resp.Body)
	log.Printf("发送 %s，status: %s", message, string(str))
}
