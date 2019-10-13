package util

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/g-Halo/go-server/conf"
	"time"
)

type MyCustomClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func CreateJWT(options map[string]interface{}) string {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(1)).Unix()
	claims["iat"] = time.Now().Unix()
	for k, v := range options {
		claims[k] = v
	}
	token.Claims = claims

	config := conf.LoadConf()
	tokenString, _ := token.SignedString([]byte(config.SecretKey))

	return tokenString
}
