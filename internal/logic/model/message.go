package model

import (
	"time"
)

type Message struct {
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
	Body      string `json:"body"`
	Room      Room
	CreatedAt time.Time `json:"created_at"`
	Status    string    `json:"status"`
}

func (Message) Create(sender *User, recipient *User, room Room, body string) *Message {
	message := &Message{
		Sender:    sender.Username,
		Recipient: recipient.Username,
		Room:      room,
		Body:      body,
		CreatedAt: time.Now(),
		Status:    "uncheck",
	}
	return message
}
