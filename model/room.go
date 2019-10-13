package model

import (
	"context"
	"github.com/g-Halo/go-server/logger"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type Room struct {
	UUID 				string 		`json:"uuid"`
	Name	 			string 		`json:"salt"`
	Members				[]string	`json:"members"`
	Type				string      `json:"type"`
	CreatedAt 			time.Time 	`json:"created_at"`
	Messages			[]*Message	`json:"messages"`
}

func (Room) New(uuid string, members []string) *Room {
	return &Room{
		UUID: uuid,
		Members: members,
		Type: "p2p",
		CreatedAt: time.Now(),
	}
}

func (r *Room) Create() {
	collection := Collection("rooms")
	_, err := collection.InsertOne(context.TODO(), r)
	if err != nil {
		logger.Error(err)
	}
}

func (r *Room) AddMembers(members []string) {
	filter := bson.D{{"uuid", r.UUID}}

	update := bson.D{
		{"$set", bson.D{
			{"members", members},
		}},
	}

	collection := Collection("rooms")
	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		logger.Error(err)
	}
}

func (r *Room) AddMessage(message *Message) {

	if message == nil {
		logger.Error("message is null")
		return
	}

	filter := bson.D{{"uuid", r.UUID}}
	update := bson.D{
		{"$push", bson.D{
			{"messages", message},
		}},
	}
	collection := Collection("rooms")
	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		logger.Error(err)
	}
}