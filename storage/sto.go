package storage

import (
	"github.com/g-Halo/go-server/model"
)

type Storage struct {
	Users []*model.User
	Rooms []*model.Room
}

var Sto *Storage

func NewStorage() *Storage {
	Sto = &Storage{
		Users: make([]*model.User, 0),
		Rooms: make([]*model.Room, 0),
	}
	return Sto
}

func AddUser(user *model.User) {
	Sto.Users = append(Sto.Users, user)
}

func GetUsers() []*model.User {
	return Sto.Users
}
