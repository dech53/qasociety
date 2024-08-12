package api

import (
	"context"
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
func LikeAnswer(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		utils.ResponseFail(c, err.Error(), http.StatusUnauthorized)
		return
	}
	AnswerIdStr := c.Param("answer_id")
	answerId, err := strconv.Atoi(AnswerIdStr)
	if err != nil {
		utils.ResponseFail(c, "不合法的回复ID", http.StatusBadRequest)
		return
	}
	//检查回复ID是否存在
	_, err = dao.GetAnswerByID(answerId)
	if err != nil {
		utils.ResponseFail(c, err.Error(), http.StatusBadGateway)
		return
	}
	// 3. 检查用户是否已经点赞过该回答
	if isLiked, err := utils.IsUserLikedAnswer(userID, answerId); err != nil {
		utils.ResponseFail(c, err.Error(), http.StatusInternalServerError)
		return
	} else if isLiked {
		err = UnlikeAnswer(c, userID, answerId)
		if err != nil {
			utils.ResponseFail(c, err.Error(), http.StatusInternalServerError)
			return
		}
		utils.ResponseSuccess(c, "取消点赞成功", http.StatusOK)
		return
	}
	err = utils.PublishLikeEvent(userID, answerId)
	if err != nil {
		utils.ResponseFail(c, "点赞失败", http.StatusInternalServerError)
		return
	}
	utils.ResponseSuccess(c, "点赞成功", http.StatusOK)
}
func UnlikeAnswer(c *gin.Context, userID, answerId int) error {
	redisKey := "answer:likes:" + strconv.Itoa(answerId)
	// 检查 userID 是否是 Redis Set 的成员
	isMember, err := dao.Rdb.SIsMember(context.Background(), redisKey, userID).Result()
	if err != nil {
		// 如果 Redis 操作失败，返回错误
		utils.ResponseFail(c, "Error checking like status in Redis", http.StatusInternalServerError)
		return err
	}
	if isMember {
		// 如果在 Redis 中找到了记录，从 Set 中移除
		err = dao.DeleteThumbRedis(redisKey, userID)
		if err != nil {
			return err
		}
	} else {
		err = dao.DeleteThumbMysql(answerId, userID)
		if err != nil {
			return err
		}
	}
	return nil
}
func GetAnswerLikesCount(c *gin.Context) {
	answerIdStr := c.Param("answer_id")
	answerId, err := strconv.Atoi(answerIdStr)
	if err != nil {
		utils.ResponseFail(c, err.Error(), http.StatusBadRequest)
		return
	}
	counts, err := dao.GetAnswerLikesCount(answerId)
	if err != nil {
		utils.ResponseFail(c, err.Error(), http.StatusBadRequest)
		return
	}
	utils.ResponseSuccess(c, counts, http.StatusOK)
}
