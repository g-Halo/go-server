package logic

import (
	"errors"
	"sync"

	"github.com/g-Halo/go-server/model"
	"github.com/g-Halo/go-server/storage"
	uuid "github.com/satori/go.uuid"
)

type roomLogic struct {
	mutex *sync.Mutex
}

var RoomLogic = &roomLogic{mutex: &sync.Mutex{}}

func (*roomLogic) FindOrCreate(usernames []string) *model.Room {
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
	room := Room.New(uuid.String(), usernames)
	storage.AddRoom(room)
	return room
}

func (*roomLogic) Push(key, username string, data string) error {
	currentUser := UserLogic.FindByUsername(key)
	user := UserLogic.FindByUsername(username)

	if user == nil || currentUser == nil {
		return errors.New("User Not Found")
	} else if user.Username == currentUser.Username {
		return errors.New("TargetUser can not be yourself")
	}

	room := RoomLogic.FindOrCreate([]string{currentUser.Username, user.Username})
	if room == nil {
		return errors.New("Room Not Found")
	}

	currentUser.Rooms = append(currentUser.Rooms, room)
	user.Rooms = append(user.Rooms, room)
	storage.UpdateUser(currentUser)
	storage.UpdateUser(user)

	var Message model.Message
	msg := Message.Create(currentUser, user, data)
	rChan, _ := RoomChannels.Get(room.UUID)
	rChan.PushMsg(room, msg)

	return nil
}
