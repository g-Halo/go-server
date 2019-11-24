package model

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"

	"github.com/g-Halo/go-server/logger"
	"github.com/g-Halo/go-server/util"

	"crypto/md5"
	"math/rand"
	"time"
)

type User struct {
	Username    string    `json:"username"`
	Salt        string    `json:"salt"`
	Password    string    `json:"password"`
	NickName    string    `json:"nickname"`
	Rooms       []*Room   `json:"rooms"`
	CreatedAt   time.Time `json:"created_at"`
	LastMessage *Message  `json:"last_message"`

	subRoom map[string]*Room
}

func (u User) Login(user *User, password string) (string, error) {
	// 密码加盐校验
	salt := user.Salt
	m5 := md5.New()
	m5.Write([]byte(salt))
	m5.Write([]byte(string(password)))
	st := m5.Sum(nil)
	if hex.EncodeToString(st) != user.Password {
		return "", errors.New("用户名或密码错误")
	}

	// 生成 JWT token
	userInfo := make(map[string]interface{})
	userInfo["username"] = user.Username
	token := util.CreateJWT(userInfo)

	return token, nil
}

func (User) New(params map[string]interface{}) *User {
	// 生成随机 salt
	rand := func(n int) string {
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
	user := &User{}
	if err := json.Unmarshal(jsonString, &user); err != nil {
		logger.Errorf("json Unmarshal fail: %s", jsonString)
		return nil
	}

	user.Password = hex.EncodeToString(st)
	user.Salt = salt
	user.CreatedAt = time.Now()

	return user
}

// 注册用户
func (u User) Create(user *User) error {
	collection := Collection("users")
	_, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		return err
	}
	return nil
}

func (user *User) SubRoom(room *Room) {
	if _, ok := user.subRoom[room.UUID]; ok {
		return
	}
	for {
		select {
		case msg := <-room.MessageChan:
			room.AddMessage(msg)
		}
	}
}

func (u *User) ToJson() map[string]interface{} {
	return map[string]interface{}{
		"username": u.Username,
		"nickname": u.NickName,
	}
}
