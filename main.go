package main

import (
	"github.com/g-Halo/go-server/conf"
	"github.com/g-Halo/go-server/logger"
	"github.com/g-Halo/go-server/rpc"
	"github.com/g-Halo/go-server/rpc/logic"
	"github.com/g-Halo/go-server/storage"
)

func init() {
	// 初始化日志文件
	logger.InitLogger("./application.log", "debug")
	// 初始化配置文件
	conf.LoadConf()
	// 初始化内存存储器
	storage.NewStorage()

	// 初始化测试用户
	userLogic := logic.UserLogic
	userLogic.SignUp(map[string]interface{}{
		"username": "test1",
		"nickname": "test1",
		"password": "123",
	})
	userLogic.SignUp(map[string]interface{}{
		"username": "test2",
		"nickname": "test2",
		"password": "123",
	})
}

func main() {
	rpc.StartServer()
}
