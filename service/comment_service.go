package service

import (
	"qasociety/model"
	"qasociety/service/dao"
	"time"
)

func CreateComment(userID, answerID int, content string) error {
	comment := model.Comment{
		UserID:      userID,
		AnswerID:    answerID,
		Content:     content,
		CreatedTime: time.Now(),
	}
	return dao.AddComment(comment)
}
