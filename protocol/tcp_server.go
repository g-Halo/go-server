package protocol

import (
	"github.com/yigger/go-server/logger"
	"net"
	"runtime"
	"strings"
)

// 定义抽象的接口，使用 tcpServer 的都需要实现这个接口
type TCPHandler interface {
	Handle(net.Conn)
}

// 第一个参数是一个 tcp 的连接
// 第二个参数是一个实现了 handle 方法的 struct
func TCPServer(listener net.Listener, handler TCPHandler) error {
	for {
		clientConn, err := listener.Accept()
		if err != nil {
			if nerr, ok := err.(net.Error); ok && nerr.Temporary() {
				runtime.Gosched()
				continue
			}
			// theres no direct way to detect this error because it is not exposed
			if !strings.Contains(err.Error(), "use of closed network connection") {
				logger.Fatalf("listener.Accept() error - %s", err)
			}
			break
		}
		go handler.Handle(clientConn)
	}

	return nil
}
