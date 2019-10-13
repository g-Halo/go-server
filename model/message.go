package model

import (
	"time"
)

type Message struct {
	SenderId		int			`json:"send_id"`
	Sender 			string 		`json:"sender"`
	Recipient	 	string 		`json:"recipient"`
	RecipientId		int			`json:"recipient_id"`
	Body 			string		`json:"body"`
	CreatedAt 		time.Time 	`json:"created_at"`
	Status			string		`json:"status"`
}

func (Message) Create(sender *User, recipient *User, body string) *Message {
	message := &Message{
		Sender: sender.Username,
		Recipient: recipient.Username,
		Body: body,
		CreatedAt: time.Now(),
		Status: "uncheck",
	}
	return message
	//collection := Collection("messages")
	//_, err := collection.InsertOne(context.TODO(), message)
	//if err != nil {
	//	logger.Error(err)
	//	return nil
	//}
	//
	//return message
}