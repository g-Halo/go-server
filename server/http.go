package server

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	"github.com/yigger/go-server/http_api"
	"github.com/yigger/go-server/model"
	"github.com/yigger/go-server/util"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/gorilla/websocket"
	"io"
	"net/http"
	"net/url"
	"time"

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

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func newHTTPServer(ctx *context) *httpServer {
	router := httprouter.New()
	router.HandleMethodNotAllowed = true

	server := &httpServer{
		ctx: ctx,
		router: router,
	}

	router.ServeFiles("/public/*filepath", http.Dir("/Users/yigger/Projects/go-server/public"))
	router.Handle("GET", "/", http_api.RenderTemplate("home"))
	router.Handle("POST", "/sign", http_api.Decorate(server.signHandler, http_api.PlainText))
	router.Handle("POST", "/login", http_api.Decorate(server.loginHandler, http_api.PlainText))
	router.HandlerFunc("POST", "/v1/info", http_api.MiddlewareHandler(server.ValidateToken, server.Info))
	router.HandlerFunc("POST", "/v1/create_room", http_api.MiddlewareHandler(server.ValidateToken, server.CreateRoom))
	router.HandlerFunc("GET", "/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, _ := upgrader.Upgrade(w, r, nil)
		go server.readPump(conn)
		//go server.writePump(conn)
	})
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

func (s *httpServer) ValidateToken(tokenString string) (*model.User, bool) {
	token, _ := jwt.ParseWithClaims(tokenString, &util.MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.ctx.chatS.conf.SecretKey), nil
	})

	if claims, ok := token.Claims.(*util.MyCustomClaims); ok && token.Valid {
		var User model.User
		client := s.ctx.chatS.mongoClient
		user, _ := User.FindByUsername(client, bson.M{"username": claims.Username})
		return &user, true
	} else {
		return nil, false
	}
}

func (s *httpServer) CreateRoom(w http.ResponseWriter, r *http.Request, p string) {
	io.WriteString(w, "test")
}

func (s *httpServer) Info(w http.ResponseWriter, r *http.Request, p string) {
	chatS := s.ctx.chatS
	roomNums := len(chatS.rooms)
	clientNums := len(chatS.clients)
	str := fmt.Sprintf("room: %d， client: %d", roomNums, clientNums)
	io.WriteString(w, str)
}

func (s *httpServer) readPump(conn *websocket.Conn) {
	for {
		_, message, _ := conn.ReadMessage()
		if message == nil {
			return
		}

		fmt.Println(string(message))
	}
}


func (s *httpServer) writePump(conn *websocket.Conn) {
	for {
		w, _ := conn.NextWriter(websocket.TextMessage)
		w.Write([]byte("test"))
		time.Sleep(5e9)
	}
}
