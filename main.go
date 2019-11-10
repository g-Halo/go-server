package main

import (
	"github.com/g-Halo/go-server/rpc"
	"github.com/g-Halo/go-server/rpc/logic"
	"github.com/g-Halo/go-server/storage"
)

func init() {
	userLogic := logic.UserLogic

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

	sto := storage.NewStorage()
	sto.AddUser(test1)
	sto.AddUser(test2)
}

func main() {
	rpc.StartServer()
}
