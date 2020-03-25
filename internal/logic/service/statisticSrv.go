package service

import (
	"log"
	"sync"
	"time"

	"github.com/g-Halo/go-server/internal/logic/chanel"
	"github.com/g-Halo/go-server/internal/logic/model"
	"github.com/g-Halo/go-server/pkg/logger"
)

type statisticService struct {
	mutex *sync.Mutex
}

var StatisticService = &statisticService{
	mutex: &sync.Mutex{},
}

// 统计当前在线的人数
var userStatistic = make(map[string]bool, 128)

func (s *statisticService) Online(username string) {
	s.mutex.Lock()
	userStatistic[username] = true
	s.mutex.Unlock()
	go Subscribe(username)
}

func (s *statisticService) Offline(username string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	userStatistic[username] = false
}

func (s *statisticService) IsOnline(username string) bool {
	isOnline, _ := userStatistic[username]
	return isOnline
}

func Subscribe(username string) {
	buff, _ := chanel.UserChannelBuffer.Get(username)

	if buff == nil {
		logger.Error("订阅失败... buf is null")
	} else {
		logger.Infof("%s 订阅成功 !!!", username)
		// logger.Info(buff)
	}

	for {
		// 检查用户是否在线
		// if !StatisticService.IsOnline(username) {
		// 	break
		// }

		if buff.MsgLength == 0 {

		} else {
			messages := make([]*model.Message, 0)
			// buff.Mutex.Lock()
			// FIXME: 如果读取消息的时候需要加锁， 如何做到不阻塞？
			headNode := buff.Head

			logger.Info(headNode)
			for headNode != nil {
				messages = append(messages, headNode.Message)
				// if headNode == lastNode {
				// 	break
				// }
				headNode = headNode.Next
			}
			// TODO: 一次性发送给消息订阅者，并且存储到 msgcache
			log.Printf("%s 接受到消息", username)
			log.Println(messages)
			// buff.Mutex.Unlock()
			time.Sleep(time.Second * 3)
		}
	}
}
