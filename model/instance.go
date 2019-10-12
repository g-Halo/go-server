package model

import (
	"context"
	"github.com/yigger/go-server/conf"
	"github.com/yigger/go-server/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client

func CreateInstance(conf *conf.Config) *mongo.Client {
	if conf.No_db() {
		return nil
	}

	if mongoClient != nil {
		return mongoClient
	}

	clientOptions := options.Client().ApplyURI(conf.MongoDbAddress)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		logger.Fatal("Fail to connect MongoDB")
	}

	mongoClient = client

	return mongoClient
}

func Collection(table string) *mongo.Collection {
	database := mongoClient.Database("chat")
	return database.Collection(table)
}