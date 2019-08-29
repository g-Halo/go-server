package server

import "net"

type Protocol struct {

}

func (p *Protocol) IOLoop(conn net.Conn) error {

	return nil
}