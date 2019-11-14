package logic

import (
	"errors"
	"sync"

	"github.com/g-Halo/go-server/model"
	"github.com/g-Halo/go-server/storage"
)

type userLogic struct {
	mutex *sync.Mutex
}

var UserLogic = &userLogic{mutex: &sync.Mutex{}}

func (logic *userLogic) GetUsers() []map[string]interface{} {
	var data []map[string]interface{}
	for _, user := range storage.GetUsers() {
		data = append(data, map[string]interface{}{
			"username":   user.Username,
			"nickname":   user.NickName,
			"created_at": user.CreatedAt,
			"unread":     "uncheck",
			"last_message": map[string]string{
				"body":       "hello",
				"created_at": "2019-10-01 12:00:00",
			},
		})
	}
	return data
}

func (logic *userLogic) Login(username, password string) (*model.User, string, error) {
	var u model.User
	user := logic.FindByUsername(username)
	if user.Username == "" {
		return nil, "", errors.New("User Not Found")
	}

	token, err := u.Login(user, password)
	if err != nil {
		return nil, "", err
	}

	return user, token, err
}

func (logic *userLogic) SignUp(params map[string]interface{}) *model.User {
	var u model.User
	if params["username"] == "" {
		return nil
	}

	username := params["username"].(string)
	user := logic.FindByUsername(username)
	if user != nil {
		return nil
	}

	user = u.New(params)
	storage.AddUser(user)
	return user
}

func (logic *userLogic) FindByUsername(username string) *model.User {
	for _, u := range storage.GetUsers() {
		if u != nil && u.Username == username {
			return u
		}
	}
	return nil
}