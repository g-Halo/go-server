package server

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	uuid "github.com/satori/go.uuid"
	"github.com/yigger/go-server/model"
	"github.com/yigger/go-server/util"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"time"

	"net/http"
	"net/url"
)

// 中间件方法，用于校验 jwt 的合法性
func (s *httpServer) ValidateToken(tokenString string) (*model.User, bool) {
	var User model.User

	// 测试环境不校验 jwt
	if s.ctx.chatS.conf.Env == "development" {
		// 测试用户的 username
		testUsername := "test-1"
		client := s.ctx.chatS.mongoClient
		user, _ := User.FindByUsername(client, bson.M{"username": testUsername})
		if user.Username == "" {
			log.Fatal("validate token: invalid user")
		} else {
			return &user, true
		}
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
	token, err := user.Login(client, reqParams.Get("username"), reqParams.Get("password"))

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
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(data)
}

func (s *httpServer) CreateChat(w http.ResponseWriter, req *http.Request, user *model.User) {
	// 1. 从 form-data 获取 username，如果有 room_id 则也一起传过来
	// TODO: 先校验收发端是否好友
	// 2. 校验聊天用户的有效性
	// 3. 查找 room 是否存在，如何查找？ 1. 有 room_id 直接找 	2. 遍历用户目前已有的 rooms, room.members.len 和 room.members 是否只包含他们两个人
	// 4. 如果 room 不存在，则新建一间 room，存储双方的信息
	// 5. 往 room 里面塞 message，并发送信号通知对方端
	var User model.User
	_ = req.ParseForm()
	username := req.PostForm.Get("username")
	content := req.PostForm.Get("content")
	roomId := req.PostForm.Get("room_id")

	chatS := s.ctx.chatS

	mgClient := chatS.mongoClient
	queryFilter := bson.M{"username": username}
	targetUser, _ := User.FindByUsername(mgClient, queryFilter)
	if targetUser.Username == "" {
		w.Write([]byte("无效的用户" + username))
		return
	}

	room := &room{}
	if len(roomId) != 0 {
		// 有房间 Id
		room = chatS.rooms[roomId]
	} else {
		// 先在用户已有的房间列表查找是否有符合的房间
		userArray := []string{user.Username, targetUser.Username}
		mroom := new(model.Room)
		for _, r := range user.Rooms {
			if r.Type == "p2p" && len(r.Members) == 2 {
				if (r.Members[0] == userArray[0] || r.Members[0] == userArray[1]) || (r.Members[1] == userArray[0] || r.Members[1] == userArray[1]) {
					log.Println("find the exist room")
					mroom = r
					break
				}
			}
		}

		// 如果没有的话，则创建 room
		uid := uuid.NewV4()
		if mroom.UUID == "" {
			mroom = &model.Room{
				UUID: uid.String(),
				Type: "p2p",
				Members: []string{targetUser.Username, user.Username},
				CreatedAt: time.Now(),
			}
			// insert room to db
			mroom.Create(mgClient)
			// room add members
			mroom.AddMembers(mgClient, userArray)
			// 把 room 加入到 user 的 rooms. user.AddRoom()
			user.AddRoom(mgClient, mroom)
		}

		// 以上都是 DB 操作，主要是为了备份，防止服务挂掉重启后什么都没了

		// 这里是在内存中创建唯一的房间号，此处与 db 的 room 一一对应
		room = chatS.GetOrCreateByRoom(mroom.UUID)
	}

	// 获取到当前正在连接的用户 Client
	client := chatS.clients[user.Username]

	// 创建 message 信息
	message := NewMessage(client, content)

	// 往房间里面扔消息
	room.AddMessage(message)

	// 因为对方已经订阅了 room 的频道，所以对方应该是可以收到的

	w.Write([]byte("OK"))
}