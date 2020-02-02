package ws_api

import (
	"context"
	"fmt"
	"github.com/g-Halo/go-server/internal/ws_conn"
	"github.com/g-Halo/go-server/pkg/logger"
	"github.com/g-Halo/go-server/pkg/pb"
)

type WsConn struct {}

func (ws *WsConn) Dispatch(ctx context.Context, in *pb.DispatchReq) (*pb.DispatchResp, error) {
	logger.Info(fmt.Sprintf("接收到来着 %s 的消息，正在通过 RPC 分发到 %s...", in.GetSender().Username, in.GetAccepter().Username))
	return &pb.DispatchResp{}, ws_conn.DispatchMessage(ctx, in)
}