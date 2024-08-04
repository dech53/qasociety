package service

//调用dao层
import (
	"errors"
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

func SearchAnswersByPattern(id int, pattern string, page, pageSize int) ([]model.Answer, error) {
	// 计算分页的起始位置
	offset := (page - 1) * pageSize
	// 调用 DAO 层函数进行分页查询
	answers, err := dao.FindAnswersByPattern(id, pattern, offset, pageSize)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return answers, nil
}
