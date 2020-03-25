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

		time.Sleep(time.Second * 3)
	}
}
