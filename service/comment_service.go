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

// GetComments 通过回复ID获取对应的评论
func GetComments(answerID, page, pageSize int, order string) ([]model.Comment, error) {
	// 计算分页的起始位置
	offset := (page - 1) * pageSize
	// 调用 DAO 层函数进行分页查询
	answers, err := dao.FindCommentsByID(answerID, offset, pageSize, order)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return answers, nil
}

// GetCommentByID 通过评论ID获取评论
func GetCommentByID(commentID int) (*model.Comment, error) {
	comment, err := dao.FindCommentByID(commentID)
	if err != nil {
		return nil, errors.New("评论不存在")
	}
	return comment, nil
}

// DeleteComment 删除单条评论
func DeleteComment(commentID int) error {
	err := dao.RemoveCommentByID(commentID)
	if err != nil {
		return errors.New("删除评论失败")
	}
	return nil
}
