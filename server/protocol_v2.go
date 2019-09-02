package server

import (
	"bytes"
	"encoding/binary"
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

		//var response []byte
		params := bytes.Split(line, separatorBytes)
		p.Exec(client, params)

		//fmt.Println(string(line))
		//
		//// 把消息转发到某(线上用户)，可考虑缓存离线的用户
		//target := p.ctx.chatS.clients[2]
		//conn := target.Conn
		//_, err = conn.Write([]byte("welcome to connect\n"))
	}

	_ = conn.Close()
	return err
}

func (p *protocolV2) Exec(client *client, params [][]byte) {
	switch {
	case bytes.Equal(params[0], []byte("LOGIN")):
		p.Login(client, params)
	case bytes.Equal(params[0], []byte("CREATE_ROOM")):
		p.CreateRoom(client, params)
	}
}

func (p *protocolV2) Login(client *client, params [][]byte) {
	messageBody := params[1]
	fmt.Println(messageBody)
}

func (p *protocolV2) CreateRoom(client *client, params [][]byte) {
	if len(params) < 2 {
		fmt.Println("无效的协议")
		return
	}

	// 房间名
	roomName := string(params[1])
	if len(roomName) < 1 || len(roomName) > 255 {
		fmt.Println("房间名字符数仅允许 1 - 255 个字符")
		return
	}

	room := p.ctx.chatS.GetOrCreateByRoom(roomName)
	// 房间记录当前连接数
	room.AddClient(client)

	// 订阅房间消息，不断从房间的内容中读取消息，并返回给用户
	go client.SubRoom(room)
}

func readLen(r io.Reader, tmp []byte) (int32, error) {
	_, err := io.ReadFull(r, tmp)
	if err != nil {
		return 0, err
	}
	return int32(binary.BigEndian.Uint32(tmp)), nil
}