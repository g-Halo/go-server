package logic

import (
	"context"
	"github.com/g-Halo/go-server/pkg/mq"

	"github.com/g-Halo/go-server/internal/logic/service"
	"github.com/g-Halo/go-server/pkg/pb"
	"github.com/g-Halo/go-server/pkg/storage"
)

type Server struct {
	Mq *mq.MQ
}

func New() *Server {
	return &Server{
		Mq: mq.New(),
	}
}

func (s *Server) PushMessage(ctx context.Context, in *pb.PushMessageReq) (*pb.PushMessageResp, error) {
	err := service.RoomService.Push(in.GetSenderUsername(), in.GetReceiverUsername(), in.GetBody())
	return &pb.PushMessageResp{}, err
}

func (s *Server) GetUser(ctx context.Context, in *pb.GetUserReq) (*pb.GetUserResp, error) {
	user := service.UserService.FindByUsername(in.GetUsername())
	if user == nil {
		return nil, nil
	}

	var rooms []*pb.Room
	for _, room := range user.Rooms {
		rooms = append(rooms, &pb.Room{
			Uuid:      room.UUID,
			Name:      room.Name,
			Members:   room.Members,
			Type:      room.Type,
			CreatedAt: room.CreatedAt.Unix(),
		})
	}

	pbUser := &pb.GetUserResp{
		User: &pb.User{
			Username: user.Username,
			Nickname: user.NickName,
			Rooms:    rooms,
		},
	}
	return pbUser, nil
}

func (s *Server) GetUsers(ctx context.Context, in *pb.GetUsersReq) (*pb.GetUsersResp, error) {
	users := service.UserService.GetUsers()
	var pbUsers []*pb.User
	for _, user := range users {
		pbUsers = append(pbUsers, &pb.User{
			Username: user.Username,
			Nickname: user.NickName,
		})
	}
	resp := &pb.GetUsersResp{
		Users: pbUsers,
	}
	return resp, nil
}

func (s *Server) FindOrCreateRoom(ctx context.Context, in *pb.FindOrCreateRoomReq) (*pb.FindOrCreateRoomResp, error) {
	room := service.RoomService.FindOrCreate(in.GetCurrentUsername(), in.GetTargetUsername())
	r := &pb.FindOrCreateRoomResp{
		Uuid: room.UUID,
		Name: room.Name,
		Type: room.Type,
	}
	return r, nil
}

func (s *Server) GetRoomById(ctx context.Context, in *pb.GetRoomByIdReq) (*pb.GetRoomByIdResp, error) {
	room := storage.GetRoom(in.Uuid)
	if room == nil {
		return nil, nil
	}

	pbRoom := &pb.GetRoomByIdResp{
		Uuid:      room.UUID,
		Name:      room.Name,
		Members:   room.Members,
		Type:      room.Type,
		CreatedAt: room.CreatedAt.Unix(),
	}
	return pbRoom, nil
}

func (s *Server) GetRoomMessages(ctx context.Context, in *pb.GetRoomMessagesReq) (*pb.GetRoomMessagesResp, error) {
	r := storage.GetRoomMsg(in.GetUuid())
	var roomMessages []*pb.RoomMessage
	for _, msg := range r.Messages {
		roomMessages = append(roomMessages, &pb.RoomMessage{
			Uuid:      in.GetUuid(),
			Sender:    msg.Sender,
			Recipient: msg.Recipient,
			Body:      msg.Body,
			Status:    msg.Status,
			CreatedAt: msg.CreatedAt.Unix(),
		})
	}

	res := &pb.GetRoomMessagesResp{
		RoomMessages: roomMessages,
	}
	return res, nil
}

func (s *Server) KeepGetMessage(ctx context.Context, in *pb.KeepGetMessageReq) (*pb.KeepGetMessageResp, error) {
	// rChan, _ := chanel.RoomChannels.Get(in.GetUuid())
	// msg := rChan.GetMsg(in.GetUsername())
	// if msg == nil {
	// 	return &pb.KeepGetMessageResp{}, errors.New("没有更多的消息了")
	// }

	// var msgItems []*pb.KeepGetMessageItem

	// // TODO: ------ 批量获取未读消息，合并发送 -------
	// room := service.RoomService.FindByUUID(in.GetUuid())
	// sender := service.UserService.FindByUsername(msg.Sender)
	// receiver := service.UserService.FindByUsername(msg.Recipient)
	// msgItems = append(msgItems, &pb.KeepGetMessageItem{
	// 	Body: msg.Body,
	// 	Recipient: &pb.User{
	// 		Username: receiver.Username,
	// 		Nickname: receiver.NickName,
	// 	},
	// 	Sender: &pb.User{
	// 		Username: sender.Username,
	// 		Nickname: sender.NickName,
	// 	},
	// 	CreatedAt: msg.CreatedAt.Unix(),
	// 	Status:    msg.Status,
	// })
	// // ------ 批量获取未读消息，合并发送 -------

	// resp := &pb.KeepGetMessageResp{
	// 	Sender: &pb.User{
	// 		Username: receiver.Username,
	// 		Nickname: receiver.NickName,
	// 	},
	// 	Accepter: &pb.User{
	// 		Username: sender.Username,
	// 		Nickname: sender.NickName,
	// 	},
	// 	Room: &pb.Room{
	// 		Uuid: room.UUID,
	// 	},
	// 	Messages: msgItems,
	// }
	// return resp, nil
	return nil, nil
}

func (s *Server) UserOnline(ctx context.Context, in *pb.UserOnlineReq) (*pb.UserOnlineResp, error) {
	service.StatisticService.Online(in.GetUsername())
	return &pb.UserOnlineResp{}, nil
}

func (s *Server) UserOffline(ctx context.Context, in *pb.UserOnlineReq) (*pb.UserOnlineResp, error) {
	service.StatisticService.Offline(in.GetUsername())
	return &pb.UserOnlineResp{}, nil
}
