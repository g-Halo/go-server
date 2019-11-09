package main

import (
	"github.com/g-Halo/go-server/model"
)

type Storage struct {
	Users	[]*model.User
	Rooms	[]*model.Room
}

var sto *Storage

func init() {
	// Q：初始化的时候应该给多少容量？
	sto = &Storage{
		Users: make([]*model.User, 0),
		Rooms: make([]*model.Room, 0),
	}

	// FIXME: 初始化测试用户，上线后删除
	var userLogic userLogic
	test1 := userLogic.SignUp(map[string]interface{}{
		"username": "test1",
		"nickname": "test1",
		"password": "123",
	})

	test2 := userLogic.SignUp(map[string]interface{}{
		"username": "test2",
		"nickname": "test2",
		"password": "123",
	})
	sto.Users = append(sto.Users, test1)
	sto.Users = append(sto.Users, test2)
}

func main() {
	StartRpc()
}
