package server

import (
	"github.com/gorilla/websocket"
	"github.com/yigger/go-server/model"
	"net/http"
)

var wsGrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
}

func (s *httpServer) WebSocketConnect(w http.ResponseWriter, req *http.Request, user *model.User) {
	conn, _ := wsGrader.Upgrade(w, req, nil)

	chatServer := s.ctx.chatS
	// 建立 tcp 连接
	tcpListen := chatServer.tcpListener
	tcpConnect, err := tcpListen.Accept()
	if err != nil {
		panic("error in connect to tcp")
	}

	// 把当前连接的客户端添加进聊天服务
	client := newClient(user.Username, tcpConnect, s.ctx)
	chatServer.AddClient(client.ID, client)
	client.addWebSocket(conn)

	go s.wsRead(client)
}

func (s *httpServer) wsRead(client *client) {
	conn := client.wsConn
	for {
		_, json, _ := conn.ReadMessage()
		if json == nil {
			return
		}

		// 往 room 写数据
		//fmt.Println(string(message))
	}
}