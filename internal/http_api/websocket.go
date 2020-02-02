package http_api

import (
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/g-Halo/go-server/internal/logic/model"
	"github.com/g-Halo/go-server/pkg/logger"
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
	// Time allowed to read the next pong message from the peer.
	pongWait = 120 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

func serveWs(w http.ResponseWriter, r *http.Request) {
	// token 校验，获取 params token，校验是否存在用户
	username := r.URL.Query().Get("username")
	if username == "" {
		w.Write([]byte("not sign in"))
		return
	}

	currentUser := getUser(username)
	if currentUser == nil {
		w.Write([]byte("-p\r\n"))
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Fatal(err)
	}
	wsWrite, err := conn.NextWriter(websocket.TextMessage)
	client := &Client{
		user:   currentUser,
		conn:   conn,
		writer: wsWrite,
	}


	go client.Done()
}

func (c *Client) Done() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Close()
	}()

	for {
		c.conn.SetReadLimit(maxMessageSize)
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			logger.Error(err)
		}

		//for _, room := range user.Rooms {
		//	msg, _ := rpc_client.LogicClient.KeepGetMessage(context.Background(), &pb.KeepGetMessageReq{Username: user.Username, Uuid: room.UUID})
		//	if msg == nil {
		//		continue
		//	} else {
		//		message, _ := json.Marshal(msg)
		//		w, err := c.conn.NextWriter(websocket.TextMessage)
		//		if err != nil {
		//			logger.Error(err)
		//			c.Close()
		//		}
		//		logger.Debugf("test: %s", message)
		//		w.Write(message)
		//	}
		//}
	}
}

func (c *Client) Output() {

}

//func messageResponse(message *model.Message) []byte {
//	sender := getUser(message.Sender)
//	receiver := getUser(message.Recipient)
//
//	res, err := json.Marshal(struct {
//		Sender   map[string]interface{} `json:"sender"`
//		Accepter map[string]interface{} `json:"accepter"`
//		Room     map[string]string      `json:"room"`
//		Message  map[string]interface{} `json:"message"`
//	}{
//		Sender:   sender.ToJson(),
//		Accepter: receiver.ToJson(),
//		Room: map[string]string{
//			"uuid": message.Room.UUID,
//		},
//		Message: map[string]interface{}{
//			"body":       message.Body,
//			"recipient":  receiver.ToJson(),
//			"sender":     sender.ToJson(),
//			"created_at": message.CreatedAt,
//			"status":     message.Status,
//		},
//	})
//
//	if err != nil {
//		logger.Error("message response json error")
//	}
//
//	return res
//}

func (c *Client) Write(message string) {
	w, err := c.conn.NextWriter(websocket.TextMessage)
	if err != nil {
		logger.Error(err)
		c.Close()
	}
	w.Write([]byte(message))
}

func (c *Client) Close() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.isClosed = true
	c.conn.Close()
}
