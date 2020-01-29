package main

import (
	"github.com/g-Halo/go-server/api/auth"
	"github.com/g-Halo/go-server/api/logic"
	"github.com/g-Halo/go-server/conf"
	"github.com/g-Halo/go-server/internal/logic/chanel"
	"github.com/g-Halo/go-server/internal/logic/service"
	"github.com/g-Halo/go-server/pkg/logger"
	"github.com/g-Halo/go-server/pkg/pb"
	"github.com/g-Halo/go-server/pkg/storage"
	"github.com/g-Halo/go-server/pkg/util"
	"google.golang.org/grpc"
	"net"
)

func init() {
	// 初始化日志文件
	logger.InitLogger("./logic.log", "debug")
	// 初始化配置文件
	conf.LoadConf()
	// 初始化内存存储器
	storage.NewStorage()
	// 初始化 roomChannel
	chanel.InitRoomChan()
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

	select {}

}