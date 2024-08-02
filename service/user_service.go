package service

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"qasociety/api/middleware"
	"qasociety/model"
	"qasociety/service/dao"
	"time"
)

// RegisterUser 处理用户注册
func RegisterUser(username, password, email string) error {
	// 检查用户名是否已存在
	exists, err := dao.SelectUserName(username)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("用户名已存在")
	}
	// 创建新用户
	user := model.User{
		Username: username,
		Password: password,
		Email:    email,
	}
	_, err = dao.AddUser(user)
	if err != nil {
		return err
	}
	return nil
}

// LoginUser 处理用户登录
func LoginUser(username, password, email string) (string, error) {
	var savedPassword string
	var err error
	// 根据用户名或邮箱查找用户密码
	if username != "" {
		savedPassword, err = dao.SelectPassword(username, "username")
	} else if email != "" {
		savedPassword, err = dao.SelectPassword(email, "email")
	} else {
		return "", errors.New("用户名或邮箱不能为空")
	}
	if err != nil {
		return "", err
	}
	// 检查密码是否正确
	if savedPassword != password {
		return "", errors.New("密码错误")
	}
	// 生成 JWT token
	claim := model.MyClaims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 2).Unix(),
			Issuer:    "YXH",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString(middleware.Secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
