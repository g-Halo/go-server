package logic

import "github.com/g-Halo/go-server/model"

type RoomChan struct {
	RoomId           string
	MsgChan          chan *model.Message
	UserDispatchChan map[string]chan *model.Message // 负责分发给某个用户
}

func NewRoomChan(RoomId string) *RoomChan {
	rc := &RoomChan{
		RoomId:           RoomId,
		MsgChan:          make(chan *model.Message, 64),
		UserDispatchChan: map[string]chan *model.Message{},
	}

	return rc
}

func (rc *RoomChan) PushMsg(key string, m *model.Message) {
	rc.MsgChan <- m
}
