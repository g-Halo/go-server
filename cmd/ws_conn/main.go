package main

import (
	"github.com/g-Halo/go-server/api/ws_api"
	"github.com/g-Halo/go-server/conf"
	"github.com/g-Halo/go-server/internal/ws_conn"
	"github.com/g-Halo/go-server/pkg/logger"
	"github.com/g-Halo/go-server/pkg/pb"
	"github.com/g-Halo/go-server/pkg/storage"
	"github.com/g-Halo/go-server/pkg/util"
	"google.golang.org/grpc"
	"net"
)

func init() {
	// 初始化日志文件
	logger.InitLogger("./ws_conn.log", "debug")
	// 初始化配置文件
	conf.LoadConf()
	// 初始化内存存储器
	storage.NewStorage()
}

func main() {
	go func() {
		defer util.RecoverPanic()
		lis, err := net.Listen("tcp", conf.Conf.WebSocketAddress)
		if err != nil {
			panic(err)
		}
		s := grpc.NewServer()
		pb.RegisterWsConnServer(s, &ws_api.WsConn{})
		logger.Info("Websocket Rpc 注册成功")
		if err := s.Serve(lis); err != nil {
			panic(err)
		}
	}()

	// 单独启动 websocket，用于与 js 端的 ws 进行连接
	ws_conn.StartServer()
}