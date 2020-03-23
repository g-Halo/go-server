package chanel

import (
	"sync"

	"github.com/g-Halo/go-server/conf"
	"github.com/g-Halo/go-server/internal/logic/model"
	"github.com/g-Halo/go-server/pkg/storage"
)

var UserChannelBuffer *ChannelList

type UCBuff struct {
	Username  string
	msgLength int
	head      map[string]*MessageBuffer
	last      map[string]*MessageBuffer
	mutex     *sync.Mutex
}

type MessageBuffer struct {
	message *model.Message
	next    *MessageBuffer
}

func InitUserChanBuffer() {
	UserChannelBuffer = NewChannelList(conf.Conf.RoomChannelsCount)
}

func NewUserChanBuff(Username string) *UCBuff {
	buff := &UCBuff{
		Username:  Username,
		msgLength: 0,
		head:      make(map[string]*MessageBuffer, 128),
		last:      make(map[string]*MessageBuffer, 128),
		mutex:     &sync.Mutex{},
	}
	return buff
}

func (buff *UCBuff) PushMessage(room *model.Room, message *model.Message) {
	// 1. 把消息记录到接收方的消息链表
	acceptorBuff, _ := UserChannelBuffer.Get(room.Acceptor)
	acceptorBuff.mutex.Lock()

	buffLink := &MessageBuffer{
		message: message,
		next:    nil,
	}
	if acceptorBuff.head == nil {
		acceptorBuff.head[room.Sender] = buffLink
	}
	lastNode := acceptorBuff.last
	if lastNode[room.Sender] == nil {
		lastNode[room.Sender] = buffLink
	} else {
		lastNode[room.Sender].next = buffLink
		acceptorBuff.last[room.Sender] = buffLink
	}
	acceptorBuff.msgLength += 1
	acceptorBuff.mutex.Unlock()
	// 记录此消息到数据存储器
	rmsg := storage.GetRoomMsg(room.UUID)
	rmsg.AddMessage(message)
	// 记录发送方的消息顺序Cache
	MessageCachedList.Put(room.Sender, message)
	// 记录接收方的消息顺序Cache
	MessageCachedList.Put(room.Acceptor, message)
}

func Subscribe(username string) {
	buff, _ := UserChannelBuffer.Get(username)
	for {
		// 检查是否下线
		if buff.msgLength == 0 {
			continue
		} else {
			messages := make([]*model.Message, buff.msgLength)
			buff.mutex.Lock()
			headNode := buff.head[username]
			lastNode := buff.last[username]
			buff.mutex.Unlock()

			for headNode != nil && headNode != lastNode {
				messages = append(messages, headNode.message)
				headNode = headNode.next
			}
			// TODO: 一次性发送给消息订阅者，并且存储到 msgcache

		}
	}
}
