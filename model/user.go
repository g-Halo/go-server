package model

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/yigger/go-server/util"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"crypto/md5"
	"math/rand"
	"time"
)

type User struct {
	Username 			string 		`json:"username"`
	Salt	 			string 		`json:"salt"`
	Password 			string		`json:"password"`
	NickName 			string		`json:"nickname"`
	CreatedAt 			time.Time 	`json:"created_at"`
	LastMessage			*Message	`json:"last_message"`
}

func (u User) Login(client *mongo.Client, username, password string) (string, error) {
	queryFilter := bson.M{"username": username}
	user, err := u.FindByUsername(client, queryFilter)
	if err != nil {
		return "", errors.New("用户名或密码错误 - 0")
	}

	// 密码加盐校验
	salt := user.Salt
	m5 := md5.New()
	m5.Write([]byte(salt))
	m5.Write([]byte(string(password)))
	st := m5.Sum(nil)
	if hex.EncodeToString(st) != user.Password {
		return "", errors.New("用户名或密码错误 - 1")
	}

	// 生成 JWT token
	userInfo := make(map[string]interface{})
	userInfo["username"] = user.Username
	token := util.CreateJWT(userInfo)

	return token, nil
}

// 注册用户
func (User) SignUp(client *mongo.Client, params map[string]interface{}) error {
	//nickname, username, password string
	// 生成随机 salt
	rand := func (n int) string {
		var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
		b := make([]rune, n)
		for i := range b {
			b[i] = letterRunes[rand.Intn(len(letterRunes))]
		}
		return string(b)
	}

	// 加密处理
	m5 := md5.New()
	salt := rand(8)
	m5.Write([]byte(salt))
	m5.Write([]byte(params["password"].(string)))
	st := m5.Sum(nil)

	jsonString, _ := json.Marshal(params)

	var user User
	if err := json.Unmarshal(jsonString, &user); err != nil {
		log.Fatal("json Unmarshal fail")
	}

	user.Password = hex.EncodeToString(st)
	user.Salt = salt
	user.CreatedAt = time.Now()

	collection := client.Database("chat").Collection("users")
	_, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		return err
	}
	return nil
}


func (User) FindByUsername(client *mongo.Client, filter bson.M) (User, error) {
	var user User
	collection := client.Database("chat").Collection("users")
	documentReturned := collection.FindOne(context.TODO(), filter)
	err := documentReturned.Decode(&user)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (User) FindAll(client *mongo.Client) []*User {
	var users []*User
	collection := client.Database("chat").Collection("users")
	cur, err := collection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(context.TODO()) {
		var user User
		err := cur.Decode(&user)
		if err != nil {
			log.Fatal(err)
		}

		users = append(users, &user)
	}

	cur.Close(context.TODO())

	return users
}