package comet

import (
	"io"
	"net"

	"github.com/g-Halo/go-server/pkg/logger"
	"github.com/g-Halo/go-server/pkg/protocol"
)

type tcpServer struct {
}

func Init(address string) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		logger.Info("register Comet fail")
	}

	logger.Infof("Comet Init And listening %s", address)
	server := &tcpServer{}
	protocol.TCPServer(listener, server)
}

func (t *tcpServer) Handle(clientConn net.Conn) {
	buf := make([]byte, 8)
	_, err := io.ReadFull(clientConn, buf)
	if err != nil {
		clientConn.Close()
	}

	protocolMagic := string(buf)

	var comet *Comet
	switch protocolMagic {
	case "  g-halo":
		comet = &Comet{}
	default:
		logger.Error("无效的连接类型")
		clientConn.Close()
		return
	}

	logger.Info("Commet 连接成功啦")
	err = comet.IOLoop(clientConn)
	if err != nil {
		return
	}
}
