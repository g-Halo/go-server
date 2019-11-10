package logic

import (
	"sync"
)

type roomLogic struct {
	mutex *sync.Mutex
}

var RoomLogic = &roomLogic{mutex: &sync.Mutex{}}
