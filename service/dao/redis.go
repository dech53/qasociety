package dao

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

// StartUpdateRedisCache 开始更新redis中存放的热门问题
func StartUpdateRedisCache() {
	ctx := context.Background()
	ticker := time.NewTicker(2 * time.Second)
	for range ticker.C {
		Rdb.ZRemRangeByRank(ctx, "questions", 0, -1)
		updateRedis()
	}
}

// 从mysql中获取新增数量前十的热门问题
func updateRedis() {
	ctx := context.Background()
	questionACs, err := GetTopQuestions("desc", 0, 10)
	if err != nil {
		log.Printf("Error finding questions: %v", err)
	}
	for _, questionAC := range questionACs {
		question, err := GetQuestionByID(questionAC.QuestionID)
		questionJson, err := json.Marshal(question)
		Rdb.ZAdd(ctx, "questions", &redis.Z{
			Score:  float64(questionAC.AnswerCount),
			Member: string(questionJson),
		})
		if err != nil {
			log.Printf("Error finding questions: %v", err)
		}
	}
}

// DeleteThumbRedis 取消点赞
func DeleteThumbRedis(redisKey string, userID int) error {
	ctx := context.Background()
	_, err := Rdb.SRem(ctx, redisKey, userID).Result()
	return err
}
