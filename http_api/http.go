package http_api

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type httpServer struct {
	ctx    *context
	router http.Handler
}

func (s *httpServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.router.ServeHTTP(w, req)
}

func Server(ctx *context) *httpServer {
	router := httprouter.New()
	router.HandleMethodNotAllowed = true

	server := &httpServer{
		ctx:    ctx,
		router: router,
	}

	// 登录/注册
	router.Handle("POST", "/v1/sign", Decorate(signHandler, PlainText))
	router.Handle("POST", "/v1/login", Decorate(loginHandler, PlainText))

	// 获取联系人列表
	router.HandlerFunc("GET", "/v1/contacts", MiddlewareHandler(ValidateToken, GetContacts))
	// 获取与其它用户聊天的历史消息
	router.HandlerFunc("GET", "/v1/contact", MiddlewareHandler(ValidateToken, GetContact))
	// 发送消息给指定用户
	router.HandlerFunc("POST", "/v1/create_chat", MiddlewareHandler(ValidateToken, CreateChat))


	// websocket 连接入口
	//router.HandlerFunc("GET", "/ws", MiddlewareHandler(server.ValidateToken, server.WebSocketConnect))
	return server
}

//func (h *httpServer) conf() *conf.Config {
//	return server.Conf
//}
