package utils

import (
	"context"
	"gorm.io/gorm"
	"qasociety/model"
	"qasociety/service/dao"
	"strconv"
)

// IsUserLikedAnswer 查找用户是否已点赞
func IsUserLikedAnswer(userID, answerId int) (bool, error) {
	ctx := context.Background()
	// 构建Redis Set 的 key
	setKey := "answer:likes:" + strconv.Itoa(answerId)
	// 检查 userID 是否是 Set setKey 的成员
	isMember, err := dao.Rdb.SIsMember(ctx, setKey, userID).Result()
	if err != nil {
		// 如果发生错误，返回错误并通知调用者
		return false, err
	}
	if isMember {
		return true, nil
	}
	var like model.Like
	result := dao.DB.Where("user_id = ? AND answer_id = ?", userID, answerId).First(&like)
	if result.Error == gorm.ErrRecordNotFound {
		return false, nil
	}
	if result.Error != nil {
		return false, result.Error
	}
	return true, nil
}
