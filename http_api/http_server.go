package http_api

import (
	"fmt"
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
		return fmt.Errorf("http.Serve() error - %s", err)
	}


	return nil
}
