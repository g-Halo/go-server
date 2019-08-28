package model

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type User struct {
	Username string
	Password string
	NickName string
}

var CurrentUser = &User{}
var DB *mongo.Client

func (User) Login(username, password string) bool {
	user := &User{Username: username, Password: password}
	collection := DB.Database("chat").Collection("users")
	insertResult, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)

	return username == "admin" && password == "123"
}

func (User) SignUp(nickname, username, password string) bool {
	user := &User{Username: username, Password: password}
	collection := DB.Database("chat").Collection("users")
	insertResult, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted a single document: ", insertResult.InsertedID)

	return true
}