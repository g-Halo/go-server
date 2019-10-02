package model

import (
	"github.com/yigger/go-server/conf"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"

	"context"
)

var mongoClient *mongo.Client

func CreateInstance(conf *conf.Config) *mongo.Client {
	if mongoClient != nil {
		return mongoClient
	}

	clientOptions := options.Client().ApplyURI(conf.MongoDbAddress)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal("Fail to connect MongoDB")
	}

	mongoClient = client

	return mongoClient
}

func Collection(table string) *mongo.Collection {
	database := mongoClient.Database("chat")
	return database.Collection(table)
}