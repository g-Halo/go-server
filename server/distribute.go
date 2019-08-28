package server

import (
	"fmt"
	"net"
)

type Server struct {
	Conn net.Conn
}

func (s *Server) Main() {
	conn := s.Conn
	buf := make([]byte, 412)
	len, err := conn.Read(buf)
	if err != nil {
		fmt.Println("error reading", err.Error())
	}
	//switch string() {
	//
	//}
}