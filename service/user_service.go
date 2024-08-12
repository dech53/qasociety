package service

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"qasociety/api/middleware"
	"qasociety/model"
	"qasociety/service/dao"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
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
	// 使用 MD5 对密码进行加密
	hashedPassword := md5.New()
	hashedPassword.Write([]byte(password))
	passwordHash := hex.EncodeToString(hashedPassword.Sum(nil))
	// 创建新用户
	user := model.User{
		Username: username,
		Password: passwordHash,
		Email:    email,
	}
	_, err = dao.AddUser(user)
	if err != nil {
		return err
	}
	return nil
}

// LoginUser 处理用户登录
func LoginUser(username, password, email, info, mode string) (string, error) {
	ctx := context.Background()
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
	// 使用 MD5 对密码进行加密
	hashedPassword := md5.New()
	hashedPassword.Write([]byte(password))
	passwordHash := hex.EncodeToString(hashedPassword.Sum(nil))
	if err != nil {
		return "", err
	}
	if mode == "" {
		// 检查密码是否正确
		if savedPassword != passwordHash {
			return "", errors.New("密码错误")
		}
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
	// 构建 Redis 键
	redisKey := "session:" + username + ":" + info
	// 存储 token 到 Redis，设置过期时间为 24 小时
	err = dao.Rdb.SetNX(ctx, redisKey, tokenString, 24*time.Hour).Err()
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
func GetUserByPattern(pattern, value string) (model.User, error) {
	return dao.GetUserByPattern(pattern, value)
}
func ResetRequest(code string, user model.User) (bool, error) {
	userID := strconv.Itoa(user.ID)
	return dao.SetCodeRedis(userID, code, "resetPassword")
}
func GetExpireTime(user model.User, pattern string) (time.Duration, error) {
	userID := strconv.Itoa(user.ID)
	return dao.GetExpireTime(userID, pattern)
}
func VerifyCode(email, code, pattern string) (bool, error) {
	return dao.VerifyCode(email, code, pattern)
}
func ResetPassword(email, newPassword string) error {
	return dao.ResetPassword(email, newPassword)
}
func LoginByCodeRequest(code string, user model.User) (bool, error) {
	userID := strconv.Itoa(user.ID)
	fmt.Println(userID)
	return dao.SetCodeRedis(userID, code, "loginCode")
}
