package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"qasociety/service"
	"qasociety/service/dao"
	"qasociety/utils"
	"strconv"
)

// CreateQuestion 创建问题
func CreateQuestion(c *gin.Context) {
	// 从上下文中获取用户名
	username, exists := c.Get("username")
	//头部验证
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

// GetQuestionByID 获取指定问题
func GetQuestionByID(c *gin.Context) {
	_, err := utils.GetUserID(c)
	if err != nil {
		utils.ResponseFail(c, err.Error(), http.StatusUnauthorized)
	}
	//获取请求ID
	idStr := c.Param("id")
	//ID转int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ResponseFail(c, "无效的问题 ID", http.StatusBadRequest)
		return
	}
	//通过ID获取问题
	question, err := service.GetQuestionByID(id)
	if err != nil {
		utils.ResponseFail(c, "获取问题失败", http.StatusInternalServerError)
		return
	}
	utils.ResponseSuccess(c, question, http.StatusOK)
}

// UpdateQuestion 更新问题
func UpdateQuestion(c *gin.Context) {
	id, err := utils.JudgeID(c)
	if err != nil {
		utils.ResponseFail(c, err.Error(), http.StatusBadRequest)
		return
	}
	title := c.PostForm("title")
	content := c.PostForm("content")
	// 参数验证
	if title == "" || content == "" {
		utils.ResponseFail(c, "标题和内容不能为空", http.StatusBadRequest)
		return
	}
	err = service.UpdateQuestion(id, title, content)
	if err != nil {
		utils.ResponseFail(c, "更新问题失败", http.StatusInternalServerError)
		return
	}
	utils.ResponseSuccess(c, "更新问题成功", http.StatusOK)
}

// DeleteQuestion 删除问题
func DeleteQuestion(c *gin.Context) {
	id, err := utils.JudgeID(c)
	if err != nil {
		utils.ResponseFail(c, err.Error(), http.StatusBadRequest)
		return
	}
	err = service.DeleteQuestion(id)
	if err != nil {
		utils.ResponseFail(c, err, http.StatusInternalServerError)
		return
	}
	utils.ResponseSuccess(c, "删除成功", http.StatusOK)
}

// TopQuestions 分页查询展示问题
//
//	直接从redis热榜中获取
func TopQuestions(c *gin.Context) {
	_, err := utils.GetUserID(c)
	if err != nil {
		utils.ResponseFail(c, err.Error(), http.StatusUnauthorized)
		return
	}
	//获取查询标志,默认值为空
	//pattern := c.PostForm("pattern")
	//分页
	pageStr := c.DefaultPostForm("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		utils.ResponseFail(c, "无效的页码", http.StatusBadRequest)
		return
	}
	//每页记录数
	//pageSize := 5
	//order := c.DefaultPostForm("order", "")
	//  存在效率问题,是否可以实现redis中的分页查询
	//此处无需做修改,只需要修改redis的增添逻辑即可
	questions, err := service.GetQuestionsByRedis()
	if err != nil {
		utils.ResponseFail(c, err.Error(), http.StatusInternalServerError)
		return
	}
	if questions == nil || len(questions) == 0 {
		utils.ResponseFail(c, "问题为空", http.StatusBadRequest)
		return
	}
	utils.ResponseSuccess(c, questions, http.StatusOK)
}

// ListQuestions 常规方式通过update_time排序获取问题
func ListQuestions(c *gin.Context) {
	_, err := utils.GetUserID(c)
	if err != nil {
		utils.ResponseFail(c, err.Error(), http.StatusUnauthorized)
		return
	}
	pageStr := c.DefaultPostForm("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		utils.ResponseFail(c, "无效的页码", http.StatusBadRequest)
		return
	}
	pattern := c.DefaultPostForm("pattern", "")
	pageSize := 5
	order := c.DefaultPostForm("order", "")
	questions, err := service.SearchQuestionsByPattern(pattern, order, page, pageSize)
	if err != nil {
		utils.ResponseFail(c, err.Error(), http.StatusInternalServerError)
		return
	}
	if questions == nil || len(questions) == 0 {
		utils.ResponseFail(c, "问题为空", http.StatusBadRequest)
		return
	}
	utils.ResponseSuccess(c, questions, http.StatusOK)
}
