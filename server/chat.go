package server

import (
	"github.com/yigger/go-server/protocol"
	"github.com/yigger/go-server/util"
	"net"
	"sync"
	"time"
)

type ClientStats struct {
	ClientID        string `json:"client_id"`
	Hostname        string `json:"hostname"`
	Version         string `json:"version"`
	RemoteAddress   string `json:"remote_address"`
	State           int32  `json:"state"`
	ReadyCount      int64  `json:"ready_count"`
	InFlightCount   int64  `json:"in_flight_count"`
	MessageCount    uint64 `json:"message_count"`
	FinishCount     uint64 `json:"finish_count"`
	RequeueCount    uint64 `json:"requeue_count"`
	ConnectTime     int64  `json:"connect_ts"`
	SampleRate      int32  `json:"sample_rate"`
	Deflate         bool   `json:"deflate"`
	Snappy          bool   `json:"snappy"`
	UserAgent       string `json:"user_agent"`
	Authed          bool   `json:"authed,omitempty"`
	AuthIdentity    string `json:"auth_identity,omitempty"`
	AuthIdentityURL string `json:"auth_identity_url,omitempty"`

	TLS                           bool   `json:"tls"`
	CipherSuite                   string `json:"tls_cipher_suite"`
	TLSVersion                    string `json:"tls_version"`
	TLSNegotiatedProtocol         string `json:"tls_negotiated_protocol"`
	TLSNegotiatedProtocolIsMutual bool   `json:"tls_negotiated_protocol_is_mutual"`
}

type Client interface {
	Stats() ClientStats
}

type ChatS struct {
	// 64bit atomic vars need to be first for proper alignment on 32bit platforms
	clientIDSequence int64

	sync.RWMutex
	tcpListener   		 net.Listener
	address				 string
	waitGroup            util.WaitGroupWrapper
	startTime			 time.Time
	clients 			 map[int64]Client

}

func NewChatS() (chat *ChatS, err error) {
	chat = &ChatS{
		startTime: time.Now(),
		address: "localhost:5000",
		clients: make(map[int64]Client),
	}
	chat.tcpListener, err = net.Listen("tcp", "localhost:5000")
	return
}

func (c *ChatS) AddClient(clientID int64, client Client) {
	c.Lock()
	c.clients[clientID] = client
	c.Unlock()
}

func (s *ChatS) Main() error {
	ctx := &context{s}
	// 定义全局的 Once
	var exitCh = make(chan error)
	var once sync.Once
	exitFunc := func(err error){
		once.Do(func() {
			exitCh <- err
		})
	}

	// 为什么需要 tcpServer 呢？原因是可以统一处理多个不同的 tcp handle，比如 tcp/ip 或者是 http/https 的handle 都可以公用这个，传入相应的上下文就可以了
	// 聊天的服务就传入 chat 的上下文
	// web服务就传入 web 的上下文
	tcpServer := &tcpServer{ctx: ctx}
	s.waitGroup.Wrap(func() {
		exitFunc(protocol.TCPServer(s.tcpListener, tcpServer))
	})

	<- exitCh
	return nil
}