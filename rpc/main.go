package rpc

import (
	"fmt"
	"net"
	"net/rpc"
	"sync"

	"github.com/g-Halo/go-server/logger"
	"github.com/g-Halo/go-server/rpc/auth"
	"github.com/g-Halo/go-server/rpc/logic"
	"github.com/g-Halo/go-server/util"
)

type Rpc struct {
	waitGroup util.WaitGroupWrapper
}

func registerRpc(rpcType interface{}, address string) error {
	if err := rpc.Register(rpcType); err != nil {
		logger.Fatal(err)
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return err
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return err
	}

	fmt.Println(address)
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
		exitFunc(registerRpc(new(auth.Token), ":7071"))
		// 逻辑层 Rpc 模块
		exitFunc(registerRpc(new(logic.Logic), ":7072"))
	})

	<-exitCh
	select {}
}
