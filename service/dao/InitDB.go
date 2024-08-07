package dao

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	DB *gorm.DB
)

// 初始化数据库
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
func setCount() {
	DB.Exec("UPDATE question_answer_counts SET answer_count = 0")
}
