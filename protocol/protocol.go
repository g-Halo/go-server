package protocol

import "net"

// Protocol describes the basic behavior of any protocol in the system
type Protocol interface {
	IOLoop(conn net.Conn) error
}

