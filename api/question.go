package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"qasociety/service"
	"qasociety/service/dao"
	"qasociety/utils"
)

// CreateQuestion 创建问题
func CreateQuestion(c *gin.Context) {
	// 从上下文中获取用户名
	username, exists := c.Get("username")
	if !exists {
		utils.ResponseFail(c, "用户未认证", http.StatusUnauthorized)
		return
	}
	// 从请求中获取数据
	title := c.PostForm("title")
	content := c.PostForm("content")
	// 参数验证
	if title == "" || content == "" {
		utils.ResponseFail(c, "标题和内容不能为空", http.StatusBadRequest)
		return
	}
	// 根据用户名获取用户ID
	userID, err := dao.GetUserIDByUsername(username.(string))
	if err != nil || userID == 0 {
		utils.ResponseFail(c, "用户不存在", http.StatusNotFound)
		return
	}
	// 调用服务层创建问题
	err = service.AddQuestion(userID, title, content)
	if err != nil {
		utils.ResponseFail(c, "创建问题失败", http.StatusInternalServerError)
		return
	}
	// 返回成功响应
	utils.ResponseSuccess(c, "问题创建成功", http.StatusOK)
}
