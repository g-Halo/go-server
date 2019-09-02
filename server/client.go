package server

import (
	"bufio"
	"net"
	"sync"
)

const defaultBufferSize = 16 * 1024

type client struct {
	ID      int64
	net.Conn

	writeLock sync.RWMutex
	metaLock  sync.RWMutex

	// reading/writing interfaces
	Reader *bufio.Reader
	Writer *bufio.Writer

	// re-usable buffer for reading the 4-byte lengths off the wire
	lenBuf   [4]byte
	lenSlice []byte

	ClientID string
}

func newClient(id int64, conn net.Conn, ctx *context) *client {
	var identifier string
	if conn != nil {
		identifier, _, _ = net.SplitHostPort(conn.RemoteAddr().String())
	}

	c := &client{
		ID:  id,
		Conn: conn,
		Reader: bufio.NewReaderSize(conn, defaultBufferSize),
		Writer: bufio.NewWriterSize(conn, defaultBufferSize),
		ClientID: identifier,
	}
	c.lenSlice = c.lenBuf[:]
	return c
}

func (c *client) SubRoom(room *room) {

}