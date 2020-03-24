package service

import (
	"log"
	"sync"

	"github.com/g-Halo/go-server/internal/logic/chanel"
	"github.com/g-Halo/go-server/internal/logic/model"
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
	for {
		// 检查用户是否在线
		// if !StatisticService.IsOnline(username) {
		// 	break
		// }

		if buff.MsgLength == 0 {
			continue
		} else {
			messages := make([]*model.Message, buff.MsgLength)
			buff.Mutex.Lock()
			headNode := buff.Head[username]
			lastNode := buff.Last[username]
			buff.Mutex.Unlock()

			for headNode != nil && headNode != lastNode {
				messages = append(messages, headNode.Message)
				headNode = headNode.Next
			}
			// TODO: 一次性发送给消息订阅者，并且存储到 msgcache
			log.Println("接受到消息")
			log.Println(messages)
			break
		}
	}
}
