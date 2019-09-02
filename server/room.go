package server

import "sync"

type room struct {
	name	string
	clients  map[int64]*client

	sync.Mutex
}

func NewRoom(roomName string) *room {
	room := &room{
		name: roomName,
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

	return c.rooms[roomName]
}

func (r *room) AddClient(client *client) {
	r.Lock()
	r.clients[client.ID] = client
	r.Unlock()
}