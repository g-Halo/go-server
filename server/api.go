package server

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	"github.com/yigger/go-server/model"
	"github.com/yigger/go-server/util"
	"go.mongodb.org/mongo-driver/bson"

	"net/http"
	"net/url"
)

func validateParams(params map[string][]string, keys []string) bool {
	for _, key := range keys {
		if _, ok := params[key]; !ok || len(params[key][0]) < 2 {
			return false
		}
	}
	return true
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
		//s.ctx.chatS
		return &user, true
	} else {
		return nil, false
	}
}

func (s *httpServer) CreateRoom(w http.ResponseWriter, r *http.Request, user *model.User) {

}

func (s *httpServer) GetContacts(w http.ResponseWriter, req *http.Request, user *model.User) {
	users := make(map[string]interface{})
	users["code"] = 200
	users["data"] = "list"
}