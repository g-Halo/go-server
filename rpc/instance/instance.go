package instance

import (
	"net/rpc"

	"github.com/g-Halo/go-server/conf"
	"github.com/g-Halo/go-server/logger"
)

var _authRPC *rpc.Client
var _logicRPC *rpc.Client

func NewInstance(address string) *rpc.Client {
	client, err := rpc.Dial("tcp", address)
	if err != nil {
		client.Close()
		logger.Fatal(err)
	}

	return client
}

func AuthRPC() *rpc.Client {
	if _authRPC == nil {
		_authRPC = NewInstance(conf.Conf.AuthRPCAddress)
	}
	return _authRPC
}

func LogicRPC() *rpc.Client {
	if _logicRPC == nil {
		_logicRPC = NewInstance(conf.Conf.LogicRPCAddress)
	}
	return _logicRPC
}
