package server

import (
	"github.com/julienschmidt/httprouter"
	"github.com/yigger/go-server/http_api"
	"github.com/yigger/go-server/model"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"net/url"
	//ctx "context"
)

type httpServer struct {
	ctx *context
	router http.Handler
}

func validateParams(params map[string][]string, keys []string) bool {
	for _, key := range keys {
		if _, ok := params[key]; !ok || len(params[key][0]) < 2 {
			return false
		}
	}

	return true
}

func newHTTPServer(ctx *context) *httpServer {
	router := httprouter.New()
	router.HandleMethodNotAllowed = true

	server := &httpServer{
		ctx: ctx,
		router: router,
	}

	router.Handle("GET", "/", http_api.Decorate(server.rootHandle, http_api.PlainText))
	router.Handle("POST", "/sign", http_api.Decorate(server.signHandler, http_api.PlainText))
	router.Handle("POST", "/login", http_api.Decorate(server.loginHandler, http_api.PlainText))

	return server
}

func (s *httpServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.router.ServeHTTP(w, req)
}

// 注册
// localhost:5001/sign?username=yigger&password=123456&nickname=yigger
func (s *httpServer) signHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) (interface{}, error) {
	// 校验字符串
	reqParams, _ := url.ParseQuery(req.URL.RawQuery)
	if !validateParams(reqParams, []string{"username", "password", "nickname"}) {
		return "参数校验错误，请检查", nil
	}

	var user model.User
	client := s.ctx.chatS.mongoClient
	if _, err := user.FindByUsername(client, bson.M{"username": reqParams["username"][0]}); err == nil {
		return "当前用户已被注册", nil
	}

	if err := user.SignUp(client, reqParams["nickname"][0], reqParams["username"][0], reqParams["password"][0] ); err != nil {
		return "注册失败", nil
	}

	return "OK", nil
}

// 登录
func (s *httpServer) loginHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) (interface{}, error) {
	reqParams, _ := url.ParseQuery(req.URL.RawQuery)
	if !validateParams(reqParams, []string{"username", "password"}) {
		return "参数校验错误，请检查", nil
	}

	var user model.User
	client := s.ctx.chatS.mongoClient
	token, err := user.Login(client, reqParams["username"][0], reqParams["password"][0])

	return token, err
}

func (s *httpServer) rootHandle(w http.ResponseWriter, req *http.Request, ps httprouter.Params) (interface{}, error) {

	return "OK", nil
}
