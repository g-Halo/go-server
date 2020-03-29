package comet

import "net"

type Comet struct {
}

func (*Comet) IOLoop(conn net.Conn) error {
	for {

	}
}
