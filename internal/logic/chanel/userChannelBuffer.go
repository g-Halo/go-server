package chanel

import (
	"sync"

	"github.com/g-Halo/go-server/conf"
	"github.com/g-Halo/go-server/internal/logic/model"
	// "github.com/g-Halo/go-server/internal/logic/service"
)

var UserChannelBuffer *ChannelList

// Remark: 实际上并不需要关注是谁发送过来的消息，这部分信息其实 message 已经包含了
type UCBuff struct {
	Username      string
	MsgLength     int
	HasNewMessage chan interface{} // 是否有新消息？
	Head          *MessageBufferNode
	Last          *MessageBufferNode
	Mutex         *sync.Mutex
}

type MessageBufferNode struct {
	Message *model.Message
	Next    *MessageBufferNode
}

func InitUserChanBuffer() {
	UserChannelBuffer = NewChannelList(conf.Conf.RoomChannelsCount)
}

func NewUserChanBuff(Username string) *UCBuff {
	buff := &UCBuff{
		Username:      Username,
		MsgLength:     0,
		HasNewMessage: make(chan interface{}, 1),
		// Head:      make(map[string]*MessageBufferNode, 128),
		// Last:      make(map[string]*MessageBufferNode, 128),
		Mutex: &sync.Mutex{},
	}
	return buff
}

func (buff *UCBuff) PushMessage(room *model.Room, message *model.Message) {
	// 1. 把消息记录到接收方的消息链表
	acceptorBuff, _ := UserChannelBuffer.Get(room.Acceptor)
	acceptorBuff.Mutex.Lock()

	// acceptorBuff.Sender[room.Sender] = true

	buffLink := &MessageBufferNode{
		Message: message,
		Next:    nil,
	}

	if acceptorBuff.Head == nil {
		acceptorBuff.Head = buffLink
	}

	if acceptorBuff.Last == nil {
		acceptorBuff.Last = buffLink
	} else {
		acceptorBuff.Last.Next = buffLink
		acceptorBuff.Last = buffLink
	}
	acceptorBuff.MsgLength += 1
	acceptorBuff.Mutex.Unlock()
	// 通知订阅者有新的消息
	acceptorBuff.HasNewMessage <- true
	// // 记录此消息到数据存储器
	// rmsg := storage.GetRoomMsg(room.UUID)
	// rmsg.AddMessage(message)
	// // 记录发送方的消息顺序Cache
	// MsgCachedPut(room.Sender, room.UUID, message)
	// // 记录接收方的消息顺序Cache
	// MsgCachedPut(room.Acceptor, room.UUID, message)
}
