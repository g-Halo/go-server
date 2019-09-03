package server

import (
	"fmt"
	"sync"
)

type room struct {
	name	string
	clients  map[int64]*client
	messages []*message
	messageChan map[int64]chan *message
	sync.Mutex
}

func NewRoom(roomName string) *room {
	room := &room{
		name: roomName,
		clients: make(map[int64]*client),
		messages: make([]*message, 1000, 2000),
		messageChan: make(map[int64](chan *message)),
	}
	return room
}

// 从 ChatS 中取，能取到就返回，不能取到则创建
func (c *ChatS) GetOrCreateByRoom(roomName string) *room {
	room, exist := c.rooms[roomName]
	if exist {
		return room
	}

	c.Lock()
	room = NewRoom(roomName)
	c.rooms[roomName] = room
	c.Unlock()

	fmt.Println("Successful create the Chat room: ", roomName)
	return c.rooms[roomName]
}

func (r *room) AddClient(client *client) {
	r.Lock()
	r.clients[client.ID] = client
	r.Unlock()
}

func (r *room) AddMessage(message *message) {
	r.messages = append(r.messages, message)
}
