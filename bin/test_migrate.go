package main

import (
	"github.com/yigger/go-server/model"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"context"
)

// 生成测试用户
func main() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ := mongo.Connect(context.TODO(), clientOptions)
	var User model.User
	params := map[string]interface{}{
		"nickname": "test-01",
		"username": "test-01",
		"password": "123456",
	}
	params["last_message"] = map[string]interface{}{

	}


	User.SignUp(client, params)
}
