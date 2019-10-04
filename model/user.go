package model

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/yigger/go-server/logger"
	"github.com/yigger/go-server/util"
	"log"

	"crypto/md5"
	"go.mongodb.org/mongo-driver/bson"
	"math/rand"
	"time"
)

type User struct {
	Username 			string 		`json:"username"`
	Salt	 			string 		`json:"salt"`
	Password 			string		`json:"password"`
	NickName 			string		`json:"nickname"`
	Rooms				[]*Room		`json:"rooms"`
	CreatedAt 			time.Time 	`json:"created_at"`
	LastMessage			*Message	`json:"last_message"`
}

func (u User) Login(username, password string) (string, error) {
	user, err := u.FindByUsername(username)
	if err != nil {
		return "", errors.New("用户名不存在")
	}

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

// 注册用户
func (User) SignUp(params map[string]interface{}) error {
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

	collection := Collection("users")
	_, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		return err
	}
	return nil
}


func (User) FindByUsername(username string) (*User, error) {
	var user User
	collection := Collection("users")

	filter := bson.M{"username": username}
	documentReturned := collection.FindOne(context.TODO(), filter)
	err := documentReturned.Decode(&user)
	if err != nil {
		return nil, err
	}

	if user.Username == "" {
		return nil, errors.New("User Not Found")
	}

	return &user, nil
}

func (User) FindAll() []*User {
	var users []*User
	collection := Collection("users")
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

func (u *User) AddRoom(room *Room) {
	u.Rooms = append(u.Rooms, room)
	filter := bson.D{{"username", u.Username}}
	update := bson.D{
		{"$set", bson.D{
			{"rooms", u.Rooms},
		}},
	}

	collection := Collection("users")
	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
}

func (u *User) FindP2PRoom(username string) *Room {
	userArray := []string{u.Username, username}
	room := &Room{}
	for _, r := range u.Rooms {
		if r.Type == "p2p" && len(r.Members) == 2 {
			if (r.Members[0] == userArray[0] || r.Members[0] == userArray[1]) || (r.Members[1] == userArray[0] || r.Members[1] == userArray[1]) {
				room = r
				break
			}
		}
	}

	// FIXME: 有可能 room 仅在 user 集合中存在，但在 Room Collection 不存在

	if room.UUID == "" {
		return nil
	} else {
		filter := bson.M{"uuid": room.UUID}
		documentReturned := Collection("rooms").FindOne(context.TODO(), filter)
		err := documentReturned.Decode(&room)
		if err != nil {
			logger.Error(err)
			return nil
		}

		return room
	}
}

func (u *User) ToJson() (map[string]interface{}) {
	return map[string]interface{}{
		"username": u.Username,
		"nickname": u.NickName,
	}
}