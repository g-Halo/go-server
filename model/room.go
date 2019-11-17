package model

import (
	"time"

	"github.com/g-Halo/go-server/logger"
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

// 分发器
// 1. 从 Chan 接收到消息
// 2. 分发到每个已经连接用户的 Chan
func (r *Room) Dispatch() {
	for {
		logger.Info("dispatch")
		select {
		case msg := <-r.MessageChan:
			logger.Info("分发")
			logger.Info(msg)
			r.AddMessage(msg)
			// 分发
		}
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
