package http_api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/g-Halo/go-server/model"
	"github.com/g-Halo/go-server/rpc/instance"
	"github.com/g-Halo/go-server/rpc/logic"
	"github.com/g-Halo/go-server/storage"

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
	mutex sync.Mutex

	isClosed bool
	conn     *websocket.Conn
	writer   io.WriteCloser
	user     *model.User
}

type WsParams struct {
	Type     string `json:"type"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

type WsResponse struct {
	Type    string                 `json:"type"`
	User    map[string]interface{} `json:"user"`
	Message map[string]interface{} `json:"message"`
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

	client := &Client{
		user:   currentUser,
		conn:   conn,
		writer: wsWrite,
	}

	// 注册 heartbeat
	// heartBeat := r.URL.Query().Get("heartbeat")
	// heartBeatTime, err := strconv.Atoi(heartBeat)
	// if err != nil || heartBeatTime < minHearbeatSec {
	// 	wsWrite.Write([]byte("-p\r\n"))
	// }

	go client.writePump()
	// go client.readPump()
}

func (c *Client) Write(message string) {
	w, err := c.conn.NextWriter(websocket.TextMessage)
	if err != nil {
		logger.Error(err)
		c.Close()
	}
	w.Write([]byte(message))
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
		c.writer.Write([]byte("+h\r\n"))
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
		c.Close()
	}()

	for {
		user := storage.GetUser(c.user.Username)
		if user == nil {
			continue
		}

		for _, room := range user.Rooms {
			rChan, _ := logic.RoomChannels.Get(room.UUID)

			msg := rChan.GetMsg(c.user.Username)
			logger.Info(msg)
			if msg == nil {
				continue
			} else {
				res := messageResponse(msg)
				logger.Info("success send")
				c.Write(string(res))
			}

		}

	}
}

func messageResponse(message *model.Message) []byte {
	logicClient := instance.LogicRPC()
	var sender model.User
	var accepter model.User
	logicClient.Call("Logic.FindByUsername", &message.Sender, &sender)
	logicClient.Call("Logic.FindByUsername", &message.Recipient, &accepter)

	res, err := json.Marshal(struct {
		Sender   map[string]interface{} `json:"sender"`
		Accepter map[string]interface{} `json:"accepter"`
		Room     map[string]string      `json:"room"`
		Message  map[string]interface{} `json:"message"`
	}{
		Sender:   sender.ToJson(),
		Accepter: accepter.ToJson(),
		Room: map[string]string{
			"uuid": message.Room.UUID,
		},
		Message: map[string]interface{}{
			"body":       message.Body,
			"recipient":  accepter.ToJson(),
			"sender":     sender.ToJson(),
			"created_at": message.CreatedAt,
			"status":     message.Status,
		},
	})

	if err != nil {
		logger.Error("message response json error")
	}

	return res
}

func (c *Client) Close() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.isClosed = true
	c.conn.Close()
}
