package main

import (
	"qasociety/api"
	"qasociety/service/dao"
	"qasociety/utils"
)

func main() {
	//初始化数据库
	dao.InitDB()
	//初始化redis
	dao.InitRdb()
	//24h后问题数量清零
	go dao.SetAnswerCountZero()
	//定期更新redis中的热榜前十问题
	go dao.StartUpdateRedisCache()
	//订阅点赞信息频道
	go utils.SubscribeToLikeChannel()
	//固定周期将redis数据写入mysql
	go dao.WriteMysqlFromRedis()
	//初始化路由
	api.InitRouter()
}
