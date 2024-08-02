package dao

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"qasociety/model"
)

// 添加用户
func AddUser(user model.User) (model.User, error) {
	result := DB.Create(&user)
	return user, result.Error
}

// 查重用户名
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

// 查找用户密码
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
