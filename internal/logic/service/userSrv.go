package service

import (
	"errors"
	"github.com/g-Halo/go-server/internal/logic/model"
	"github.com/g-Halo/go-server/pkg/storage"
)

type userService struct {}

var UserService = new(userService)

func (s *userService) Login(username, password string) (*model.User, string, error) {
	var u model.User
	user := s.FindByUsername(username)
	if user.Username == "" {
		return nil, "", errors.New("User Not Found")
	}

	token, err := u.Login(user, password)
	if err != nil {
		return nil, "", err
	}

	return user, token, err
}

func (s *userService) FindByUsername(username string) *model.User {
	for _, u := range storage.GetUsers() {
		if u != nil && u.Username == username {
			return u
		}
	}
	return nil
}

func (s *userService) GetUsers() []*model.User {
	return storage.GetUsers()
}

func (s *userService) SignUp(params map[string]interface{}) *model.User {
	var u model.User
	if params["username"] == "" {
		return nil
	}

	username := params["username"].(string)
	user := s.FindByUsername(username)
	if user != nil {
		return nil
	}

	user = u.New(params)
	storage.AddUser(user)
	return user
}