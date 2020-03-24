package ws_conn

import (
	"context"
	"sync"

	"github.com/g-Halo/go-server/pkg/pb"
	"github.com/g-Halo/go-server/pkg/rpc_client"
)

var manage sync.Map

func store(username string, conn *Client) {
	rpc_client.LogicClient.UserOnline(context.Background(), &pb.UserOnlineReq{
		Username: username,
	})
	manage.Store(username, conn)
}

func getConn(username string) *Client {
	v, ok := manage.Load(username)
	if ok {
		return v.(*Client)
	}

	return nil
}
