package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"qasociety/model"
	"time"
)

// AddQuestion 创建问题
func AddQuestion(question model.Question) error {
	questionACount := model.QuestionAnswerCount{
		QuestionID:  question.ID,
		AnswerCount: 0,
	}
	DB.Model(&model.QuestionAnswerCount{}).Create(&questionACount)
	return DB.Model(&model.Question{}).Create(&question).Error
}

// GetQuestionByID 通过ID查找问题
func GetQuestionByID(questionID int) (*model.Question, error) {
	var question model.Question
	result := DB.Take(&question, questionID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &question, nil
}

// UpdateQuestion 更新问题
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

// DeleteQuestion 删除问题
func DeleteQuestion(questionID int) error {
	ctx := context.Background()
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
	question, err := GetQuestionByID(questionID)
	if err != nil {
		return err
	}
	questionJson, err := json.Marshal(question)
	DB.Delete(&model.QuestionAnswerCount{}, questionID)
	Rdb.ZRem(ctx, "questions", string(questionJson))
	questions, err := GetTopQuestions("desc", 0, 1)
	question, err = GetQuestionByID(questions[0].QuestionID)
	questionJson, err = json.Marshal(question)
	Rdb.ZAdd(ctx, "questions", &redis.Z{
		Score:  float64(questions[0].AnswerCount),
		Member: string(questionJson),
	})
	return DB.Delete(&model.Question{}, questionID).Error
}

// FindQuestionByPattern 通过特征查找问题
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

// GetQuestionsByRedis 从redis中获取前十的问题
func GetQuestionsByRedis() ([]model.Question, error) {
	//返回一个Object
	var questions []model.Question
	ctx := context.Background()
	//需要修改成类似gorm中分页查询的形式,添加offset之类的
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

// GetTopQuestions 获取新增回复数前十的问题
func GetTopQuestions(order string, offset, pageSize int) ([]model.QuestionAnswerCount, error) {
	var questions []model.QuestionAnswerCount
	err := DB.Order("answer_count " + order).
		Offset(offset).
		Limit(pageSize).
		Find(&questions).Error
	if err != nil {
		return nil, err
	}
	return questions, nil
}

// GetQuestionLikesCount 获取问题点赞数
func GetQuestionLikesCount(id int) (int, error) {
	var counts int
	answers, err := GetAllAnswers(id)
	if err != nil {
		return 0, err
	}
	for _, answer := range answers {
		answerCount, err := GetAnswerLikesCount(answer.ID)
		if err != nil {
			return 0, err
		}
		counts += answerCount
	}
	return counts, nil
}
