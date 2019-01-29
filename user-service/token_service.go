package main

import (
	pb "github.com/porcorosso/shippy/user-service/proto/user"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type Authable interface {
	Decode(tokenStr string) (*CustomClaims, error)
	Encode(user *pb.User) (string, error)
}

// 定义哈希密码所用的盐，要保证其生成和保存都足够安全，可用md5来生成
var privateKey = []byte("`xsa$39(")

// 自定义的metadata，在加密后作为JWT的第二部分返回给客户端
type CustomClaims struct {
	User *pb.User
	jwt.StandardClaims // 使用标准的payload
}

type TokenService struct {
	repo Repository
}

// 将User用户信息加密为JWT字符串
func (srv *TokenService) Encode(user *pb.User) (string, error) {
	// 3天后过期
	expireTime := time.Now().Add(time.Hour * 24 *3).Unix()
	claims := CustomClaims{
		user,
		jwt.StandardClaims{
			Issuer: "go.micro.srv.user",
			ExpiresAt: expireTime,
		},
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return jwtToken.SignedString(privateKey)
}

// 将JWT字符串解密为CustomClaims对象
func (srv *TokenService) Decode(tokenStr string) (*CustomClaims, error) {
	t, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return privateKey, nil
	})

	// 解密转换类型并返回值
	if claims, ok := t.Claims.(*CustomClaims); ok && t.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}