package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"qasociety/service"
	"qasociety/service/dao"
	"qasociety/utils"
	"strconv"
)

// 创建新的回复
func CreateAnswer(c *gin.Context) {
	// 从上下文获取用户名
	username, exists := c.Get("username")
	if !exists {
		utils.ResponseFail(c, "用户未认证", http.StatusUnauthorized)
		return
	}
	// 根据用户名获取用户ID
	userID, err := dao.GetUserIDByUsername(username.(string))
	if err != nil || userID == 0 {
		utils.ResponseFail(c, "用户不存在", http.StatusUnauthorized)
		return
	}
	// 获取请求ID
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ResponseFail(c, "无效的问题ID", http.StatusBadRequest)
		return
	}
	content := c.PostForm("content")
	if content == "" {
		utils.ResponseFail(c, "回答内容不能为空", http.StatusBadRequest)
		return
	}
	err = service.AddAnswer(userID, id, content)
	if err != nil {
		utils.ResponseFail(c, err.Error(), http.StatusBadGateway)
		return
	}
	utils.ResponseSuccess(c, "添加回复成功", http.StatusOK)
}

// 搜索回复,分页查询
func SearchAnswers(c *gin.Context) {
	// 从上下文获取用户名
	username, exists := c.Get("username")
	if !exists {
		utils.ResponseFail(c, "用户未认证", http.StatusUnauthorized)
		return
	}
	// 根据用户名获取用户ID
	userID, err := dao.GetUserIDByUsername(username.(string))
	if err != nil || userID == 0 {
		utils.ResponseFail(c, "用户不存在", http.StatusUnauthorized)
		return
	}
	//获取查询的问题ID
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ResponseFail(c, "无效的 ID", http.StatusBadRequest)
		return
	}
	//获取查询标志
	pattern := c.PostForm("pattern")
	if pattern == "" {
		utils.ResponseFail(c, "查询标志为空", http.StatusBadRequest)
		return
	}
	//分页
	pageStr := c.DefaultPostForm("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		utils.ResponseFail(c, "无效的页码", http.StatusBadRequest)
		return
	}
	//每页记录数
	pageSize := 5
	answers, err := service.SearchAnswersByPattern(id, pattern, page, pageSize)
	if err != nil {
		utils.ResponseFail(c, "搜索错误", http.StatusInternalServerError)
		return
	}
	if len(answers) == 0 {
		utils.ResponseFail(c, "搜索结果为空", http.StatusBadRequest)
		return
	}
	utils.ResponseSuccess(c, answers, http.StatusOK)
}
