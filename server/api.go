package server

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	"github.com/yigger/go-server/model"
	"github.com/yigger/go-server/util"
	"go.mongodb.org/mongo-driver/bson"

	"net/http"
	"net/url"
)

// 中间件方法，用于校验 jwt 的合法性
func (s *httpServer) ValidateToken(tokenString string) (*model.User, bool) {
	var User model.User

	// 测试环境不校验 jwt
	if s.ctx.chatS.conf.Env == "development" {
		// 测试用户的 username
		testUsername := "test"
		client := s.ctx.chatS.mongoClient
		user, _ := User.FindByUsername(client, bson.M{"username": testUsername})
		return &user, true
	}

	token, _ := jwt.ParseWithClaims(tokenString, &util.MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.ctx.chatS.conf.SecretKey), nil
	})

	if claims, ok := token.Claims.(*util.MyCustomClaims); ok && token.Valid {

		client := s.ctx.chatS.mongoClient
		user, _ := User.FindByUsername(client, bson.M{"username": claims.Username})
		return &user, true
	} else {
		return nil, false
	}
}

func validateParams(params url.Values, keys []string) bool {
	for _, key := range keys {
		str := params.Get(key)
		if len(str) < 2 {
			return false
		}
	}
	return true
}

// 注册
// localhost:5001/sign?username=yigger&password=123456&nickname=yigger
func (s *httpServer) signHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) (interface{}, error) {
	// 校验字符串
	if !validateParams(req.URL.Query(), []string{"username", "password", "nickname"}) {
		return "参数校验错误，请检查", nil
	}

	var user model.User
	client := s.ctx.chatS.mongoClient
	if _, err := user.FindByUsername(client, bson.M{"username": req.URL.Query().Get("username")}); err == nil {
		return "当前用户已被注册", nil
	}

	params := map[string]interface{}{
		"nickname": req.URL.Query().Get("nickname"),
		"username": req.URL.Query().Get("username"),
		"password": req.URL.Query().Get("password"),
	}
	if err := user.SignUp(client, params); err != nil {
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

func (s *httpServer) CreateRoom(w http.ResponseWriter, r *http.Request, user *model.User) {

}

func (s *httpServer) GetContacts(w http.ResponseWriter, req *http.Request, user *model.User) {
	var User model.User
	client := s.ctx.chatS.mongoClient
	users := User.FindAll(client)
	data, _ := json.Marshal(users)
	w.Header().Set("content-type", "application/json")
	w.Write(data)
}