package model

import "time"

// Question 代表问题表
type Question struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement;unique;not null"`
	UserID      int       `json:"user_id" gorm:"index;not null"` // 关联用户ID
	Title       string    `json:"title" gorm:"not null"`
	Content     string    `json:"content" gorm:"type:text"`
	CreatedTime time.Time `json:"created_time" gorm:"autoCreateTime"` // 自动创建时间
	UpdatedTime time.Time `json:"updated_time"`                       // 最新评论的创建时间
}

// Answer 代表回答表
type Answer struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement;unique;not null"`
	QuestionID  int       `json:"question_id" gorm:"index;not null"` // 关联问题ID
	UserID      int       `json:"user_id" gorm:"index;not null"`     // 关联用户ID
	Content     string    `json:"content" gorm:"type:text;not null"`
	CreatedTime time.Time `json:"created_time" gorm:"autoCreateTime"` // 自动创建时间
}

// Comment 代表评论表
type Comment struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement;unique;not null"`
	AnswerID    int       `json:"answer_id" gorm:"index;not null"` // 关联回答ID
	UserID      int       `json:"user_id" gorm:"index;not null"`   // 关联用户ID
	Content     string    `json:"content" gorm:"type:text;not null"`
	CreatedTime time.Time `json:"created_time" gorm:"autoCreateTime"` // 自动创建时间
}
