package http_api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type httpServer struct {
	//ctx    *context
	router http.Handler
}

func (s *httpServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.router.ServeHTTP(w, req)
}

func StartServer() *httpServer {
	router := httprouter.New()
	router.HandleMethodNotAllowed = true

	server := &httpServer{
		router: router,
	}

	// 登录/注册
	router.Handle("POST", "/v1/sign", Decorate(signHandler, PlainText))
	router.Handle("POST", "/v1/login", Decorate(loginHandler, PlainText))

	// 获取联系人列表
	router.HandlerFunc("GET", "/v1/contacts", MiddlewareHandler(ValidateToken, GetContacts))
	// 发送消息
	router.HandlerFunc("POST", "/v1/room/push", MiddlewareHandler(ValidateToken, PushMessage))
	// 创建房间
	router.HandlerFunc("POST", "/v1/room/create", MiddlewareHandler(ValidateToken, CreateRoom))
	// 获取与用户聊天的历史消息
	router.HandlerFunc("GET", "/v1/room/contact", MiddlewareHandler(ValidateToken, GetMessages))

	// websocket 连接入口
	router.HandlerFunc("GET", "/v1/ws", serveWs)
	return server
}
