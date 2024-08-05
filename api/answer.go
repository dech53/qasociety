package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"qasociety/service"
	"qasociety/service/dao"
	"qasociety/utils"
	"strconv"
)

// CreateAnswer 创建新的回复
func CreateAnswer(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		utils.ResponseFail(c, err.Error(), http.StatusUnauthorized)
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
	//通过问题ID更新question数据的更新时间
	question, err := dao.GetQuestionByID(id)
	if err != nil {
		utils.ResponseFail(c, err.Error(), http.StatusBadGateway)
		return
	}
	err = dao.UpdateQuestion(question.ID, question.Title, question.Content)
	utils.ResponseSuccess(c, "添加回复成功", http.StatusOK)
}

// SearchAnswers 搜索回复,分页查询
func SearchAnswers(c *gin.Context) {
	_, err := utils.GetUserID(c)
	if err != nil {
		utils.ResponseFail(c, err.Error(), http.StatusUnauthorized)
		return
	}
	//获取查询的问题ID
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ResponseFail(c, "无效的 ID", http.StatusBadRequest)
		return
	}
	//获取查询标志,默认值为空
	pattern := c.DefaultPostForm("pattern", "")
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

// DeleteAnswer 删除回答
func DeleteAnswer(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		utils.ResponseFail(c, err.Error(), http.StatusUnauthorized)
		return
	}
	AnswerIdStr := c.Param("answer_id")
	answerId, err := strconv.Atoi(AnswerIdStr)
	if err != nil {
		utils.ResponseFail(c, "无效的ID", http.StatusBadRequest)
		return
	}
	answer, err := dao.GetAnswerByID(answerId)
	if err != nil {
		utils.ResponseFail(c, err.Error(), http.StatusBadGateway)
		return
	}
	if answer.UserID != userID {
		utils.ResponseFail(c, "无权删除回复", http.StatusUnauthorized)
		return
	}
	err = service.DeleteAnswer(answer)
	if err != nil {
		utils.ResponseFail(c, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.ResponseSuccess(c, "删除成功", http.StatusOK)
}
func DeleteAnswers(c *gin.Context) {

}
