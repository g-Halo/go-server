package main

import (
	"github.com/g-Halo/go-server/conf"
	"github.com/g-Halo/go-server/internal/http_api"
	"github.com/g-Halo/go-server/internal/logic/service"
	"github.com/g-Halo/go-server/pkg/logger"
	"github.com/g-Halo/go-server/pkg/rpc_client"
	"github.com/g-Halo/go-server/pkg/storage"
	"net"
	"net/http"
	"strings"
)

func init() {
	// 初始化日志文件
	logger.InitLogger("./application.log", "debug")
	// 初始化配置文件
	conf.LoadConf()
	// 初始化内存存储器
	storage.NewStorage()

	// 初始化测试用户
	service.UserService.SignUp(map[string]interface{}{
		"username": "test1",
		"nickname": "test1",
		"password": "123",
	})
	service.UserService.SignUp(map[string]interface{}{
		"username": "test2",
		"nickname": "test2",
		"password": "123",
	})
}

func Serve(listener net.Listener, handler http.Handler, proto string) error {
	server := &http.Server{
		Handler:  handler,
	}
	err := server.Serve(listener)
	if err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
		logger.Fatal("http.Serve() error - %s", err)
	}
	return nil
}

func main() {
	rpc_client.InitAuthClient(conf.Conf.AuthRPCAddress)
	rpc_client.InitLogicClient(conf.Conf.LogicRPCAddress)

	httpListener, err := net.Listen("tcp", conf.Conf.HttpApiAddress)
	if err != nil {
		logger.Fatal("error")
	}
	logger.Infof("start Listen web api in %s", conf.Conf.HttpApiAddress)
	Serve(httpListener, http_api.StartServer(), "HTTP")
}

