package server

import (
	uuid "github.com/satori/go.uuid"
	"github.com/yigger/go-server/logger"
	"github.com/yigger/go-server/model"
	"sync"
	"time"
)

type room struct {
	name	string
	clients  map[string]*client
	messages []*message
	messageChan map[string]chan *message
	sync.Mutex
}

func (s *ChatS) NewRoom(roomName string, members []string) *model.Room {
	//room := &room{
	//	name: roomName,
	//	clients: make(map[string]*client),
	//	messages: make([]*message, 0),
	//	messageChan: make(map[string](chan *message)),
	//}
	uid := uuid.NewV4()
	room := &model.Room{
		UUID:      uid.String(),
		Name:		roomName,
		Type:      "p2p",
		Members:   members,
		CreatedAt: time.Now(),
	}
	return room
}

func (s *ChatS) FindRoomByName(roomName string) *model.Room {
	if room, exist := s.rooms[roomName]; exist {
		return room
	} else {
		return nil
	}
}

// 从 ChatS 中取，能取到就返回，不能取到则创建
func (s *ChatS) GetOrCreateByRoom(roomName string) *model.Room {
	room, exist := s.rooms[roomName]
	if exist {
		return room
	}

	s.Lock()
	room = NewRoom(roomName)
	s.rooms[roomName] = room
	s.Unlock()

	logger.Infof("Successful Create the Chat room: %s", roomName)
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
