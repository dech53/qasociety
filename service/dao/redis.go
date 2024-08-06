package dao

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

// StartCacheCleanup 10s 后清空 answerCount
func StartCacheCleanup() {
	ctx := context.Background()
	// 清理周期，测试阶段用 10s，实际使用采取 24h 刷新制git
	ticker := time.NewTicker(24 * time.Hour)
	for range ticker.C {
		err := cleanupRedisCache(ctx)
		if err != nil {
			log.Printf("Error cleaning up Redis cache: %v", err)
		}
	}
}

// cleanupRedisCache 将 answerCount 清零
func cleanupRedisCache(ctx context.Context) error {
	//删除整个list
	Rdb.Del(ctx, "questions")
	pageSize := 5
	questionCount, _ := GetQuestionsCount()
	offset := 0
	for offset < questionCount+pageSize {
		questions, err := FindQuestionByPattern("", "", offset, pageSize)
		if err != nil {
			return err
		}
		for _, question := range questions {
			// 将 question 序列化为 JSON 字符串
			questionJson, err := json.Marshal(question)
			// 将 answerCount 重置为 0
			_, err = Rdb.ZAdd(ctx, "questions", &redis.Z{
				Score:  0,
				Member: string(questionJson),
			}).Result()
			if err != nil {
				return err
			}
		}
		offset += pageSize
	}
	return nil
}
