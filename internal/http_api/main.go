package http_api

import (
	"github.com/g-Halo/go-server/conf"
	"github.com/g-Halo/go-server/pkg/logger"
	"net"
	"net/http"
	"strings"
)

func Serve(listener net.Listener, handler http.Handler, proto string) error {
	server := &http.Server{
		Handler:  handler,
	}
	err := server.Serve(listener)
	if err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
		logger.Fatal("http.Serve() error - %s", err)
	}
	return nil
}

func Main() {
	httpListener, err := net.Listen("tcp", conf.Conf.HttpApiAddress)
	if err != nil {
		logger.Fatal("error")
	}

	logger.Infof("start Listen web api in %s", conf.Conf.HttpApiAddress)
	Serve(httpListener, StartServer(), "HTTP")
}
