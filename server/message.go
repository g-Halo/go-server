package server

import "time"

type message struct {
	sender *client
	accepter *client
	content string

	createdAt time.Time
}

func NewMessage(sender *client, content string) *message {
	message := &message{
		content: content,
		sender: sender,
		createdAt: time.Now(),
	}
	return message
}