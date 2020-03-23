package model

import (
	"time"
)

type Room struct {
	UUID      string    `json:"uuid"`
	Name      string    `json:"name"`
	Members   []string  `json:"members"`
	Acceptor  string    `json:"acceptor"`
	Sender    string    `json:"sender"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}

type RoomMessage struct {
	UUID     string     `json:"uuid"`
	Messages []*Message `json:"messages"`
	// MessageChan chan *Message
}

func (*Room) New(uuid, sender, acceptor string) (*Room, *RoomMessage) {
	members := []string{sender, acceptor}
	room := &Room{
		UUID:      uuid,
		Members:   members,
		Type:      "p2p",
		CreatedAt: time.Now(),
		Acceptor:  acceptor,
		Sender:    sender,
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
