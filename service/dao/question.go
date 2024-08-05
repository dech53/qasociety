package dao

import (
	"errors"
	"gorm.io/gorm"
	"qasociety/model"
	"time"
)

func AddQuestion(question model.Question) error {
	return DB.Model(&model.Question{}).Create(&question).Error
}
func GetQuestionByID(questionID int) (*model.Question, error) {
	var question model.Question
	result := DB.Take(&question, questionID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &question, nil
}
func UpdateQuestion(questionID int, title, content string) error {
	var question model.Question
	// 查找问题
	if err := DB.First(&question, questionID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("不存在该问题")
		}
		return err
	}
	// 更新问题的标题和内容
	question.Title = title
	question.Content = content
	question.UpdatedTime = time.Now()
	// 保存更新
	if err := DB.Save(&question).Error; err != nil {
		return err
	}
	return nil
}
func DeleteQuestion(questionID int) error {
	answers, err := GetAllAnswers(questionID)
	if err != nil {
		return err
	}
	for _, answer := range answers {
		err = RemoveAnswer(&answer)
	}
	if err != nil {
		return err
	}
	return DB.Delete(&model.Question{}, questionID).Error
}
