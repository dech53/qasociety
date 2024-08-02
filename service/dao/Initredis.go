package dao

//初始化redis
import (
	"github.com/go-redis/redis/v8"
)

var rdb *redis.Client

func InitRdb() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "123456",
		DB:       0,
	})
}
