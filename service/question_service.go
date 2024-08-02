package service

import (
	"qasociety/model"
	"qasociety/service/dao"
	"time"
)

// CreateQuestion 创建问题
func AddQuestion(userID int, title, content string) error {
	// 创建问题实例
	question := model.Question{
		UserID:      userID,
		Title:       title,
		Content:     content,
		CreatedTime: time.Now(),
		UpdatedTime: time.Now(),
	}
	// 将问题存储到数据库
	return dao.AddQuestion(question)
}
