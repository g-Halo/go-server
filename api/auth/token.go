package auth

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/g-Halo/go-server/internal/logic/service"
	"github.com/g-Halo/go-server/pkg/pb"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/g-Halo/go-server/conf"
	"github.com/g-Halo/go-server/pkg/util"
	"golang.org/x/net/context"
)

type AuthServer struct {}

type Token struct {
	Username string
	Password string
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func (s *AuthServer) SignIn(ctx context.Context, in *pb.AuthReq) (*pb.AuthResp, error) {
	username := in.GetUsername()
	password := in.GetPassowrd()

	user := service.UserService.FindByUsername(username)
	if user == nil {
		return &pb.AuthResp{Code: util.Fail, Msg: "无效的用户"}, nil
	}

	salt := user.Salt
	m5 := md5.New()
	m5.Write([]byte(salt))
	m5.Write([]byte(password))
	st := m5.Sum(nil)
	if hex.EncodeToString(st) != user.Password {
		return &pb.AuthResp{Code: util.Fail, Msg: "用户名或者密码错误"}, nil
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)

	// FIXME: 调试阶段：token有效期 24 小时
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(24)).Unix()
	claims["iat"] = time.Now().Unix()
	claims["username"] = username

	token.Claims = claims

	tokenString, _ := token.SignedString([]byte(conf.Conf.SecretKey))

	return &pb.AuthResp{Code: util.Success, Data: tokenString, Msg: "登录成功"}, nil
}

func (s *AuthServer) SignUp(ctx context.Context, in *pb.SignUpReq) (*pb.SignUpResp, error) {
	params := map[string]interface{}{
		"nickname": in.GetUsername(),
		"username": in.GetNickname(),
		"password": in.GetPassword(),
	}

	user := service.UserService.SignUp(params)
	if user == nil {
		return &pb.SignUpResp{Code: util.Fail, Msg: "注册失败"}, nil
	} else {
		return &pb.SignUpResp{Code: util.Success, Username: user.Username, Msg: "注册成功"}, nil
	}
}

func (s *AuthServer) Validate(ctx context.Context, in *pb.ValidateReq) (*pb.ValidateResp, error) {
	token, _ := jwt.ParseWithClaims(in.GetTokenStr(), &util.MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(conf.Conf.SecretKey), nil
	})

	if claims, ok := token.Claims.(*util.MyCustomClaims); ok && token.Valid {
		user := service.UserService.FindByUsername(claims.Username)
		if user == nil {
			return &pb.ValidateResp{Code: util.Fail, Msg: "无效的 Token"}, nil
		} else {
			return &pb.ValidateResp{Code: util.Success, Username: user.Username, Msg: "登录成功"}, nil
		}
	}
	return &pb.ValidateResp{Code: util.Fail, Msg: "未知的 Token 错误"}, nil
}