package ws_conn

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/g-Halo/go-server/internal/logic/model"
	"github.com/g-Halo/go-server/pkg/logger"
	"github.com/g-Halo/go-server/pkg/pb"
	"github.com/g-Halo/go-server/pkg/rpc_client"
	"github.com/g-Halo/go-server/pkg/storage"
	"github.com/gorilla/websocket"
)

const (
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 4 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
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

	conn   *websocket.Conn
	writer io.WriteCloser
	user   *model.User
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	// token 校验，获取 params token，校验是否存在用户
	username := r.URL.Query().Get("username")
	if username == "" {
		w.Write([]byte("not sign in"))
		logger.Error("not sign in")
	}

	currentUser := storage.GetUser(username)
	if currentUser == nil {
		logger.Error("user not found")
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Fatal(err)
	}
	client := &Client{
		user: currentUser,
		conn: conn,
	}

	client.write(websocket.TextMessage, []byte("success connect to websocket!!!"))

	// pings headtbeat
	go client.ping()

	// 存储当前连接
	store(currentUser.Username, client)
	client.Run()
}

// write writes a message with the given message type and payload.
func (c *Client) write(mt int, payload []byte) error {
	// 写入超时
	c.conn.SetWriteDeadline(time.Now().Add(writeWait))
	return c.conn.WriteMessage(mt, payload)
}

func (client *Client) ping() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case <-ticker.C:
			if err := client.write(websocket.PingMessage, []byte{}); err != nil {
				logger.Error("对方连接已经断开")
				// TODO: 尝试重新连接
				return
			} else {
				logger.Infof("heart beat to %s", client.user.Username)
			}
		}
	}
}

func (c *Client) DispatchMessage(ctx context.Context, req *pb.DispatchReq) {
	v, err := json.Marshal(req)
	if err != nil {
		logger.Error(err)
		return
	}

	if c.conn == nil {
		logger.Error("connect is null")
		return
	}
	w, err := c.conn.NextWriter(websocket.TextMessage)
	if w == nil {
		logger.Error("writer is null")
		return
	}

	if err != nil {
		logger.Error(err)
		c.Close()
	}
	c.conn.WriteMessage(websocket.TextMessage, v)
}

func (c *Client) Run() {
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}
		logger.Info("接收到来自客户端的消息:", message)
	}
}

func (c *Client) Close() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	rpc_client.LogicClient.UserOffline(context.Background(), &pb.UserOnlineReq{
		Username: c.user.Username,
	})
	c.conn.Close()
}
