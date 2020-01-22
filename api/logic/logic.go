package logic

import (
	"context"
	"github.com/g-Halo/go-server/internal/logic/service"
	"github.com/g-Halo/go-server/pkg/pb"
)

type LogicServer struct {}

func (s *LogicServer) PushMessage(ctx context.Context, in *pb.PushMessageReq) (*pb.PushMessageResp, error) {
	err := service.RoomService.Push(in.GetSenderUsername(), in.GetReceiverUsername(), in.GetBody())
	return nil, err
}