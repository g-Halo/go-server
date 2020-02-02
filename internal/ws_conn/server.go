package ws_conn

import (
	"github.com/g-Halo/go-server/pkg/logger"
	"net/http"
)

func StartServer() {
	http.HandleFunc("/v1/ws", wsHandler)
	logger.Info("websocket init connect")
	err := http.ListenAndServe(":7771", nil)
	if err != nil {
		panic(err)
	}
}