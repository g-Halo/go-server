package http_api

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/g-Halo/go-server/model"
	"github.com/g-Halo/go-server/rpc/instance"
	"github.com/g-Halo/go-server/util"
	"github.com/julienschmidt/httprouter"
)

type loginParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func renderSuccess(data interface{}) (res map[string]interface{}) {
	return map[string]interface{}{
		"status": 200,
		"data":   data,
	}
}

func renderError(err string) (res map[string]interface{}) {
	return map[string]interface{}{
		"status": 400,
		"error":  err,
	}
}

// 中间件方法，用于校验 jwt 的合法性
func ValidateToken(tokenString string) (*model.User, bool) {
	if tokenString == "" {
		return nil, false
	}

	var user *model.User
	authClient := instance.AuthRPC()
	authClient.Call("Token.Validate", &tokenString, &user)
	if user != nil {
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
func signHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) (interface{}, error) {
	// 校验字符串
	if !validateParams(req.URL.Query(), []string{"username", "password", "nickname"}) {
		return renderError("参数校验错误，请检查"), nil
	}

	params := map[string]interface{}{
		"nickname": req.URL.Query().Get("nickname"),
		"username": req.URL.Query().Get("username"),
		"password": req.URL.Query().Get("password"),
	}

	var user model.User
	if err := instance.LogicRPC().Call("Logic.SignUp", params, &user); err != nil {
		return renderError("注册失败"), nil
	}
	return renderSuccess("OK"), nil
}

// 登录
func loginHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) (interface{}, error) {
	var params loginParams
	var res util.Response
	logicClient := instance.LogicRPC()
	authClient := instance.AuthRPC()

	err := json.NewDecoder(req.Body).Decode(&params)
	if err != nil {
		return renderError("参数有误"), err
	}

	var user *model.User
	if err := logicClient.Call("Logic.FindByUsername", &params.Username, &user); err != nil {
		return renderError("Login Fail -2"), err
	}

	if err := authClient.Call("Token.Create", params, &res); err != nil {
		return renderError("Login Fail -1"), err
	}

	if res.Code != util.Success {
		return renderError(res.Msg), err
	}

	uJSON := user.ToJson()
	resHs := map[string]interface{}{
		"user":  uJSON,
		"token": res.Data,
	}
	return renderSuccess(resHs), nil
}

// 获取联系人列表
func GetContacts(w http.ResponseWriter, req *http.Request, currentUser *model.User) (interface{}, error) {
	client := instance.LogicRPC()
	var users map[string]interface{}
	client.Call("Logic.GetUsers", nil, &users)

	return renderSuccess(users), nil
}

// 获取与某用户的聊天信息
//func (s *httpServer) GetContact(w http.ResponseWriter, req *http.Request, currentUser *model.User) (interface{}, error) {
//	username := req.URL.Query().Get("username")
//	if username == "" {
//		return renderError("无效的用户"), nil
//	}
//
//	var User model.User
//	user, err := User.FindByUsername(username)
//	if err != nil {
//		return nil, err
//	}
//
//	// 获取他们之间的聊天消息内容
//	// 1. 获取他们之间的房间号
//	room := currentUser.FindP2PRoom(user.Username)
//	if room == nil {
//		return renderError("empty"), nil
//	}
//
//	// 2. 从 DB 拉取历史聊天记录
//	var chatData []map[string]interface{}
//	for _, msg := range room.Messages {
//		chatData = append(chatData, map[string]interface{}{
//			"recipient": msg.Recipient,
//			"sender": msg.Sender,
//			"body": msg.Body,
//			"created_at": msg.CreatedAt,
//			"status": "check",
//		})
//	}
//
//	data := map[string]interface{}{
//		"user": user,
//		"messages": chatData,
//	}
//
//	return renderSuccess(data), nil
//}
//
//func (s *httpServer) CreateChat(w http.ResponseWriter, req *http.Request, currentUser *model.User) (interface{}, error) {
//	// 1. 从 form-data 获取 username，如果有 room_id 则也一起传过来
//	// TODO: 先校验收发端是否好友
//	// 2. 校验聊天用户的有效性
//	// 3. 查找 room 是否存在，如何查找？ 1. 有 room_id 直接找 	2. 遍历用户目前已有的 rooms, room.members.len 和 room.members 是否只包含他们两个人
//	// 4. 如果 room 不存在，则新建一间 room，存储双方的信息
//	// 5. 往 room 里面塞 message，并发送信号通知对方端
//	_ = req.ParseForm()
//	username := req.PostForm.Get("username")
//	content := req.PostForm.Get("content")
//	roomId := req.PostForm.Get("room_id")
//
//	//chatS := s.ctx.chatS
//
//	var User model.User
//	targetUser, _ := User.FindByUsername(username)
//	if targetUser.Username == "" {
//		err := "无效的用户" + username
//		return renderError(err), nil
//	}
//
//	room := &model.Room{}
//	if len(roomId) != 0 {
//		// 有房间 Id
//		//room = server.findRoomById(roomId)
//	} else {
//		// 先在用户已有的房间列表查找是否有符合的房间
//		userArray := []string{currentUser.Username, targetUser.Username}
//		p2pRoom := currentUser.FindP2PRoom(targetUser.Username)
//		if p2pRoom != nil {
//			room = p2pRoom
//		} else {
//			uid := uuid.NewV4()
//			room = room.New(uid.String(), userArray)
//			// TODO: 把用户添加到房间里面去
//
//			//server.Rooms = append(server.Rooms, room)
//		}
//	}
//
//	// 获取到当前正在连接的用户 Client
//	//client := chatS.clients[currentUser.Username]
//
//	// 创建 message 信息
//	var Message model.Message
//	message := Message.Create(currentUser, targetUser, content)
//
//	// 往房间里面扔消息
//	room.AddMessage(message)
//
//	// 因为对方已经订阅了 room 的频道，所以对方应该是可以收到的
//	return renderSuccess("OK"), nil
//}
