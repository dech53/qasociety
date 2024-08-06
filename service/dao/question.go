package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"qasociety/model"
	"time"
)

func AddQuestion(question model.Question) error {
	id := GetTheLastID()
	question.ID = id + 1
	//缓存中设置初始ID值为0,过期时间为永久
	ctx := context.Background()
	questionJson, err := json.Marshal(question)
	if err != nil {
		return err
	}
	//answerCount := 0
	// 设置set以score排序
	_, err = Rdb.ZAdd(ctx, "questions", &redis.Z{
		Score:  5,
		Member: string(questionJson),
	}).Result()
	if err != nil {
		return err
	}
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
	ctx := context.Background()
	var question model.Question
	// 查找问题
	DB.First(&question, questionID)
	questionJson, err := json.Marshal(question)
	if err != nil {
		return err
	}
	//{"id":1,"user_id":1,"title":"测试问题1","content":"测试问题1内容","created_time":"2024-08-02T22:33:45.328+08:00","updated_time":"2024-08-06T13:22:33.961+08:00"}
	fmt.Println(string(questionJson))
	questionScore, _ := Rdb.ZScore(ctx, "questions", string(questionJson)).Result()
	fmt.Println(questionScore)
	Rdb.ZRem(ctx, "questions", string(questionJson))
	// 更新问题的标题和内容
	question.Title = title
	question.Content = content
	question.UpdatedTime = time.Now()
	//更新mysql中的question
	if err := DB.Save(&question).Error; err != nil {
		return err
	}
	// 保存更新
	DB.First(&question, questionID)
	questionJson, err = json.Marshal(question)
	Rdb.ZAdd(ctx, "questions", &redis.Z{
		Score:  questionScore,
		Member: questionJson,
	})
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
func FindQuestionByPattern(pattern, order string, offset, pageSize int) ([]model.Question, error) {
	var questions []model.Question
	err := DB.Where("content LIKE ?", "%"+pattern+"%").
		Order("updated_time " + order).
		Offset(offset).
		Limit(pageSize).
		Find(&questions).Error
	if err != nil {
		return nil, err
	}
	return questions, nil
}
func GetQuestionsCount() (int, error) {
	var count int64
	err := DB.Model(&model.Question{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}
func GetTheLastID() int {
	var max int
	DB.Model(&model.Question{}).Select("max(id) as id").Scan(&max)
	return max
}
func GetQuestionsByRedis() ([]model.Question, error) {
	//返回一个Object
	var questions []model.Question
	ctx := context.Background()
	member, err := Rdb.ZRange(ctx, "questions", 0, -1).Result()
	if err != nil {
		return nil, err
	}
	for _, v := range member {
		var question model.Question
		err = json.Unmarshal([]byte(v), &question)
		questions = append(questions, question)
	}
	if err != nil {
		return nil, err
	}
	return questions, nil
}
