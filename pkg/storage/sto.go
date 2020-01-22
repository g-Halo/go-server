package storage

import (
	"github.com/g-Halo/go-server/internal/logic/model"
)

type Storage struct {
	Users    map[string]*model.User
	Rooms    map[string]*model.Room
	RoomMsgs map[string]*model.RoomMessage
}

var Sto *Storage

func NewStorage() *Storage {
	Sto = &Storage{
		Users:    map[string]*model.User{},
		Rooms:    map[string]*model.Room{},
		RoomMsgs: map[string]*model.RoomMessage{},
	}
	return Sto
}

func AddUser(user *model.User) {
	Sto.Users[user.Username] = user
}

func GetUsers() []*model.User {
	var users []*model.User
	for _, v := range Sto.Users {
		users = append(users, v)
	}
	return users
}

func UpdateUser(user *model.User) {
	_, exist := Sto.Users[user.Username]
	if exist {
		Sto.Users[user.Username] = user
	}
}

func GetUser(key string) *model.User {
	user, ok := Sto.Users[key]
	if ok {
		return user
	} else {
		return nil
	}
}

func AddRoom(room *model.Room) {
	Sto.Rooms[room.UUID] = room
}

func AddRoomMsg(rmsg *model.RoomMessage) {
	Sto.RoomMsgs[rmsg.UUID] = rmsg
}

func GetRoomMsg(uuid string) *model.RoomMessage {
	return Sto.RoomMsgs[uuid]
}

func GetRooms() []*model.Room {
	var rooms []*model.Room
	for _, v := range Sto.Rooms {
		rooms = append(rooms, v)
	}
	return rooms
}

func GetRoom(key string) *model.Room {
	room, ok := Sto.Rooms[key]
	if ok {
		return room
	} else {
		return nil
	}
}

func UpdateRoom(room *model.Room) {
	_, exist := Sto.Rooms[room.UUID]
	if exist {
		Sto.Rooms[room.UUID] = room
	}
}
