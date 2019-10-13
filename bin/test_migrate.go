package main

import (
	"github.com/g-Halo/go-server/model"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"context"
	"strconv"
	"time"
)


func main() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ := mongo.Connect(context.TODO(), clientOptions)

	// 生成 100 名测试用户
	users := [100]int{}
	for k, _ := range users {
		var User model.User
		params := map[string]interface{}{
			"nickname": "test-" + strconv.Itoa(k),
			"username": "test-" + strconv.Itoa(k),
			"password": "123456",
			"last_message": map[string]interface{}{
				"sender": "test-" + strconv.Itoa(k),
				"recipient": "test-" + strconv.Itoa(k + 1),
				"body": "Hi, 你好",
				"created_at": time.Now(),
				"status": "uncheck",
			},
		}
		User.SignUp(client, params)
	}
}
