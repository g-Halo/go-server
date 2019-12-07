package http_api

import (
	"github.com/g-Halo/go-server/rpc/logic"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/g-Halo/go-server/model"
	"github.com/g-Halo/go-server/rpc/instance"

	"github.com/g-Halo/go-server/logger"
	"github.com/g-Halo/go-server/storage"
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
	conn   *websocket.Conn
	writer io.WriteCloser
	user   *model.User
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

	minHearbeatSec = 30

	Second = int64(time.Second)
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

	wsWrite, err := conn.NextWriter(websocket.TextMessage)
	if currentUser == nil {
		errorMsg := []byte("-p\r\n")
		wsWrite.Write(errorMsg)
	}

	// 注册 client，包括 user
	client := &Client{user: currentUser, conn: conn, writer: wsWrite}

	// 注册 heartbeat
	heartBeat := r.URL.Query().Get("heartbeat")
	heartBeatTime, err := strconv.Atoi(heartBeat)
	if err != nil || heartBeatTime < minHearbeatSec {
		wsWrite.Write([]byte("-p\r\n"))
	}

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

	begin := time.Now().UnixNano()
	end := begin + Second
	for {
		if end-begin >= Second {
			// 超过 1s 重置定时器?
		}
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		if string(message) == "h" {
			// 回应给 客户端
			c.writer.Write([]byte("+h\r\n"))
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		user := storage.GetUser(c.user.Username)
		if user == nil {
			continue
		}

		for _, room := range user.Rooms {
			rChan, _ := logic.RoomChannels.Get(room.UUID)

			msg := rChan.GetMsg(c.user.Username)
			if msg == nil {
				continue
			} else {
				logger.Info("get message")
				//logger.Info("Get the message: %s", string(msg.Body))
			}

		}
	}
}
