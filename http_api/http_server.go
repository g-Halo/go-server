package http_api

import (
	"github.com/g-Halo/go-server/logger"
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
