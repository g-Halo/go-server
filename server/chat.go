package server

import (
	"github.com/g-Halo/go-server/conf"
	"github.com/g-Halo/go-server/http_api"
	"github.com/g-Halo/go-server/logger"
	"github.com/g-Halo/go-server/model"
	"github.com/g-Halo/go-server/protocol"
	"github.com/g-Halo/go-server/util"
	"go.mongodb.org/mongo-driver/mongo"
	//"go.mongodb.org/mongo-driver/mongo/options"
	"net"
	"sync"
	"time"
)

type ChatS struct {
	// 64bit atomic vars need to be first for proper alignment on 32bit platforms
	ClientIDSequence 		int64
	Conf 				 	*conf.Config

	Clients 			 	map[string]*Client
	Rooms				 	[]*model.Room
	Users					[]*model.User

	tcpListener   		 	net.Listener
	httpListener		 	net.Listener

	mongoClient 			*mongo.Client
	startTime			 	time.Time

	sync.RWMutex
	waitGroup            	util.WaitGroupWrapper
}

type LogicInterface interface {
	Register(*ChatS)
}

func (s *ChatS) RegisterLogic(i LogicInterface) {
	i.Register(s)
}

func NewChatS(config *conf.Config) (chat *ChatS, err error) {
	chat = &ChatS{
		startTime: time.Now(),
		Clients: make(map[string]*Client),
		Conf: config,
	}
	return
}

// Initialize http server
func (c *ChatS) CreateHttpListen() net.Listener {
	var err error
	c.httpListener, err = net.Listen("tcp", c.Conf.HttpAddress)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Infof("Start to listening Http %s ...\n", c.Conf.HttpAddress)
	return c.httpListener
}

// Initialize tcp server
func (c *ChatS) CreateTcpListen() net.Listener {
	var err error
	c.tcpListener, err = net.Listen("tcp", c.Conf.TcpAddress)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Infof("Start to listening Tcp %s ...\n", c.Conf.TcpAddress)
	return c.tcpListener
}

// Initialize mongodb
func (c *ChatS) CreateDatabase() {
	config := c.Conf
	if config.No_db() {
		//clientOptions := options.Client().ApplyURI(config.MongoDbAddress)
		//mongoClient, err := mongo.Connect(ctx.TODO(), clientOptions)
		//if err != nil {
		//	logger.Fatal("mongo client error:", err)
		//}
		//c.mongoClient = mongoClient
	}
}

func (c *ChatS) AddClient(clientID string, client *Client) {
	c.Lock()
	defer c.Unlock()
	c.Clients[clientID] = client
}

func (c *ChatS) findRoomById(uuid string) *model.Room {
	for _, v := range c.Rooms {
		if v.UUID == uuid {
			return v
		}
	}
	return nil
}

// 拉起 tcp 服务和 http 服务
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

	tcpServer := &tcpServer{ctx: contxt}
	s.waitGroup.Wrap(func() {
		exitFunc(protocol.TCPServer(s.tcpListener, tcpServer))
	})

	context := http_api.NewContext(s)
	httpServer := http_api.Server(context)
	s.waitGroup.Wrap(func() {
		exitFunc(http_api.Serve(s.httpListener, httpServer, "HTTP"))
	})

	<- exitCh
	return nil
}
