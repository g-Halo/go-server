package server

import (
	"github.com/julienschmidt/httprouter"
	"github.com/yigger/go-server/http_api"
	"net/http"
)

type httpServer struct {
	ctx    *context
	router http.Handler
}

func newHTTPServer(ctx *context) *httpServer {
	router := httprouter.New()
	router.HandleMethodNotAllowed = true

	server := &httpServer{
		ctx:    ctx,
		router: router,
	}

	// 登录/注册
	router.Handle("POST", "/sign", http_api.Decorate(server.signHandler, http_api.PlainText))
	router.Handle("POST", "/login", http_api.Decorate(server.loginHandler, http_api.PlainText))

	// 获取联系人列表
	router.HandlerFunc("GET", "/v1/contacts", http_api.MiddlewareHandler(server.ValidateToken, server.GetContacts))
	// 获取与其它用户聊天的历史消息
	router.HandlerFunc("GET", "/v1/contact", http_api.MiddlewareHandler(server.ValidateToken, server.GetContact))
	// 发送消息给指定用户
	router.HandlerFunc("POST", "/v1/create_chat", http_api.MiddlewareHandler(server.ValidateToken, server.CreateChat))


	// websocket 连接入口
	router.HandlerFunc("GET", "/ws", http_api.MiddlewareHandler(server.ValidateToken, server.WebSocketConnect))
	return server
}
