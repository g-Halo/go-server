package main

import (
	"github.com/g-Halo/go-server/rpc"
	"github.com/g-Halo/go-server/rpc/logic"
	"github.com/g-Halo/go-server/storage"
)

func init() {
	// 初始化内存存储器
	storage.NewStorage()

	// 初始化测试用户
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

	storage.AddUser(test1)
	storage.AddUser(test2)
}

func main() {
	rpc.StartServer()
}
