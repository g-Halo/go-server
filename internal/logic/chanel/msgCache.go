package chanel

import (
	"github.com/g-Halo/go-server/internal/logic/model"
)

var MessageCachedList map[string]([]*model.Message)

func InitMessageCachedList() {
	MessageCachedList = make(map[string]([]*model.Message), 128)
}

func MsgCachedPut(username string, roomUUID string, message *model.Message) {
	messageKey := username + "-" + roomUUID
	_, ok := MessageCachedList[messageKey]
	if !ok {
		MessageCachedList[messageKey] = make([]*model.Message, 64)
	}
	MessageCachedList[messageKey] = append(MessageCachedList[messageKey], message)
}
