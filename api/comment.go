package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"qasociety/service"
	"qasociety/service/dao"
	"qasociety/utils"
	"strconv"
)

// CreateComment 创建评论
func CreateComment(c *gin.Context) {
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
	// 获取请求ID和评论对应的AnswerID
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	question, _ := dao.GetQuestionByID(id)
	if question == nil {
		utils.ResponseFail(c, "无效的问题ID", http.StatusNotFound)
		return
	}
	// 获取回复
	answerIDStr := c.Param("answer_id")
	answerID, err := strconv.Atoi(answerIDStr)
	answer, err := dao.GetAnswerByID(answerID)
	if err != nil {
		utils.ResponseFail(c, "获取回答信息失败", http.StatusInternalServerError)
		return
	}
	if answer == nil || answer.QuestionID != id {
		utils.ResponseFail(c, "回答不属于该问题", http.StatusBadRequest)
		return
	}
	// 获取评论内容
	content := c.PostForm("content")
	if content == "" {
		utils.ResponseFail(c, "评论内容不能为空", http.StatusBadRequest)
		return
	}
	// 插入评论数据
	err = service.CreateComment(userID, answerID, content)
	if err != nil {
		utils.ResponseFail(c, "创建评论失败", http.StatusInternalServerError)
		return
	}
	utils.ResponseSuccess(c, "评论创建成功", http.StatusOK)
}

// ListComments 分页查询获取评论列表
func ListComments(c *gin.Context) {
	_, exists := c.Get("username")
	if !exists {
		utils.ResponseFail(c, "用户未认证", http.StatusUnauthorized)
		return
	}
	// 获取请求ID和评论对应的AnswerID
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)
	question, _ := dao.GetQuestionByID(id)
	if question == nil {
		utils.ResponseFail(c, "无效的问题ID", http.StatusNotFound)
		return
	}
	//获取请求连接中的answerID
	answerIDStr := c.Param("answer_id")
	answerID, err := strconv.Atoi(answerIDStr)
	answer, err := dao.GetAnswerByID(answerID)
	if err != nil {
		utils.ResponseFail(c, "获取回答信息失败", http.StatusInternalServerError)
		return
	}
	if answer == nil || answer.QuestionID != id {
		utils.ResponseFail(c, "回答不属于该问题", http.StatusBadRequest)
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
	comments, err := service.GetCommentsByAnswerID(answerID, page, pageSize)
	if err != nil {
		utils.ResponseFail(c, "搜索错误", http.StatusInternalServerError)
		return
	}
	if len(comments) == 0 {
		utils.ResponseFail(c, "搜索结果为空", http.StatusBadRequest)
		return
	}
	utils.ResponseSuccess(c, comments, http.StatusOK)
}
