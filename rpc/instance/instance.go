package instance

import (
	"github.com/g-Halo/go-server/conf"
	"github.com/g-Halo/go-server/logger"
	"net/rpc"
)


func NewInstance(address string) *rpc.Client {
	client, err := rpc.Dial("tcp", address)
	if err != nil {
		client.Close()
		logger.Fatal(err)
	}

	return client
}

func AuthRPC() *rpc.Client {
	return NewInstance(conf.Conf.AuthRPCAddress)
}

func LogicRPC() *rpc.Client {
	return NewInstance(conf.Conf.LogicRPCAddress)
}