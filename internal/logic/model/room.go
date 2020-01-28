package model

import (
	"time"
)

type Room struct {
	UUID      string    `json:"uuid"`
	Name      string    `json:"name"`
	Members   []string  `json:"members"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}

type RoomMessage struct {
	UUID     string     `json:"uuid"`
	Messages []*Message `json:"messages"`
	// MessageChan chan *Message
}

func (*Room) New(uuid string, members []string) (*Room, *RoomMessage) {
	room := &Room{
		UUID:      uuid,
		Members:   members,
		Type:      "p2p",
		CreatedAt: time.Now(),
		// MessageChan: make(chan *Message, 256),
		// Messages: make([]*Message, 32), // RPC 不能回传指针
	}
	rmsg := &RoomMessage{
		UUID:     uuid,
		Messages: make([]*Message, 0),
	}
	return room, rmsg
}

func (rmsg *RoomMessage) AddMessage(message *Message) {
	rmsg.Messages = append(rmsg.Messages, message)
}
