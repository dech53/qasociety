package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"qasociety/service"
	"qasociety/service/dao"
	"qasociety/utils"
	"strconv"
)

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
