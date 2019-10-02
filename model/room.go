package model

import (
	"context"
	"fmt"
	"github.com/yigger/go-server/logger"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type Room struct {
	UUID 				string 		`json:"uuid"`
	Name	 			string 		`json:"salt"`
	Members				[]string	`json:"members"`
	Type				string      `json:"type"`
	CreatedAt 			time.Time 	`json:"created_at"`
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
	fmt.Println(members)

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