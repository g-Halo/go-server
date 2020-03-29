package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/g-Halo/go-server/api/auth"
	"github.com/g-Halo/go-server/api/logic"
	"github.com/g-Halo/go-server/conf"
	"github.com/g-Halo/go-server/internal/logic/chanel"
	"github.com/g-Halo/go-server/internal/logic/service"
	"github.com/g-Halo/go-server/pkg/logger"
	"github.com/g-Halo/go-server/pkg/pb"
	"github.com/g-Halo/go-server/pkg/rpc_client"
	"github.com/g-Halo/go-server/pkg/storage"
	"github.com/g-Halo/go-server/pkg/util"
	"google.golang.org/grpc"
)

func init() {
	// 初始化日志文件
	logger.InitLogger("./logic.log", "debug")
	// 初始化配置文件
	conf.LoadConf()
	// 初始化内存存储器
	storage.NewStorage()
	// 初始化 roomChannel
	chanel.InitUserChanBuffer()
	chanel.InitMessageCachedList()
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

func main() {
	// 注册 Auth 层 RPC
	go func() {
		defer util.RecoverPanic()
		// grpc Auth Logic
		lis, err := net.Listen("tcp", conf.Conf.AuthRPCAddress)
		if err != nil {
			panic(err)
		}
		s := grpc.NewServer()
		pb.RegisterAuthServer(s, &auth.AuthServer{})
		logger.Info("Auth Rpc 注册成功")
		if err := s.Serve(lis); err != nil {
			panic(err)
		}
	}()

	// 注册逻辑层 RPC
	go func() {
		defer util.RecoverPanic()
		lis, err := net.Listen("tcp", conf.Conf.LogicRPCAddress)
		if err != nil {
			panic(err)
		}
		s := grpc.NewServer()
		pb.RegisterLogicServer(s, &logic.LogicServer{})
		logger.Info("Logic Rpc 注册成功")
		if err := s.Serve(lis); err != nil {
			panic(err)
		}
	}()

	// 连接 websocket 注册的 RPC
	rpc_client.InitWsClient(conf.Conf.WebSocketAddress)
	// 与 commet 层进行连接
	// rpc_client.ConnectToComet(conf.Conf.CommetAddress)
	go chanel.Subscribe()

	c := make(chan os.Signal)
	signal.Notify(c)
	for s := range c {
		switch s {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL:
			fmt.Println("退出", s)
			os.Exit(0)
		case syscall.SIGUSR1:
			fmt.Println("usr1", s)
		case syscall.SIGUSR2:
			fmt.Println("usr2", s)
		default:
			fmt.Println("other", s)
		}
	}
}
