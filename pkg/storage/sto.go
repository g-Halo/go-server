package storage

import (
	"encoding/json"
	"github.com/g-Halo/go-server/internal/logic/model"
	"github.com/go-redis/redis/v7"
)

type Storage struct {
	Users    map[string]*model.User
	Rooms    map[string]*model.Room
	RoomMsgs map[string]*model.RoomMessage
}

var redisCli *redis.Client

func NewStorage() {
	if redisCli == nil {
		redisCli = redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		})
	}
}

func AddUser(user *model.User) {
	v, _ := json.Marshal(user)
	redisCli.HSet("users", user.Username, v)
}

func GetUsers() []*model.User {
	users := make([]*model.User, 0)
	v := redisCli.HGetAll("users")
	if v.Err() != nil {
		return users
	}

	for _, item := range v.Val() {
		var user *model.User
		_ = json.Unmarshal([]byte(item), &user)
		users = append(users, user)
	}
	return users
}

func UpdateUser(user *model.User) {
	v, _ := json.Marshal(user)
	redisCli.HSet("users", user.Username, v)
}

func GetUser(key string) (user *model.User) {
	v := redisCli.HGet("users", key)
	if v.Err() != nil {
		return nil
	}
	bytes, _ := v.Bytes()
	if err := json.Unmarshal(bytes, &user); err != nil {
		return nil
	}
	return user
}

func AddRoom(room *model.Room) {
	v, _ := json.Marshal(room)
	redisCli.HSet("rooms", room.UUID, v)
}

func AddRoomMsg(msg *model.RoomMessage) {
	v, _ := json.Marshal(msg)
	redisCli.HSet("room_msgs", msg.UUID, v)
}

func GetRoomMsg(key string) (roomMsg *model.RoomMessage) {
	v := redisCli.HGet("room_msgs", key)
	if v.Err() != nil {
		return nil
	}
	bytes, _ := v.Bytes()
	if err := json.Unmarshal(bytes, &roomMsg); err != nil {
		return nil
	}

	return roomMsg
}

func GetRooms() []*model.Room {
	rooms := make([]*model.Room, 0)
	v := redisCli.HGetAll("rooms")
	if v.Err() != nil {
		return rooms
	}

	for _, item := range v.Val() {
		var room *model.Room
		_ = json.Unmarshal([]byte(item), &room)
		rooms = append(rooms, room)
	}
	return rooms
}

func GetRoom(key string) (room *model.Room) {
	v := redisCli.HGet("rooms", key)
	if v.Err() != nil {
		return nil
	}
	bytes, _ := v.Bytes()
	if err := json.Unmarshal(bytes, &room); err != nil {
		return nil
	}
	return room
}

func UpdateRoom(room *model.Room) {
	v, _ := json.Marshal(room)
	redisCli.HSet("rooms", room.UUID, v)
}
