package main

import (
	"github.com/g-Halo/go-server/conf"
	"github.com/g-Halo/go-server/http_api"
	"github.com/g-Halo/go-server/logger"
	"github.com/g-Halo/go-server/logic"
	"github.com/g-Halo/go-server/model"
	"github.com/g-Halo/go-server/server"
	"runtime"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	logger.InitLogger("./application.log", "debug")
}

func main() {
	config := conf.LoadConf()
	model.CreateInstance(config)

	chatServer, err := server.NewChatS(config)
	if err != nil {
		logger.Fatal(err)
	}

	// 注册逻辑层的上下文环境
	chatServer.RegisterLogic(logic.UserLogic)
	chatServer.RegisterLogic(logic.RoomLogic)

	// 暂无注册功能，开发环境暂时预创建两用户. test-1 test-2
	if config.Env == "development" {
		test1 := logic.UserLogic.SignUp(map[string]interface{}{
			"username": "test-1",
			"nickname": "test-1",
			"password": "123",
		})

		test2 := logic.UserLogic.SignUp(map[string]interface{}{
			"username": "test-2",
			"nickname": "test-2",
			"password": "123",
		})
		chatServer.AddUser(test1)
		chatServer.AddUser(test2)
	}

	// 初始化 http api 服务
	context := http_api.NewContext(chatServer)
	httpServer := http_api.Server(context)
	httpLister := chatServer.CreateHttpListen()
	http_api.Serve(httpLister, httpServer, "HTTP")

	// 初始化 tcp 服务
	chatServer.CreateTcpListen()
	err = chatServer.Main()
	if err != nil {
		logger.Fatal(err)
	}
}
