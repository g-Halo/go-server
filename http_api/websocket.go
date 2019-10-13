package http_api

import (
	"github.com/g-Halo/go-server/model"
	"github.com/g-Halo/go-server/server"
	"github.com/gorilla/websocket"
	"net/http"
)

var wsGrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
}

func (s *httpServer) WebSocketConnect(w http.ResponseWriter, req *http.Request, user *model.User) (interface{}, error)  {
	//conn, _ := wsGrader.Upgrade(w, req, nil)
	//
	//chatServer := s.ctx.chatS
	//// 建立 tcp 连接
	//tcpListen := chatServer.TcpListener
	//tcpConnect, err := tcpListen.Accept()
	//if err != nil {
	//	logger.Error("error in connect to tcp")
	//	return "fail to connect tcp address", nil
	//}
	//
	//var client *server.client
	//for _, c := range chatServer.Clients {
	//	if c.ID == user.Username {
	//		client = c
	//	}
	//}
	//
	//if client == nil || server.ID == "" {
	//	client = server.newClient(user.Username, tcpConnect, s.ctx)
	//	chatServer.AddClient(server.ID, client)
	//}
	//
	//server.addWebSocket(conn)
	//go s.wsRead(client)
	//go s.wsWrite(client)

	return "OK", nil
}

func (s *httpServer) wsRead(client *server.Client) {
	//conn := server.wsConn
	//for {
	//	_, json, _ := conn.ReadMessage()
	//	if json == nil {
	//		return
	//	}
	//
	//}
}

func (s *httpServer) wsWrite(client *server.Client) {
	//conn := client.wsConn
	//for {
	//	for _, room := range client.rooms {
	//		msg := <- room.messageChan[client.ID]
	//
	//		w, err := conn.NextWriter(websocket.TextMessage)
	//		if err != nil {
	//			logger.Error(err)
	//		}
	//		w.Write([]byte(msg.content))
	//	}
	//}
}