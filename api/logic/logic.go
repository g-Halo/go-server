package logic

import (
	"context"
	"github.com/g-Halo/go-server/internal/logic/service"
	"github.com/g-Halo/go-server/pkg/pb"
	"github.com/g-Halo/go-server/pkg/storage"
)

type LogicServer struct {}

func (s *LogicServer) PushMessage(ctx context.Context, in *pb.PushMessageReq) (*pb.PushMessageResp, error) {
	err := service.RoomService.Push(in.GetSenderUsername(), in.GetReceiverUsername(), in.GetBody())
	return nil, err
}

func (s *LogicServer) GetUser(ctx context.Context, in *pb.GetUserReq) (*pb.GetUserResp, error) {
	user := service.UserService.FindByUsername(in.GetUsername())
	if user == nil {
		return nil, nil
	}

	pbUser := &pb.GetUserResp{
		User:                 &pb.User{
			Username:             user.Username,
			Nickname:             user.NickName,
		},
	}
	return pbUser, nil
}

func (s *LogicServer) GetUsers(ctx context.Context, in *pb.GetUsersReq) (*pb.GetUsersResp, error) {
	users := service.UserService.GetUsers()
	var pbUsers []*pb.User
	for _, user := range users {
		pbUsers = append(pbUsers, &pb.User{
			Username:             user.Username,
			Nickname:             user.NickName,
		})
	}
	resp := &pb.GetUsersResp{
		Users:                pbUsers,
	}
	return resp, nil
}

func (s *LogicServer) FindOrCreateRoom(ctx context.Context, in *pb.FindOrCreateRoomReq) (*pb.FindOrCreateRoomResp, error) {
	room := service.RoomService.FindOrCreate([]string{in.GetCurrentUsername(), in.GetTargetUsername()})
	r := &pb.FindOrCreateRoomResp{
		Uuid:                 room.UUID,
		Name:                 room.Name,
		Type:                 room.Type,
	}
	return r, nil
}

func (s *LogicServer) GetRoomById(ctx context.Context, in *pb.GetRoomByIdReq) (*pb.GetRoomByIdResp, error) {
	room := storage.GetRoom(in.Uuid)
	if room == nil {
		return nil, nil
	}

	pbRoom := &pb.GetRoomByIdResp{
		Uuid:                 room.UUID,
		Name:                 room.Name,
		Members:              room.Members,
		Type:                 room.Type,
		CreatedAt:            room.CreatedAt.Unix(),
	}
	return pbRoom, nil
}

func (s *LogicServer) GetRoomMessages(ctx context.Context, in *pb.GetRoomMessagesReq) (*pb.GetRoomMessagesResp, error) {
	r := storage.GetRoomMsg(in.GetUuid())
	var roomMessages []*pb.RoomMessage
	for _, msg := range r.Messages {
		roomMessages = append(roomMessages,&pb.RoomMessage{
			Uuid:                 in.GetUuid(),
			Sender:               msg.Sender,
			Recipient:            msg.Recipient,
			Body:                 msg.Body,
			Status:               msg.Status,
			CreatedAt:            msg.CreatedAt.Unix(),
		})
	}

	res := &pb.GetRoomMessagesResp{
		RoomMessages: roomMessages,
	}
	return res, nil
}