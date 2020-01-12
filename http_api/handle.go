package http_api

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/g-Halo/go-server/storage"

	"github.com/g-Halo/go-server/logger"
	"github.com/g-Halo/go-server/model"
	"github.com/g-Halo/go-server/rpc/instance"
	"github.com/g-Halo/go-server/util"
	"github.com/julienschmidt/httprouter"
)

type loginParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateRoomParams struct {
	Username string `json:"username"`
}

type RoomPushParams struct {
	RoomId   string `json:"room_id"`
	Username string `json:"username"`
	Message  string `json:"message"`
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
	var users []*model.User
	if err := client.Call("Logic.GetUsers", "", &users); err != nil {
		return renderError("Get Contacts Fail -1"), err
	}

	var data []map[string]interface{}
	for _, user := range users {
		if user.Username == currentUser.Username {
			continue
		}

		var room *model.Room
		if err := client.Call("Logic.FindOrCreate", []string{currentUser.Username, user.Username}, &room); err != nil {
			logger.Error(err.Error())
			continue
		}

		data = append(data, map[string]interface{}{
			"username":   user.Username,
			"nickname":   user.NickName,
			"room_id":    room.UUID,
			"created_at": user.CreatedAt,
			"unread":     "uncheck",
			"last_message": map[string]string{
				"body":       "hello",
				"created_at": "2019-10-01 12:00:00",
			},
		})
	}

	return renderSuccess(data), nil
}

// 获取与某用户的聊天信息
func GetMessages(w http.ResponseWriter, req *http.Request, currentUser *model.User) (interface{}, error) {
	username := req.URL.Query().Get("username")
	if username == "" {
		return renderError("无效的用户"), nil
	}

	logicClient := instance.LogicRPC()
	var user *model.User
	if err := logicClient.Call("Logic.FindByUsername", &username, &user); err != nil {
		return renderError("User not found"), err
	}

	// 获取他们之间的聊天消息内容
	// 1. 获取他们之间的房间号
	roomID := req.URL.Query().Get("room_id")
	room := storage.GetRoom(roomID)
	if room == nil {
		return renderError("Room not found"), nil
	}

	// 2. 从 DB 拉取历史聊天记录
	var chatData []map[string]interface{}
	for _, msg := range room.Messages {
		chatData = append(chatData, map[string]interface{}{
			"recipient":  msg.Recipient,
			"sender":     msg.Sender,
			"body":       msg.Body,
			"created_at": msg.CreatedAt,
			"status":     "check",
		})
	}

	data := map[string]interface{}{
		"room":     room.UUID,
		"user":     user.ToJson(),
		"messages": chatData,
	}

	return renderSuccess(data), nil
}

// 创建房间接口
// POST /v1/room/create
// type: JSON
// params: {
// 	username: "xxx"
// }
func CreateRoom(w http.ResponseWriter, req *http.Request, currentUser *model.User) (interface{}, error) {
	var params CreateRoomParams
	err := json.NewDecoder(req.Body).Decode(&params)
	if err != nil {
		return renderError("Params Error"), err
	}
	if params.Username == "" {
		return renderError("User Not Found"), nil
	}

	logicClient := instance.LogicRPC()
	var user *model.User
	logicClient.Call("Logic.FindByUsername", &params.Username, &user)
	if user == nil {
		return renderError("User Not Found"), nil
	} else if user.Username == currentUser.Username {
		return renderError("TargetUser can not be yourself"), nil
	}

	var room *model.Room
	if err := logicClient.Call("Logic.FindOrCreate", []string{currentUser.Username, user.Username}, &room); err != nil {
		return renderError("创建失败， -1"), nil
	}

	if room == nil {
		return renderError("创建失败， -2"), nil
	}

	// User 结构下的 Room 仅代表当前聊天窗口的 Room
	// FIXME: 检测重复添加
	currentUser.Rooms = append(currentUser.Rooms, room)
	user.Rooms = append(user.Rooms, room)

	return renderSuccess("创建成功"), nil
}

// 向房间发送新消息
// POST /v1/room/push
// type: JSON
// params: {
// 	username: "xxx",
//  room_id: "xxx",
//  message: "xxx"
// }
func PushMessage(w http.ResponseWriter, req *http.Request, currentUser *model.User) (interface{}, error) {
	var params RoomPushParams
	err := json.NewDecoder(req.Body).Decode(&params)
	if err != nil {
		return renderError("Params Error"), err
	}
	if params.Username == "" && params.RoomId == "" {
		return renderError("User Not Found"), nil
	}

	logicClient := instance.LogicRPC()
	logicClient.Call("Logic.Push", []string{currentUser.Username, params.Username, params.Message}, &err)

	if err != nil {
		return renderError(err.Error()), nil
	}

	// FIXME: 此处无法保证一定能返回正确的消息给用户
	return renderSuccess("发送成功"), nil
}
