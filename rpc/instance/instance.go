package instance

import (
	"net/rpc"
	"os"

	"github.com/g-Halo/go-server/logger"
)

var (
	AuthRPC  *rpc.Client
	LogicRPC *rpc.Client
)

func NewInstance(serverName string, address string) *rpc.Client {
	client, err := rpc.Dial("tcp", address)
	// defer client.Close()
	if err != nil {
		logger.Fatal(err)
		os.Exit(0)
	}

	if serverName == "auth" {
		AuthRPC = client
	} else if serverName == "logic" {
		LogicRPC = client
	}
	return client
}
