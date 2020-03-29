package chanel

import (
	"log"
	"sync"

	"github.com/g-Halo/go-server/internal/logic/model"
	"github.com/g-Halo/go-server/pkg/logger"
)

// 统计当前在线的人数
var UsersOnlineSts = &sync.Map{}
var Subs = &sync.Map{}
var SubsChans = make(chan string, 64)

func Subscribe() {
	for {
		select {
		case username := <-SubsChans:
			_, alreadySub := Subs.Load(username)
			if !alreadySub {
				Subs.Store(username, true)
				go keepReceive(username)
			}
		}
	}
}

func keepReceive(username string) {
	buff, _ := UserChannelBuffer.Get(username)

	if buff == nil {
		logger.Error("订阅失败... buf is null")
	} else {
		logger.Infof("%s 订阅成功 !!!", username)
	}

	for {
		select {
		case <-buff.HasNewMessage:
			// 检查用户是否在线，如果已经下线，则跳出循环
			isOnline, _ := UsersOnlineSts.Load(username)
			if isOnline == nil || !isOnline.(bool) {
				logger.Info("%s 已下线", username)
				Subs.Delete(username)
				break
			}

			if buff.MsgLength == 0 {
				continue
			}

			messages := make([]*model.Message, 0)

			// FIXME: 如果读取消息的时候需要加锁， 如何做到不阻塞？
			buff.Mutex.Lock()
			headNode := buff.Head
			lastNode := buff.Last
			buff.Mutex.Unlock()

			for headNode != nil {
				messages = append(messages, headNode.Message)
				if headNode == lastNode {
					break
				}
				headNode = headNode.Next
			}

			// 移动头节点到最新的消息
			if headNode != nil {
				buff.Mutex.Lock()
				buff.Head = headNode.Next
				buff.Mutex.Unlock()
			}

			for _, item := range messages {
				log.Printf("服务端 statisticSrv. %s 接受到消息: %s", username, item.Body)
			}
		}
	}

}
