package model

import "time"

type Message struct {
	SenderId		int			`json:"send_id"`
	Sender 			string 		`json:"sender"`
	Recipient	 	string 		`json:"recipient"`
	RecipientId		int			`json:"recipient_id"`
	Body 			string		`json:"body"`
	CreatedAt 		time.Time 	`json:"created_at"`
	Status			string		`json:"status"`
}

func (m Message) Create(sender string, recipient string, body string) *Message {
	message := &Message{
		Sender: sender,
		Recipient: recipient,
		Body: body,
		CreatedAt: time.Now(),
		Status: "uncheck",
	}

	return message
}