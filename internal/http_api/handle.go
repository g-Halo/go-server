package http_api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/g-Halo/go-server/pkg/logger"
	"github.com/g-Halo/go-server/pkg/pb"
	"github.com/g-Halo/go-server/pkg/rpc_client"
	"github.com/g-Halo/go-server/pkg/storage"

	"github.com/g-Halo/go-server/internal/logic/model"
	"github.com/g-Halo/go-server/pkg/util"
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

	r, err := rpc_client.AuthClient.Validate(context.Background(), &pb.ValidateReq{TokenStr: tokenString})
	if err != nil {
		return nil, false
	}

	if r.GetCode() == util.Success {
		user := getUser(r.GetUsername())
		if user != nil {
			return user, true
		}
	}
	return nil, false
}

func getUser(username string) *model.User {
	//userResp, err := rpc_client.LogicClient.GetUser(context.Background(), &pb.GetUserReq{
	//	Username:             username,
	//})
	//
	//if err != nil {
	//	logger.Error(err)
	//	return nil
	//}
	//
	//if userResp == nil {
	//	return nil
	//}
	//
	//user := &model.User{
	//	Username:    userResp.User.Username,
	//	NickName:    userResp.User.Nickname,
	//}
	user := storage.GetUser(username)
	return user
}

func getRoom(currentUsername, targetUsername string) *model.Room {
	r, err := rpc_client.LogicClient.FindOrCreateRoom(context.Background(), &pb.FindOrCreateRoomReq{
		CurrentUsername: currentUsername,
		TargetUsername:  targetUsername,
	})
	if err != nil {
		logger.Error(err)
		return nil
	}

	return &model.Room{
		UUID:    r.Uuid,
		Name:    r.Name,
		Members: []string{currentUsername, targetUsername},
		Type:    r.Type,
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

	resp, err := rpc_client.AuthClient.SignUp(context.Background(), &pb.SignUpReq{
		Nickname: req.URL.Query().Get("nickname"),
		Username: req.URL.Query().Get("username"),
		Password: req.URL.Query().Get("password"),
	})
	if err != nil {
		return renderError(err.Error()), nil
	}

	if resp.Code == util.Success {
		return renderSuccess(resp.Msg), nil
	}

	return renderError(resp.Msg), nil
}

// 登录
func loginHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) (interface{}, error) {
	var params loginParams
	err := json.NewDecoder(req.Body).Decode(&params)
	if err != nil {
		return renderError("参数有误"), err
	}

	user := getUser(params.Username)
	if user == nil {
		return renderError("无效的用户"), nil
	}

	r, err := rpc_client.AuthClient.SignIn(context.Background(), &pb.AuthReq{
		Username: params.Username,
		Passowrd: params.Password,
	})

	if err != nil {
		return renderError("未知的错误"), err
	}

	if r.Code == util.Success {
		uJSON := user.ToJson()
		resHs := map[string]interface{}{
			"user":  uJSON,
			"token": r.Data,
		}
		return renderSuccess(resHs), nil
	}

	return renderError(r.Msg), nil
}

// 获取联系人列表
func GetContacts(w http.ResponseWriter, req *http.Request, currentUser *model.User) (interface{}, error) {
	resp, err := rpc_client.LogicClient.GetUsers(context.Background(), &pb.GetUsersReq{})
	if err != nil {
		logger.Error(err)
		return renderError("无法获取联系人"), err
	}

	var users []*model.User
	for _, u := range resp.Users {
		users = append(users, &model.User{
			Username: u.Username,
			NickName: u.Nickname,
		})
	}

	var data []map[string]interface{}
	for _, user := range users {
		if user.Username == currentUser.Username {
			continue
		}

		room := getRoom(currentUser.Username, user.Username)
		if room == nil {
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

	user := getUser(username)
	if user == nil {
		return renderError("无效的用户"), nil
	}
	// 获取他们之间的聊天消息内容
	// 1. 获取他们之间的房间号
	roomID := req.URL.Query().Get("room_id")
	pbRoom, _ := rpc_client.LogicClient.GetRoomById(context.Background(), &pb.GetRoomByIdReq{Uuid: roomID})
	if pbRoom == nil {
		return renderError("Room not found"), nil
	}
	room := &model.Room{
		UUID:      pbRoom.Uuid,
		Name:      pbRoom.Name,
		Members:   pbRoom.Members,
		Type:      pbRoom.Type,
		CreatedAt: time.Unix(pbRoom.CreatedAt, 0),
	}

	// 2. 从 DB 拉取历史聊天记录
	var chatData []map[string]interface{}
	chatData = make([]map[string]interface{}, 0)

	m, _ := rpc_client.LogicClient.GetRoomMessages(context.Background(), &pb.GetRoomMessagesReq{
		Uuid: roomID,
	})
	for _, msg := range m.RoomMessages {
		chatData = append(chatData, map[string]interface{}{
			"recipient":  msg.Recipient,
			"sender":     msg.Sender,
			"body":       msg.Body,
			"created_at": time.Unix(msg.CreatedAt, 0),
			"status":     msg.Status,
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

	user := getUser(params.Username)
	if user == nil {
		return renderError("User Not Found"), nil
	} else if user.Username == currentUser.Username {
		return renderError("TargetUser can not be yourself"), nil
	}

	room := getRoom(currentUser.Username, user.Username)
	if room == nil {
		return renderError("房间创建失败"), nil
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

	_, err = rpc_client.LogicClient.PushMessage(context.Background(), &pb.PushMessageReq{
		ReceiverUsername: params.Username,
		SenderUsername:   currentUser.Username,
		Body:             params.Message,
	})

	if err != nil {
		return renderError("发送失败. Fail -1"), nil
	}

	// FIXME: 此处无法保证一定能返回正确的消息给用户
	return renderSuccess("发送成功"), nil
}
