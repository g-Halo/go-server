package server

import (
	"github.com/gorilla/websocket"
	"github.com/yigger/go-server/logger"
	"github.com/yigger/go-server/model"
	"net/http"
)

var wsGrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
}

func (s *httpServer) WebSocketConnect(w http.ResponseWriter, req *http.Request, user *model.User) (interface{}, error)  {
	conn, _ := wsGrader.Upgrade(w, req, nil)

	chatServer := s.ctx.chatS
	// 建立 tcp 连接
	tcpListen := chatServer.tcpListener
	tcpConnect, err := tcpListen.Accept()
	if err != nil {
		logger.Error("error in connect to tcp")
		return "fail to connect tcp address", nil
	}

	var client *client
	for _, c := range chatServer.clients {
		if c.ID == user.Username {
			client = c
		}
	}

	if client == nil || client.ID == "" {
		client = newClient(user.Username, tcpConnect, s.ctx)
		chatServer.AddClient(client.ID, client)
	}

	client.addWebSocket(conn)
	go s.wsRead(client)
	go s.wsWrite(client)

	return "OK", nil
}

func (s *httpServer) wsRead(client *client) {
	conn := client.wsConn
	for {
		_, json, _ := conn.ReadMessage()
		if json == nil {
			return
		}

	}
}

func (s *httpServer) wsWrite(client *client) {
	conn := client.wsConn
	for {
		for _, room := range client.rooms {
			msg := <- room.messageChan[client.ID]

			w, err := conn.NextWriter(websocket.TextMessage)
			if err != nil {
				logger.Error(err)
			}
			w.Write([]byte(msg.content))
		}
	}
}