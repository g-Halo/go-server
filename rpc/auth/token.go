package auth

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/rpc"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/g-Halo/go-server/logger"
	"github.com/g-Halo/go-server/model"
	"github.com/g-Halo/go-server/util"
)

type Token struct {
	Username string
	Password string
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

const SecretKey = "b66630cedd9cd5c34cf91e7"

// 传入用户名、密码，创建 token
func (Token) Create(t *Token, reply *util.Response) error {
	// TODO: 用户名、密码校验

	client, err := rpc.Dial("tcp", ":7302")
	if err != nil {
		fmt.Println("无效的地址")
		os.Exit(0)
	}
	defer client.Close()

	var user *model.User
	if err = client.Call("Logic.FindByUsername", &t.Username, &user); err != nil {
		logger.Error(err.Error())
	}

	if user == nil {
		*reply = util.Response{Code: util.Fail, Msg: "无效的用户"}
		return nil
	} else {
		salt := user.Salt
		m5 := md5.New()
		m5.Write([]byte(salt))
		m5.Write([]byte(string(t.Password)))
		st := m5.Sum(nil)
		if hex.EncodeToString(st) != user.Password {
			*reply = util.Response{Code: util.Fail, Msg: "用户名或者密码错误"}
			return nil
		}

		// 颁发 token
		token := jwt.New(jwt.SigningMethodHS256)
		claims := make(jwt.MapClaims)
		claims["exp"] = time.Now().Add(time.Hour * time.Duration(1)).Unix()
		claims["iat"] = time.Now().Unix()
		claims["username"] = t.Username

		token.Claims = claims

		tokenString, _ := token.SignedString([]byte(SecretKey))

		*reply = util.Response{Code: util.Success, Data: tokenString, Msg: "登录成功"}
	}

	return nil
}

// 传入 token，并检验 token 的合法性
func (Token) Validate(arg *string, reply *int) error {
	tokenString := *arg

	token, _ := jwt.ParseWithClaims(tokenString, &util.MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	checkResult := false
	if _, ok := token.Claims.(*util.MyCustomClaims); ok && token.Valid {
		//user := logic.UserLogic.FindByUsername(claims.Username)
		checkResult = true
	}

	if checkResult {
		*reply = 1
	} else {
		*reply = 0
	}
	return nil
}
