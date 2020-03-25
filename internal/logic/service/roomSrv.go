package service

import (
	"errors"

	"github.com/g-Halo/go-server/internal/logic/chanel"
	"github.com/g-Halo/go-server/internal/logic/model"
	"github.com/g-Halo/go-server/pkg/storage"
	uuid "github.com/satori/go.uuid"
)

type roomService struct{}

var RoomService = new(roomService)

func (*roomService) FindOrCreate(sender, acceptor string) *model.Room {
	var targetRoom *model.Room
	for _, r := range storage.GetRooms() {
		if (r.Members[0] == sender && r.Members[1] == acceptor) ||
			(r.Members[0] == acceptor && r.Members[1] == sender) {
			targetRoom = r
			break
		}
	}

	if targetRoom != nil {
		return targetRoom
	}

	var Room model.Room
	uuid := uuid.NewV4()
	room, roomMsg := Room.New(uuid.String(), sender, acceptor)
	storage.AddRoom(room)
	storage.AddRoomMsg(roomMsg)
	return room
}

func (s *roomService) Push(senderUsername, receiverUsername string, data string) error {
	currentUser := UserService.FindByUsername(senderUsername)
	user := UserService.FindByUsername(receiverUsername)

	if user == nil || currentUser == nil {
		return errors.New("User Not Found")
	}

	room := s.FindOrCreate(currentUser.Username, user.Username)
	if room == nil {
		return errors.New("Room Not Found")
	}

	buff, _ := chanel.UserChannelBuffer.Get(currentUser.Username)
	// currentUser.Rooms = append(currentUser.Rooms, room)
	// user.Rooms = append(user.Rooms, room)
	// storage.UpdateUser(currentUser)
	// storage.UpdateUser(user)
	var Message model.Message
	msg := Message.Create(currentUser, user, *room, data)
	buff.PushMessage(room, msg)

	return nil
}

func (s *roomService) FindByUUID(uuid string) *model.Room {
	return storage.GetRoom(uuid)
}
