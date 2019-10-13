package logic

import (
	"github.com/g-Halo/go-server/server"
	"sync"
)

type roomLogic struct {
	mutex *sync.Mutex
	server *server.ChatS
}

var RoomLogic = &roomLogic{mutex: &sync.Mutex{}}

func (logic *roomLogic) Register(server *server.ChatS) {
	logic.server = server
}