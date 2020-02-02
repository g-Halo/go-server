package ws_conn

import (
	"context"
	"encoding/json"
	"github.com/g-Halo/go-server/internal/logic/model"
	"github.com/g-Halo/go-server/pkg/logger"
	"github.com/g-Halo/go-server/pkg/pb"
	"github.com/g-Halo/go-server/pkg/storage"
	"github.com/gorilla/websocket"
	"io"
	"net/http"
	"sync"
	"time"
)

const (
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

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

	isClosed bool
	conn     *websocket.Conn
	writer   io.WriteCloser
	user     *model.User
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
		w.Write([]byte("-p\r\n"))
		logger.Error("user not found")
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

	// 存储当前连接
	store(currentUser.Username, client)
	//client.Run()
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
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Close()
	}()

	//for {
	//	c.conn.SetReadLimit(maxMessageSize)
	//	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	//	//_, _, err := c.conn.ReadMessage()
	//	//if err != nil {
	//	//	logger.Error(err)
	//	//}
	//}
}


func (c *Client) Close() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.isClosed = true
	c.conn.Close()
}