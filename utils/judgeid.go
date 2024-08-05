package utils

import (
	"errors"
	"github.com/gin-gonic/gin"
	"qasociety/service"
	"qasociety/service/dao"
	"strconv"
)

func JudgeID(c *gin.Context) (int, error) {
	// 从上下文获取用户名
	username, exists := c.Get("username")
	if !exists {
		return 0, errors.New("用户未认证")
	}
	// 根据用户名获取用户ID
	userID, err := dao.GetUserIDByUsername(username.(string))
	if err != nil || userID == 0 {
		return 0, errors.New("用户不存在")
	}
	// 获取请求ID
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, errors.New("无效的问题 ID")
	}
	// 获取问题
	question, err := service.GetQuestionByID(id)
	if err != nil {
		return 0, errors.New("获取问题失败")
	}
	if question.UserID != userID {
		return 0, errors.New("无权变更问题")
	}
	return id, nil
}
func GetUserID(c *gin.Context) (int, error) {
	username, exists := c.Get("username")
	if !exists {
		return 0, errors.New("用户未认证")
	}
	return dao.GetUserIDByUsername(username.(string))
}
