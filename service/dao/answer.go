package dao

import "qasociety/model"

// 添加回复
func AddAnswer(answer model.Answer) error {
	return DB.Model(&model.Answer{}).Create(&answer).Error
}

// 通过pattern查找相似回复
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

// 通过id查找answer
func GetAnswerByID(answerID int) (*model.Answer, error) {
	var answer model.Answer
	err := DB.First(&answer, "id = ?", answerID).Error
	if err != nil {
		return nil, err
	}
	return &answer, nil
}
