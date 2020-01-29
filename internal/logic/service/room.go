package service

import (
	"errors"
	"github.com/g-Halo/go-server/internal/logic/chanel"
	"github.com/g-Halo/go-server/internal/logic/model"
	"github.com/g-Halo/go-server/pkg/storage"
	uuid "github.com/satori/go.uuid"
)

type roomService struct {}

var RoomService = new(roomService)

func (*roomService) FindOrCreate(usernames []string) *model.Room {
	var targetRoom *model.Room
	for _, r := range storage.GetRooms() {
		if (r.Members[0] == usernames[0] && r.Members[1] == usernames[1]) ||
			(r.Members[0] == usernames[1] && r.Members[1] == usernames[0]) {
			targetRoom = r
			break
		}
	}

	if targetRoom != nil {
		return targetRoom
	}

	var Room model.Room
	uuid := uuid.NewV4()
	room, roomMsg := Room.New(uuid.String(), usernames)
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

	room := s.FindOrCreate([]string{currentUser.Username, user.Username})
	if room == nil {
		return errors.New("Room Not Found")
	}

	rChan, _ := chanel.RoomChannels.Get(room.UUID)
	currentUser.Rooms = append(currentUser.Rooms, room)
	user.Rooms = append(user.Rooms, room)
	storage.UpdateUser(currentUser)
	storage.UpdateUser(user)

	var Message model.Message
	msg := Message.Create(currentUser, user, *room, data)
	rChan.PushMsg(room, msg)

	return nil
}