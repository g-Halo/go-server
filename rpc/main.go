package rpc

import (
	"net"
	"net/rpc"
	"sync"

	"github.com/g-Halo/go-server/conf"
	"github.com/g-Halo/go-server/logger"
	"github.com/g-Halo/go-server/rpc/auth"
	"github.com/g-Halo/go-server/rpc/instance"
	"github.com/g-Halo/go-server/rpc/logic"
	"github.com/g-Halo/go-server/util"
)

type Rpc struct {
	waitGroup util.WaitGroupWrapper
}

var Client map[string]*rpc.Client

func registerRpc(serverName string, rpcType interface{}, address string) error {
	if err := rpc.Register(rpcType); err != nil {
		logger.Fatal(err)
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		logger.Fatal(err)
		return err
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		logger.Fatal(err)
		return err
	}

	logger.Infof("Start listen in %s", address)
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				logger.Error("accept error")
				continue
			}
			rpc.ServeConn(conn)
		}
	}()

	instance.NewInstance(serverName, address)
	return nil
}

func StartServer() {
	var exitCh = make(chan error)
	var once sync.Once
	exitFunc := func(err error) {
		once.Do(func() {
			exitCh <- err
		})
	}

	waitGroup := &util.WaitGroupWrapper{}
	waitGroup.Wrap(func() {
		// 权限校验 Rpc 模块
		exitFunc(registerRpc("auth", new(auth.Token), conf.Conf.AuthRPCAddress))
		// 逻辑层 Rpc 模块
		exitFunc(registerRpc("logic", new(logic.Logic), conf.Conf.LogicRPCAddress))
	})

	<-exitCh
	select {}
}
