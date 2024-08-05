package dao

import (
	"qasociety/model"
)

// AddAnswer 添加回复
func AddAnswer(answer model.Answer) error {
	return DB.Model(&model.Answer{}).Create(&answer).Error
}

// FindAnswersByPattern 通过pattern查找相似回复
func FindAnswersByPattern(questionID int, pattern string, offset int, pageSize int) ([]model.Answer, error) {
	var answers []model.Answer
	err := DB.Where("question_id = ? AND content LIKE ?", questionID, "%"+pattern+"%").
		Offset(offset).
		Limit(pageSize).
		Find(&answers).Error
	if err != nil {
		return nil, err
	}
	return answers, nil
}

// GetAnswerByID 通过id查找answer
func GetAnswerByID(answerID int) (*model.Answer, error) {
	var answer model.Answer
	err := DB.First(&answer, "id = ?", answerID).Error
	if err != nil {
		return nil, err
	}
	return &answer, nil
}

// RemoveAnswer 删除回复
func RemoveAnswer(answer *model.Answer) error {
	comments, err := GetAllComments(answer.ID)
	if err != nil {
		return err
	}
	err = RemoveComments(comments)
	if err != nil {
		return err
	}
	return DB.Delete(answer).Error
}

// GetAllAnswers 通过问题ID获取其ID下的所有回复
func GetAllAnswers(questionID int) ([]model.Answer, error) {
	var answers []model.Answer
	err := DB.Where("question_id = ?", questionID).Find(&answers).Error
	if err != nil {
		return nil, err
	}
	return answers, nil
}
