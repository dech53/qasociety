package dao

import (
	"context"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"qasociety/model"
	"strconv"
	"strings"
	"time"
)

var (
	DB *gorm.DB
)

// InitDB 初始化数据库
func InitDB() {
	//数据库基本参数配置
	username := "root"
	password := "root"
	host := "127.0.0.1"
	post := 3306
	database := "qasociety"
	timeout := "10s"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?timeout=%s&charset=utf8mb4&loc=Local&parseTime=true", username, password, host, post, database, timeout)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		panic("数据库连接失败")
	}
	DB = db
	fmt.Println("数据库连接成功")
}

// SetAnswerCountZero 24h后将question_answer_counts中的新增count数清零
func SetAnswerCountZero() {
	//时间可以调整实现半小时内新增之类的
	ticker := time.NewTicker(24 * time.Hour)
	for range ticker.C {
		setCount()
	}
}

// 清空点赞数
func setCount() {
	DB.Exec("UPDATE question_answer_counts SET answer_count = 0")
}

// WriteMysqlFromRedis 从redis中读取数据并写入mysql中
func WriteMysqlFromRedis() {
	ticker := time.NewTicker(5 * time.Minute)
	for range ticker.C {
		ctx := context.Background()
		// 获取所有匹配的Redis键
		keys, err := Rdb.Keys(ctx, "answer:likes:*").Result()
		if err != nil {
			fmt.Println("未找到相应的键:", err)
			continue
		}
		for _, key := range keys {
			members, _ := Rdb.SMembers(ctx, key).Result()
			for _, member := range members {
				memberInt, _ := strconv.Atoi(member)
				like := model.Like{
					AnswerID: ExtractAnswerIDFromKey(key),
					UserID:   memberInt,
				}
				DB.Create(&like).Model(&model.Like{})
			}
			Rdb.Del(ctx, key)
		}
	}
}

// DeleteThumbMysql 从mysql中删除点赞信息
func DeleteThumbMysql(answerID, userID int) error {
	like := model.Like{
		AnswerID: answerID,
		UserID:   userID,
	}
	return DB.Delete(&like).Model(&model.Like{}).Error
}

// ExtractAnswerIDFromKey 工具函数
func ExtractAnswerIDFromKey(key string) int {
	parts := strings.Split(key, ":")
	answerID, _ := strconv.Atoi(parts[len(parts)-1])
	return answerID
}
