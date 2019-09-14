package server

import (
	ctx "context"
	"fmt"
	"github.com/yigger/go-server/conf"
	"github.com/yigger/go-server/http_api"
	"github.com/yigger/go-server/protocol"
	"github.com/yigger/go-server/util"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net"
	"net/http"
	"sync"
	"time"
)

type ChatS struct {
	// 64bit atomic vars need to be first for proper alignment on 32bit platforms
	clientIDSequence 		int64
	conf 				 	*conf.Config

	clients 			 	map[string]*client
	rooms				 	map[string]*room

	sync.RWMutex
	waitGroup            	util.WaitGroupWrapper

	tcpListener   		 	net.Listener
	httpListener		 	net.Listener

	mongoClient 			*mongo.Client
	startTime			 	time.Time
}

func NewChatS(conf *conf.Config) (chat *ChatS, err error) {
	chat = &ChatS{
		startTime: time.Now(),
		clients: make(map[string]*client),
		rooms: make(map[string]*room),
		conf: conf,
	}
	fmt.Printf("Start listening Tcp %s ...\n", conf.TcpAddress)
	chat.tcpListener, err = net.Listen("tcp", conf.TcpAddress)

	// 初始化 http servers
	fmt.Printf("Start listening Http %s ...\n", conf.HttpAddress)
	chat.httpListener, err = net.Listen("tcp", conf.HttpAddress)

	// 初始化 mongodb 的 client
	clientOptions := options.Client().ApplyURI(conf.MongoDbAddress)
	mongoClient, err := mongo.Connect(ctx.TODO(), clientOptions)
	chat.mongoClient = mongoClient

	return
}

func (c *ChatS) AddClient(clientID string, client *client) {
	c.Lock()
	c.clients[clientID] = client
	c.Unlock()
}

func (s *httpServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.router.ServeHTTP(w, req)
}

func (s *ChatS) Main() error {
	contxt := &context{s}
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
	tcpServer := &tcpServer{ctx: contxt}
	s.waitGroup.Wrap(func() {
		exitFunc(protocol.TCPServer(s.tcpListener, tcpServer))
	})

	// 初始化 http api 服务
	httpServer := newHTTPServer(contxt)
	s.waitGroup.Wrap(func() {
		exitFunc(http_api.Serve(s.httpListener, httpServer, "HTTP"))
	})

	<- exitCh
	return nil
}
