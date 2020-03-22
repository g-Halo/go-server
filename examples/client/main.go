package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
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

func main() {
	// 登录
	log.Print("Start Login...\n\n")
	login()
	log.Println("token:", Token)
	// ws 连接

	// 获取当前用户的房间信息
	// 获取客户端的输入，往对方发送消息
	log.Print("\n\nStart Push Message...")
	pushTextMessage("test2", "hi")
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
	Token = response.Data["token"].(string)
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