package chanel

import (
	"github.com/g-Halo/go-server/internal/logic/model"
	"github.com/g-Halo/go-server/pkg/storage"
	"sync"
)

type RoomChan struct {
	RoomId           string
	MsgChan          chan *model.Message
	UserDispatchChan map[string]chan *model.Message // 负责分发给某个用户

	mutex *sync.Mutex
}

func NewRoomChan(RoomId string) *RoomChan {
	rc := &RoomChan{
		RoomId:           RoomId,
		MsgChan:          make(chan *model.Message, 512),
		UserDispatchChan: map[string]chan *model.Message{},
		mutex:            &sync.Mutex{},
	}

	room := storage.GetRoom(RoomId)
	if room != nil {
		for _, username := range room.Members {
			rc.UserDispatchChan[username] = make(chan *model.Message, 512)
		}
	}

	// 监听 MsgChan，并分发给用户
	// TODO: 如何确保用户都收到了消息？
	go func() {
		for {
			msg := <-rc.MsgChan
			for _, item := range rc.UserDispatchChan {
				item <- msg
			}
		}
	}()

	return rc
}

func (rc *RoomChan) PushMsg(room *model.Room, message *model.Message) {
	rc.MsgChan <- message
	roomMsg := storage.GetRoomMsg(room.UUID)
	roomMsg.AddMessage(message)
}

func (rc *RoomChan) GetMsg(key string) *model.Message {
	c, ok := rc.UserDispatchChan[key]
	if !ok {
		return nil
	}

	return <-c
}
