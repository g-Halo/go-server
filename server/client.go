package server

import (
	"bufio"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const defaultBufferSize = 16 * 1024

type client struct {
	ID      string

	sync.Mutex

	net.Conn
	wsConn *websocket.Conn

	writeLock sync.RWMutex
	metaLock  sync.RWMutex

	// reading/writing interfaces
	Reader *bufio.Reader
	Writer *bufio.Writer

	lenBuf   [4]byte
	lenSlice []byte

	rooms     []*room
	ClientID string

	// 上下文，一般情况下是 *Chats
	ctx *context
}

func newClient(id string, conn net.Conn, ctx *context) *client {
	var identifier string
	if conn != nil {
		identifier, _, _ = net.SplitHostPort(conn.RemoteAddr().String())
	}

	c := &client{
		ID:  id,
		Conn: conn,
		Reader: bufio.NewReaderSize(conn, defaultBufferSize),
		Writer: bufio.NewWriterSize(conn, defaultBufferSize),
		ClientID: identifier,
		rooms: make([]*room, 100),
		ctx: ctx,
	}
	c.lenSlice = c.lenBuf[:]
	return c
}

func (c *client) addWebSocket(conn *websocket.Conn) {
	c.Lock()
	defer c.Unlock()
	c.wsConn = conn
}

func (c *client) close() {
	c.Conn.Close()
	c.Close()
}

func (c *client) SubRoom(room *room) {
	c.rooms = append(c.rooms, room)
	for {
		select {
		case message := <-room.messageChan[c.ID]:
			if _, err := c.Write([]byte(message.content)); err != nil {
				fmt.Sprintln("SubRoom Get Message Fail, err = %V", err.Error())
			}
		case <-time.After(time.Millisecond):
		}
	}
}

func (c *client) SendMessage(room *room, msg *message) {
	fmt.Printf("Get Message from %d, content %s:", c.ID, msg.content)
	room.AddMessage(msg)

	// 初始化 channel
	if _, ok := room.messageChan[c.ID]; !ok {
		msgChan := make(chan *message)
		room.messageChan[c.ID] = msgChan
	}

	// 给所有在房间内的用户发送消息
	for _, client := range room.clients {
		room.messageChan[client.ID] <- msg
	}
}