package main

import (
	"qasociety/api"
	"qasociety/service/dao"
)

func main() {
	//初始化数据库
	dao.InitDB()
	//初始化redis
	dao.InitRdb()
	//多线程,设置定时器,定期清除redis中的当日回复数缓存
	go dao.StartCacheCleanup()
	//初始化路由
	api.InitRouter()
}
