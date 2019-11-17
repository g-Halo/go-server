package logic

import (
	"github.com/g-Halo/go-server/model"
)

type Logic struct{}

// rpc 方法
func (logic *Logic) FindByUsername(username *string, user *model.User) error {
	u := UserLogic.FindByUsername(*username)
	if u != nil {
		*user = *u
	}
	return nil
}

func (logic *Logic) SignUp(params map[string]interface{}, user *model.User) error {
	u := UserLogic.SignUp(params)
	if u != nil {
		*user = *u
	}
	return nil
}

func (logic *Logic) GetUsers(params string, reply *[]*model.User) error {
	us := UserLogic.GetUsers()
	*reply = us
	return nil
}
