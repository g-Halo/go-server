package server

import (
	"fmt"
	"io"
	"net"
	"sync/atomic"
)

var separatorBytes = []byte(" ")

type protocolV2 struct {
	ctx *context
}

func (p *protocolV2) IOLoop(conn net.Conn) error {
	var err error
	var line []byte

	clientID := atomic.AddInt64(&p.ctx.chatS.clientIDSequence, 1)
	client := newClient(clientID, conn, p.ctx)
	p.ctx.chatS.AddClient(client.ID, client)
	fmt.Println("localClient:", client.ID)
	for {
		line, err = client.Reader.ReadSlice('\n')
		if err != nil {
			if err == io.EOF {
				err = nil
			} else {
				err = fmt.Errorf("failed to read command - %s", err)
			}
			break
		}
		// trim the '\n'
		line = line[:len(line)-1]
		// optionally trim the '\r'
		if len(line) > 0 && line[len(line)-1] == '\r' {
			line = line[:len(line)-1]
		}
		//params := bytes.Split(line, separatorBytes)
		fmt.Println(string(line))

	}
	return err
}