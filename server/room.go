package server

import (
	"fmt"
	"sync"
)

type room struct {
	name	string
	clients  map[string]*client
	messages []*message
	messageChan map[string]chan *message
	sync.Mutex
}

func NewRoom(roomName string) *room {
	room := &room{
		name: roomName,
		clients: make(map[string]*client),
		messages: make([]*message, 10),
		messageChan: make(map[string](chan *message)),
	}
	return room
}

// 从 ChatS 中取，能取到就返回，不能取到则创建
func (s *ChatS) GetOrCreateByRoom(roomName string) *room {
	room, exist := s.rooms[roomName]
	if exist {
		return room
	}

	s.Lock()
	room = NewRoom(roomName)
	s.rooms[roomName] = room
	s.Unlock()

	fmt.Println("Successful Create the Chat room: ", roomName)
	return s.rooms[roomName]
}

func (r *room) AddClient(client *client) {
	r.Lock()
	r.clients[client.ID] = client
	r.Unlock()
}

func (r *room) AddMessage(message *message) {
	r.messages = append(r.messages, message)

}
