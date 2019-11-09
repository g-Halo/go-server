package main

import (
	"github.com/g-Halo/go-server/model"
	"sync"
)

type userLogic struct {
	mutex *sync.Mutex
}

//var UserLogic = &userLogic{mutex: &sync.Mutex{}}

//func (logic *userLogic) GetUsers() []map[string]interface{} {
//	var data []map[string]interface{}
//	for _, user := range logic.server.Users {
//		data = append(data, map[string]interface{}{
//			"username": user.Username,
//			"nickname": user.NickName,
//			"created_at": user.CreatedAt,
//			"unread": "uncheck",
//			"last_message": map[string]string{
//				"body": "hello",
//				"created_at": "2019-10-01 12:00:00",
//			},
//		})
//	}
//	return data
//}
//
//func (logic *userLogic) Login(username, password string) (*model.User, string, error) {
//	var u model.User
//	user := logic.FindByUsername(username)
//	if user == nil {
//		return nil, "", errors.New("User Not Found")
//	}
//
//	token, err := u.Login(user, password)
//	if err != nil {
//		return nil, "", err
//	}
//
//	return user, token, err
//}
//
func (logic *userLogic) SignUp(params map[string]interface{}) *model.User {
	var u model.User
	if params["username"] == "" {
		return nil
	}

	user := &model.User{}
	username := params["username"].(string)
	if err := logic.FindByUsername(&username, user); err != nil {
		return nil
	}

	if user.Username != "" {
		return user
	}


	user = u.New(params)
	return user
}

func (logic *userLogic) FindByUsername(username *string, user *model.User) error {
	for _, u := range sto.Users {
		if u != nil && u.Username == *username {
			*user = *u
			return nil
		}
	}
	return nil
}
