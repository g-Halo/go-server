package model

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"time"
)

type Room struct {
	UUID 				string 		`json:"uuid"`
	Name	 			string 		`json:"salt"`
	Members				[]string	`json:"members"`
	Type				string      `json:"type"`
	CreatedAt 			time.Time 	`json:"created_at"`
}

func (r *Room) Create(client *mongo.Client) {
	collection := client.Database("chat").Collection("rooms")
	_, err := collection.InsertOne(context.TODO(), r)
	if err != nil {
		fmt.Println("err", err)
		log.Fatal(err)
	}
}

func (r *Room) AddMembers(client *mongo.Client, members []string) {
	filter := bson.D{{"uuid", r.UUID}}
	fmt.Println(members)

	update := bson.D{
		{"$set", bson.D{
			{"members", members},
		}},
	}

	collection := client.Database("chat").Collection("rooms")
	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Println(err)
	}
}