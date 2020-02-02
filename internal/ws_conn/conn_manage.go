package ws_conn

import (
	"sync"
)

var manage sync.Map

func store(username string, conn *Client) {
	manage.Store(username, conn)
}

func getConn(username string) *Client {
	v, ok := manage.Load(username)
	if ok {
		return v.(*Client)
	}

	return nil
}