package service

import (
	"qasociety/model"
	"qasociety/service/dao"
	"time"
)

func AddAnswer(userID, questionid int, content string) error {
	answer := model.Answer{
		QuestionID:  questionid,
		UserID:      userID,
		Content:     content,
		CreatedTime: time.Now(),
	}
	return dao.AddAnswer(answer)
}
