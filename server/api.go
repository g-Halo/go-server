package server

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	uuid "github.com/satori/go.uuid"
	"github.com/yigger/go-server/logger"
	"github.com/yigger/go-server/model"
	"github.com/yigger/go-server/util"
	"time"

	"net/http"
	"net/url"
)

type loginParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func renderSuccess(data interface{}) (res map[string]interface{}) {
	return map[string]interface{}{
		"status": 200,
		"data": data,
	}
}

func renderError(err string) (res map[string]interface{}) {
	return map[string]interface{}{
		"status": 400,
		"error": err,
	}
}

// 中间件方法，用于校验 jwt 的合法性
func (s *httpServer) ValidateToken(tokenString string) (*model.User, bool) {
	var User model.User

	// FIXME: 测试环境不校验 jwt
	//if s.ctx.chatS.conf.Env == "development" {
	//	// 测试用户的 username
	//	testUsername := "test-1"
	//	user, _ := User.FindByUsername(testUsername)
	//	if user.Username == "" {
	//		logger.Error("validate token: invalid user")
	//	} else {
	//		return user, true
	//	}
	//}
	if tokenString == "" {
		return nil, false
	}

	token, _ := jwt.ParseWithClaims(tokenString, &util.MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.ctx.chatS.conf.SecretKey), nil
	})

	if claims, ok := token.Claims.(*util.MyCustomClaims); ok && token.Valid {
		user, _ := User.FindByUsername(claims.Username)
		return user, true
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
		return renderError("参数校验错误，请检查"), nil
	}

	var user model.User
	if _, err := user.FindByUsername(req.URL.Query().Get("username")); err == nil {
		return renderError("当前用户已被注册"), nil
	}

	params := map[string]interface{}{
		"nickname": req.URL.Query().Get("nickname"),
		"username": req.URL.Query().Get("username"),
		"password": req.URL.Query().Get("password"),
	}
	if err := user.SignUp(params); err != nil {
		return renderError("注册失败"), nil
	}

	return renderSuccess("OK"), nil
}

// 登录
func (s *httpServer) loginHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) (interface{}, error) {
	var params loginParams
	err := json.NewDecoder(req.Body).Decode(&params)
	if err != nil {
		return renderError("参数有误"), err
	}

	if s.conf().No_db() {
		client := s.ctx.chatS.clients[params.Username]

		return client, nil
	} else {
		var User model.User
		if token, err := User.Login(params.Username, params.Password); err != nil {
			return renderError(err.Error()), nil
		} else {
			user, _ := User.FindByUsername(params.Username)
			uJson := user.ToJson()
			res := map[string]interface{}{
				"user": uJson,
				"token": token,
			}
			return renderSuccess(res), nil
		}
	}
}

// 获取联系人列表
func (s *httpServer) GetContacts(w http.ResponseWriter, req *http.Request, currentUser *model.User) (interface{}, error) {
	var User model.User
	users := User.FindAll()
	return renderSuccess(users), nil
}

// 获取与某用户的聊天信息
func (s *httpServer) GetContact(w http.ResponseWriter, req *http.Request, currentUser *model.User) (interface{}, error) {
	username := req.URL.Query().Get("username")
	if username == "" {
		return renderError("无效的用户"), nil
	}

	var User model.User
	user, err := User.FindByUsername(username)
	if err != nil {
		return nil, err
	}

	// 获取他们之间的聊天消息内容
	// 1. 获取他们之间的房间号
	room := currentUser.FindP2PRoom(user.Username)
	if room == nil {
		return renderError("empty"), nil
	}

	// 2. 从 DB 拉取历史聊天记录
	var chatData []map[string]interface{}
	for _, msg := range room.Messages {
		chatData = append(chatData, map[string]interface{}{
			"recipient": msg.Recipient,
			"sender": msg.Sender,
			"body": msg.Body,
			"created_at": msg.CreatedAt,
			"status": "check",
		})
	}

	data := map[string]interface{}{
		"user": user,
		"messages": chatData,
	}

	return renderSuccess(data), nil
}

func (s *httpServer) CreateChat(w http.ResponseWriter, req *http.Request, currentUser *model.User) (interface{}, error) {
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

	targetUser, _ := User.FindByUsername(username)
	if targetUser.Username == "" {
		err := "无效的用户" + username
		return renderError(err), nil
	}

	room := &model.Room{}
	if len(roomId) != 0 {
		// 有房间 Id
		room = chatS.rooms[roomId]
	} else {
		// 先在用户已有的房间列表查找是否有符合的房间
		userArray := []string{currentUser.Username, targetUser.Username}
		p2pRoom := currentUser.FindP2PRoom(targetUser.Username)
		if p2pRoom != nil {
			room = p2pRoom
		}

		if room.UUID == "" {
			chatS.NewRoom()
		}

		// 如果没有的话，则创建 room
		//uid := uuid.NewV4()
		//logger.Info(db_room)
		//if db_room.UUID == "" {
		//	db_room = &model.Room{
		//		UUID:      uid.String(),
		//		Type:      "p2p",
		//		Members:   []string{targetUser.Username, currentUser.Username},
		//		CreatedAt: time.Now(),
		//	}
		//	// insert room to db
		//	logger.Info("start to create room to db")
		//	db_room.Create()
		//	// room add members
		//	db_room.AddMembers(userArray)
		//	// 把 room 加入到 user 的 rooms. user.AddRoom()
		//	currentUser.AddRoom(db_room)
		//}

		// 以上都是 DB 操作，主要是为了备份，防止服务挂掉重启后什么都没了

		// 这里是在内存中创建唯一的房间号，此处与 db 的 room 一一对应
		room = chatS.GetOrCreateByRoom(db_room.UUID)
	}

	// 获取到当前正在连接的用户 Client
	client := chatS.clients[currentUser.Username]

	// 创建 message 信息
	var Message model.Message
	message := NewMessage(client, content)
	db_message := Message.Create(currentUser, targetUser, content)

	// 往房间里面扔消息
	room.AddMessage(message)
	db_room.AddMessage(db_message)

	// 因为对方已经订阅了 room 的频道，所以对方应该是可以收到的
	return renderSuccess("OK"), nil
}
