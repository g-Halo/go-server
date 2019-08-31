package server

import (
	"fmt"
	"github.com/yigger/go-server/protocol"
	"io"
	"net"
)

type tcpServer struct {
	ctx *context
}

func (p *tcpServer) Handle(conn net.Conn) {
	buf := make([]byte, 6)
	_, err := io.ReadFull(conn, buf)
	if err != nil {
		fmt.Println("error reading", err.Error())
	}

	var prot protocol.Protocol
	switch string(buf) {
	case "  CHAT":
		prot = &protocolV2{ctx: p.ctx}
	default:
		//fmt.Println("Tcp file: handle switch default")
		panic("Tcp file: handle switch default")
	}

	err = prot.IOLoop(conn)
	if err != nil {
		fmt.Println(err)
	}

}
