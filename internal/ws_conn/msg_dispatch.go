package ws_conn

import (
	"context"
	"errors"
	"fmt"
	"github.com/g-Halo/go-server/pkg/logger"
	"github.com/g-Halo/go-server/pkg/pb"
)

func DispatchMessage(ctx context.Context, req *pb.DispatchReq) error {
	client := getConn(req.GetAccepter().Username)
	if client == nil {
		logger.Error("getWebSocketClient is null")
		return errors.New(fmt.Sprintf("empty websocket connect, username %s", req.GetAccepter().Username))
	}

	client.DispatchMessage(ctx, req)
	return nil
}