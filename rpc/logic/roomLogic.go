package logic

import (
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
