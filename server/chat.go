package server

import (
	"fmt"
	"io"
	"net"
)

type ChatS struct {

}

func NewChatS() *ChatS {
	return &ChatS{}
}

func (s *ChatS) Main(listener net.Listener) error {
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("connect err")
			panic(err)
		}

		buf := make([]byte, 4)
		_, err = io.ReadFull(conn, buf)
		if err != nil {
			fmt.Println("error reading", err.Error())
		}

		prot := &Protocol{}
		switch string(buf) {
		case "  V2":
			err := prot.IOLoop(conn)
			if err != nil {
				fmt.Println(err)
			}
		default:
			fmt.Println("default")
			return nil
		}
	}

	return nil
}