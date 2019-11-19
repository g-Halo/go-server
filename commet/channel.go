package commet

import (
	"github.com/g-Halo/go-server/model"
)

type Chan struct {
	RoomId   string
	MsgChan  chan *model.Message
	UserChan map[string]chan *model.Message
}

var (
	RoomChannel = make(map[string]*Chan, 256)
	UserRooms   = make(map[string][]string, 256)
)

func InitRoomChan(room *model.Room) *Chan {
	userChan := make(map[string]chan *model.Message)
	for _, username := range room.Members {
		userChan[username] = make(chan *model.Message, 256)
	}

	ch := &Chan{
		RoomId:   room.UUID,
		MsgChan:  make(chan *model.Message, 256),
		UserChan: userChan,
	}
	RoomChannel[room.UUID] = ch
	return RoomChannel[room.UUID]
}

func PushMsg(room *model.Room, msg *model.Message) {
	c, ok := RoomChannel[room.UUID]
	if !ok {
		c = InitRoomChan(room)
	}

	c.MsgChan <- msg
	for _, username := range room.Members {
		c.UserChan[username] <- msg
	}
}

func AddUserTo(user *model.User, room *model.Room) {
	UserRooms[user.Username] = append(UserRooms[user.Username], room.UUID)
}
