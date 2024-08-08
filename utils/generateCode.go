package utils

import (
	"math/rand"
	"time"
)

// GenerateCode 生成验证码
func GenerateCode() string {
	numeric := [10]string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())
	var sb string
	for i := 0; i < 6; i++ {
		sb += numeric[rand.Intn(r)]
	}
	return sb
}
