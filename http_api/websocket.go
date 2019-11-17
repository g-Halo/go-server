package http_api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/g-Halo/go-server/model"
	"github.com/g-Halo/go-server/rpc/instance"

	"github.com/g-Halo/go-server/logger"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	conn *websocket.Conn
	send chan []byte
	user *model.User
}

type WsParams struct {
	Type     string `json:"type"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

func serveWs(w http.ResponseWriter, r *http.Request) {
	// token 校验，获取 params token，校验是否存在用户
	username := "test1"
	logicClient := instance.LogicRPC()
	var currentUser *model.User
	logicClient.Call("Logic.FindByUsername", &username, &currentUser)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Fatal(err)
	}

	// 注册 client，包括 user
	client := &Client{user: currentUser, conn: conn, send: make(chan []byte, 256)}

	go client.writePump()
	go client.readPump()
}

func (c *Client) readPump() {
	defer func() {
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	var params *WsParams
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		if err := json.Unmarshal(message, &params); err != nil {
			logger.Error(err)
			continue
		}

		// 查找发送消息的目标对象
		logicClient := instance.LogicRPC()
		var user *model.User
		logicClient.Call("Logic.FindByUsername", &params.Username, &user)
		if user == nil {
			logger.Error("username not found")
			continue
		}

		// 创建房间, 然后双方同时订阅房间
		var room *model.Room
		if err := logicClient.Call("Logic.FindOrCreate", []string{user.Username, c.user.Username}, &room); err != nil {
			logger.Error(err)
		}
		c.user.Rooms = append(c.user.Rooms, room)
		user.Rooms = append(user.Rooms, room)

		// 分发器
		go room.Dispatch()

		// 往房间发送消息
		var Message model.Message
		msg := Message.Create(c.user, user, string(message))
		// room.MessageChan <- msg
		logger.Info(msg)
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		for _, room := range c.user.Rooms {
			logger.Info("room")
			logger.Info(<-room.MessageChan)
			logger.Info("message chan")
			select {
			case msg := <-room.MessageChan:
				logger.Info("read")
				logger.Info(msg)
			}
		}
	}
}
