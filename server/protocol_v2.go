package server

import (
	"fmt"
	"io"
	"net"
	"sync/atomic"
)

//var separatorBytes = []byte(" ")

type protocolV2 struct {
	ctx *context
}

func (p *protocolV2) IOLoop(conn net.Conn) error {
	var err error
	var line []byte

	// 为客户端注册并保存到内存
	clientID := atomic.AddInt64(&p.ctx.chatS.clientIDSequence, 1)
	client := newClient(clientID, conn, p.ctx)
	p.ctx.chatS.AddClient(client.ID, client)
	fmt.Println("local ClientID is:", client.ID)

	// 随后循环遍历获取消息
	for {
		line, err = client.Reader.ReadSlice('\n')
		// 用户从终端 CTRL-C 退出
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

		fmt.Println(string(line))

		// 把消息转发到某(线上用户)，可考虑缓存离线的用户
		target := p.ctx.chatS.clients[2]
		conn := target.Conn
		_, err = conn.Write([]byte("welcome to connect\n"))
	}

	_ = conn.Close()
	return err
}