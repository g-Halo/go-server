package chanel

import (
	"sync"

	"github.com/g-Halo/go-server/conf"
	"github.com/g-Halo/go-server/internal/logic/model"

	// "github.com/g-Halo/go-server/internal/logic/service"
	"github.com/g-Halo/go-server/pkg/storage"
)

var UserChannelBuffer *ChannelList

type UCBuff struct {
	Username  string
	MsgLength int
	Head      map[string]*MessageBuffer
	Last      map[string]*MessageBuffer
	Mutex     *sync.Mutex
}

type MessageBuffer struct {
	Message *model.Message
	Next    *MessageBuffer
}

func InitUserChanBuffer() {
	UserChannelBuffer = NewChannelList(conf.Conf.RoomChannelsCount)
}

func NewUserChanBuff(Username string) *UCBuff {
	buff := &UCBuff{
		Username:  Username,
		MsgLength: 0,
		Head:      make(map[string]*MessageBuffer, 128),
		Last:      make(map[string]*MessageBuffer, 128),
		Mutex:     &sync.Mutex{},
	}
	return buff
}

func (buff *UCBuff) PushMessage(room *model.Room, message *model.Message) {
	// 1. 把消息记录到接收方的消息链表
	acceptorBuff, _ := UserChannelBuffer.Get(room.Acceptor)
	acceptorBuff.Mutex.Lock()

	buffLink := &MessageBuffer{
		Message: message,
		Next:    nil,
	}
	if acceptorBuff.Head == nil {
		acceptorBuff.Head[room.Sender] = buffLink
	}
	lastNode := acceptorBuff.Last
	if lastNode[room.Sender] == nil {
		lastNode[room.Sender] = buffLink
	} else {
		lastNode[room.Sender].Next = buffLink
		acceptorBuff.Last[room.Sender] = buffLink
	}
	acceptorBuff.MsgLength += 1
	acceptorBuff.Mutex.Unlock()
	// 记录此消息到数据存储器
	rmsg := storage.GetRoomMsg(room.UUID)
	rmsg.AddMessage(message)
	// 记录发送方的消息顺序Cache
	MsgCachedPut(room.Sender, room.UUID, message)
	// 记录接收方的消息顺序Cache
	MsgCachedPut(room.Acceptor, room.UUID, message)
}
