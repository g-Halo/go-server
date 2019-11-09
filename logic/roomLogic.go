package main

import (
	"sync"
)

type roomLogic struct {
	mutex *sync.Mutex
}

var RoomLogic = &roomLogic{mutex: &sync.Mutex{}}
