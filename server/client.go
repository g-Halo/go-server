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

func (c *client) Stats() ClientStats {
	c.metaLock.RLock()
	clientID := c.ClientID
	//hostname := c.Hostname
	//userAgent := c.UserAgent
	//var identity string
	//var identityURL string
	//if c.AuthState != nil {
	//	identity = c.AuthState.Identity
	//	identityURL = c.AuthState.IdentityURL
	//}
	//pubCounts := make([]PubCount, 0, len(c.pubCounts))
	//for topic, count := range c.pubCounts {
	//	pubCounts = append(pubCounts, PubCount{
	//		Topic: topic,
	//		Count: count,
	//	})
	//}
	//c.metaLock.RUnlock()
	stats := ClientStats{
		Version:         "V2",
		RemoteAddress:   c.RemoteAddr().String(),
		ClientID:        clientID,
		//Hostname:        hostname,
		//UserAgent:       userAgent,
		//State:           atomic.LoadInt32(&c.State),
		//ReadyCount:      atomic.LoadInt64(&c.ReadyCount),
		//InFlightCount:   atomic.LoadInt64(&c.InFlightCount),
		//MessageCount:    atomic.LoadUint64(&c.MessageCount),
		//FinishCount:     atomic.LoadUint64(&c.FinishCount),
		//RequeueCount:    atomic.LoadUint64(&c.RequeueCount),
		//ConnectTime:     c.ConnectTime.Unix(),
		//SampleRate:      atomic.LoadInt32(&c.SampleRate),
		//TLS:             atomic.LoadInt32(&c.TLS) == 1,
		//Deflate:         atomic.LoadInt32(&c.Deflate) == 1,
		//Snappy:          atomic.LoadInt32(&c.Snappy) == 1,
		//Authed:          c.HasAuthorizations(),
		//AuthIdentity:    identity,
		//AuthIdentityURL: identityURL,
		//PubCounts:       pubCounts,
	}
	return stats
}