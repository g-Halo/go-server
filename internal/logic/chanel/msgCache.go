package chanel

import (
	"sync"

	"github.com/g-Halo/go-server/internal/logic/model"
)

var MessageCachedList *MessageCached

type MessageCached struct {
	container map[string]*UMessageCached
}

type UMessageCached struct {
	message []*model.Message
	mutex   *sync.Mutex
}

func InitMessageCachedList() {
	MessageCachedList = &MessageCached{
		container: make(map[string]*UMessageCached, 64),
	}
}

func (mc *MessageCached) Put(username string, message *model.Message) {
	c, ok := mc.container[username]
	if !ok {
		c = &UMessageCached{
			message: make([]*model.Message, 64),
			mutex:   &sync.Mutex{},
		}
	}
	c.mutex.Lock()
	c.message = append(c.message, message)
	c.mutex.Unlock()
}
