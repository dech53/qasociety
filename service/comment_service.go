package service

import (
	"errors"
	"qasociety/model"
	"qasociety/service/dao"
	"time"
)

// CreateComment 创建评论
func CreateComment(userID, answerID int, content string) error {
	comment := model.Comment{
		UserID:      userID,
		AnswerID:    answerID,
		Content:     content,
		CreatedTime: time.Now(),
	}
	return dao.AddComment(comment)
}

// GetCommentsByAnswerID 通过回复ID获取对应的评论
func GetCommentsByAnswerID(answerID, page, pageSize int) ([]model.Comment, error) {
	// 计算分页的起始位置
	offset := (page - 1) * pageSize
	// 调用 DAO 层函数进行分页查询
	answers, err := dao.FindCommentsByID(answerID, offset, pageSize)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return answers, nil
}
