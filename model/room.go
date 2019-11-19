package model

import (
	"time"
)

type Room struct {
	UUID      string     `json:"uuid"`
	Name      string     `json:"salt"`
	Members   []string   `json:"members"`
	Type      string     `json:"type"`
	CreatedAt time.Time  `json:"created_at"`
	Messages  []*Message `json:"messages"`

	MessageChan chan *Message
}

func (Room) New(uuid string, members []string) *Room {
	return &Room{
		UUID:        uuid,
		Members:     members,
		Type:        "p2p",
		CreatedAt:   time.Now(),
		MessageChan: make(chan *Message, 256),
	}
}

func (r *Room) AddMessage(message *Message) {
	r.Messages = append(r.Messages, message)
	// if message == nil {
	// 	logger.Error("message is null")
	// 	return
	// }

	// filter := bson.D{{"uuid", r.UUID}}
	// update := bson.D{
	// 	{"$push", bson.D{
	// 		{"messages", message},
	// 	}},
	// }
	// collection := Collection("rooms")
	// _, err := collection.UpdateOne(context.TODO(), filter, update)
	// if err != nil {
	// 	logger.Error(err)
	// }
}
