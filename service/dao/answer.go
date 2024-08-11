package dao

import (
	"context"
	"fmt"
	"qasociety/model"
	"strconv"
)

// AddAnswer 添加回复
func AddAnswer(answer model.Answer) error {
	DB.Exec("UPDATE question_answer_counts SET answer_count = answer_count + 1 WHERE question_id = ?", answer.QuestionID)
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

// GetAnswerLikesCount 获取回复点赞数
func GetAnswerLikesCount(answerID int) (int, error) {
	mysqlCountResult := DB.Exec("SELECT COUNT(*) FROM likes where answer_id = ?", answerID)
	if mysqlCountResult.Error != nil {
		return 0, mysqlCountResult.Error
	}
	mysqlCount := mysqlCountResult.RowsAffected
	fmt.Println(mysqlCount)
	redisKey := "answer:likes:" + strconv.Itoa(answerID)
	redisCount, err := Rdb.SCard(context.Background(), redisKey).Result()
	if err != nil {
		return 0, err
	}
	totalLikes := int(mysqlCount) + int(redisCount) + 1
	return totalLikes, nil
}
