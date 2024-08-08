package model

import "github.com/dgrijalva/jwt-go"

// User 用户结构体
type User struct {
	ID       int    `gorm:"primaryKey;autoIncrement;unique;not null" json:"id"`
	Username string `gorm:"size:50;not null" json:"username"`      // 增加长度到 50
	Password string `gorm:"size:60;not null" json:"password"`      // 增加长度到 60
	Email    string `gorm:"size:255;unique;not null" json:"email"` // 增加长度到 255
}
type MyClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
