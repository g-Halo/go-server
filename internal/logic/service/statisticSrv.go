package service

import (
	"sync"

	"github.com/g-Halo/go-server/internal/logic/chanel"
)

type statisticService struct {
	mutex *sync.Mutex
}

var StatisticService = &statisticService{
	mutex: &sync.Mutex{},
}

func (s *statisticService) Online(username string) {
	chanel.UsersOnlineSts.Store(username, true)
	chanel.SubsChans <- username
}

func (s *statisticService) Offline(username string) {
	chanel.UsersOnlineSts.Store(username, false)
}

func (s *statisticService) IsOnline(username string) bool {
	isOnline, _ := chanel.UsersOnlineSts.Load(username)
	if isOnline == nil {
		return false
	} else {
		return isOnline.(bool)
	}
}
