package http_api

import "github.com/g-Halo/go-server/server"

type context struct {
	chatS *server.ChatS
}

func NewContext(s *server.ChatS) *context {
	return &context{
		chatS: s,
	}
}