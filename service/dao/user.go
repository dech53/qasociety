package dao

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"qasociety/model"
	"strconv"
	"time"
)

// AddUser 添加用户
func AddUser(user model.User) (model.User, error) {
	result := DB.Create(&user)
	return user, result.Error
}

// SelectUserName 查重用户名
func SelectUserName(username string) (bool, error) {
	result := DB.Where("username = ?", username).First(&model.User{})
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// 记录未找到则表明用户名未重复
		return false, nil
	}
	if result.Error != nil {
		// 查询出现错误，返回 false 和错误
		return true, result.Error
	}
	if result.RowsAffected != 0 {
		// 用户名存在
		return true, nil
	}
	// 用户名不存在
	return false, nil
}

// SelectPassword 查找用户密码
func SelectPassword(key string, pattern string) (string, error) {
	var user model.User
	result := DB.Where(pattern+" = ?", key).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return "", result.Error
	}
	return user.Password, nil
}
func GetUserIDByUsername(username string) (int, error) {
	var user model.User
	result := DB.Model(&model.User{}).Where("username = ?", username).First(&user)
	fmt.Println(user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return 0, nil // 用户未找到
		}
		return 0, result.Error // 查询错误
	}
	return user.ID, nil // 返回用户ID
}
func GetUserByPattern(pattern, value string) (model.User, error) {
	var user model.User
	err := DB.Model(&model.User{}).Where(pattern+" = ?", value).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return user, nil
	}
	return user, err
}
func SetCodeRedis(userID string, code string) (bool, error) {
	ctx := context.Background()
	result, err := Rdb.SetNX(ctx, userID, code, 3*time.Minute).Result()
	return result, err
}
func GetExpireTime(userID string) (time.Duration, error) {
	ctx := context.Background()
	restTime, err := Rdb.TTL(ctx, userID).Result()
	return restTime, err
}
func VerifyCode(email, code string) (bool, error) {
	ctx := context.Background()
	user, err := GetUserByPattern("email", email)
	if err != nil {
		return false, err
	}
	userID := strconv.Itoa(user.ID)
	realCode, err := Rdb.Get(ctx, userID).Result()
	if realCode != code {
		return false, errors.New("验证码错误")
	}
	return true, nil
}
func ResetPassword(email, newPassword string) error {
	user, err := GetUserByPattern("email", email)
	if err != nil {
		return err
	}
	if newPassword == "" {
		return errors.New("新密码不能为空")
	}
	user.Password = newPassword
	err = DB.Save(&user).Error
	if err != nil {
		return err
	}
	return nil
}
