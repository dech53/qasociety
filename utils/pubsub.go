package utils

import (
	"context"
	"encoding/json"
	"log"
	"qasociety/model"
	"qasociety/service/dao"
	"strconv"
)

// PublishLikeEvent 发送订阅信息
func PublishLikeEvent(userID int, answerId int) error {
	likeEvent := model.Like{
		AnswerID: answerId,
		UserID:   userID,
	}
	jsonEvent, err := json.Marshal(likeEvent)
	if err != nil {
		return err
	}
	return dao.Rdb.Publish(context.Background(), "like_channel", string(jsonEvent)).Err()
}

// SubscribeToLikeChannel 订阅Redis点赞频道，并更新Redis Set
func SubscribeToLikeChannel() {
	pubsub := dao.Rdb.Subscribe(context.Background(), "like_channel")
	for msg := range pubsub.Channel() {
		likeEvent := model.Like{}
		if err := json.Unmarshal([]byte(msg.Payload), &likeEvent); err != nil {
			log.Printf("Error decoding like event: %v", err)
			continue
		}
		// 更新Redis Set
		redisKey := "answer:likes:" + strconv.Itoa(likeEvent.AnswerID)
		if _, err := dao.Rdb.SAdd(context.Background(), redisKey, likeEvent.UserID).Result(); err != nil {
			log.Printf("Error updating Redis set for like: %v", err)
		}
	}
}
