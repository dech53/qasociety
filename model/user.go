package model

import "github.com/dgrijalva/jwt-go"

type User struct {
	ID       int    `gorm:"primaryKey;autoIncrement;unique;not null" json:"id"`
	Username string `gorm:"not null;type:varchar(20)" json:"username"`
	Password string `gorm:"not null;type:varchar(20)" json:"password"`
	Email    string `gorm:"unique;type:varchar(25);not null" json:"email"`
}
type MyClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
