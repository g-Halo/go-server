package chanel

import (
	"sync"

	"github.com/g-Halo/go-server/conf"
	"github.com/g-Halo/go-server/internal/logic/model"
	"github.com/g-Halo/go-server/pkg/storage"
)

var UserChannelBuffer *ChannelList

// Remark: 实际上并不需要关注是谁发送过来的消息，这部分信息其实 message 已经包含了
type UCBuff struct {
	Username      string
	MsgLength     int
	HasNewMessage chan interface{} // 是否有新消息
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
		Mutex:         &sync.Mutex{},
	}
	return buff
}

func (buff *UCBuff) PushMessage(room *model.Room, message *model.Message) {
	// 1. 把消息记录到接收方的消息链表
	acceptorBuff, _ := UserChannelBuffer.Get(room.Acceptor)
	acceptorBuff.Mutex.Lock()

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

	// 处理上线
	acceptorBuff.Online()
	// 通知订阅者有新的消息
	acceptorBuff.HasNewMessage <- true

	// 把消息记录到 "消息-房间" 的关联关系
	rmsg := storage.GetRoomMsg(room.UUID)
	rmsg.AddMessage(message)

	// 维护发送方的消息列表
	MsgCachedPut(room.Sender, room.UUID, message)
	// 维护接收方的消息列表
	MsgCachedPut(room.Acceptor, room.UUID, message)
}

func (buff *UCBuff) Online() {
	username := buff.Username
	UsersOnlineSts.Store(username, true)
	// FIXME: SubsChans 限制了同一时间只能有 xx 个一起发送消息？
	SubsChans <- username
}
