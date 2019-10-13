package logic

import (
	"errors"
	"github.com/g-Halo/go-server/logger"
	"github.com/g-Halo/go-server/model"
	"github.com/g-Halo/go-server/server"
	"sync"
)

type userLogic struct {
	mutex *sync.Mutex
	server *server.ChatS
}

var UserLogic = &userLogic{mutex: &sync.Mutex{}}

func (logic *userLogic) Register(server *server.ChatS) {
	logic.server = server
}

func (logic *userLogic) Login(username, password string) (*model.User, string, error) {
	var u model.User
	user := logic.FindByUsername(username)
	if user == nil {
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

	if user := logic.FindByUsername(params["username"].(string)); user != nil {
		return user
	}

	user := u.New(params)
	if logic.server.Conf.No_db() {
		return user
	}

	u.Create(user)
	return user
}

func (logic *userLogic) FindByUsername(username string) *model.User {
	server := logic.server
	if server.Conf.No_db() {
		for _, user := range server.Users {
			if user.Username == username {
				return user
			}
		}
	} else {
		// find in db
		var u model.User
		user, err := u.FindByUsername(username)
		if err != nil {
			logger.Error(err)
		} else {
			return user
		}
	}
	return nil
}
