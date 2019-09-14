package server

import (
	"github.com/julienschmidt/httprouter"
	"github.com/yigger/go-server/http_api"
	"net/http"
)

type httpServer struct {
	ctx *context
	router http.Handler
}

func newHTTPServer(ctx *context) *httpServer {
	router := httprouter.New()
	router.HandleMethodNotAllowed = true

	server := &httpServer{
		ctx: ctx,
		router: router,
	}

	// 静态资源
	router.ServeFiles("/public/*filepath", http.Dir("/Users/yigger/Projects/go-server/public"))
	// 主入口 html
	router.Handle("GET", "/", http_api.RenderTemplate("home"))


	// --------- 下面的都是 API 请求接口 --------------

	// 登录/注册
	router.Handle("POST", "/sign", http_api.Decorate(server.signHandler, http_api.PlainText))
	router.Handle("POST", "/login", http_api.Decorate(server.loginHandler, http_api.PlainText))

	// 含有中间件的 Api 方法，需要进行校验 token
	router.HandlerFunc("POST", "/v1/create_room", http_api.MiddlewareHandler(server.ValidateToken, server.CreateRoom))
	router.HandlerFunc("GET", "/v1/contacts", http_api.MiddlewareHandler(server.ValidateToken, server.GetContacts))

	// websocket 连接入口
	router.HandlerFunc("GET", "/ws", http_api.MiddlewareHandler(server.ValidateToken, server.WebSocketConnect))

	return server
}
