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
	//初始化路由
	dao.GetUserIDByUsername("wx")
	api.InitRouter()
	//test区域

}
