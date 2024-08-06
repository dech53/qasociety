package service

import (
	"errors"
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

// 通过ID查找问题
func GetQuestionByID(questionID int) (*model.Question, error) {
	return dao.GetQuestionByID(questionID)
}

// 通过ID更新问题
func UpdateQuestion(questionID int, title, content string) error {
	return dao.UpdateQuestion(questionID, title, content)
}

// 通过ID删除问题
func DeleteQuestion(questionID int) error {
	return dao.DeleteQuestion(questionID)
}

// pattern不为空时使用这个
func SearchQuestionsByPattern(pattern, order string, page, pageSize int) ([]model.Question, error) {
	// 计算分页的起始位置
	offset := (page - 1) * pageSize
	// 调用 DAO 层函数进行分页查询
	questions, err := dao.FindQuestionByPattern(pattern, order, offset, pageSize)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return questions, nil
}

// pattern为空时使用这个,从redis中获取question
func GetQuestionsByRedis() ([]model.Question, error) {
	return dao.GetQuestionsByRedis()
}
