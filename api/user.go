package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"qasociety/mail"
	"qasociety/service"
	"qasociety/service/dao"
	"qasociety/utils"
	"strconv"
)

// Register 用户注册
func Register(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	email := c.PostForm("email")
	if username == "" || password == "" || email == "" {
		utils.ResponseFail(c, "有字段为空", http.StatusBadRequest)
		return
	}
	err := utils.MatchStr(password)
	if err != nil {
		utils.ResponseFail(c, err.Error(), http.StatusBadRequest)
		return
	}
	err = service.RegisterUser(username, password, email)
	if err != nil {
		utils.ResponseFail(c, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.ResponseSuccess(c, "用户注册成功", http.StatusOK)
}

// Login 用户登录
func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	email := c.PostForm("email")
	if (username == "" && password == "") || (password == "" && email == "") {
		utils.ResponseFail(c, "有字段为空", http.StatusBadRequest)
		return
	}
	info := utils.GetUserAgent(c)
	// 构建 Redis 键
	redisKey := "session:" + username + ":" + info
	// 尝试从 Redis 中获取现有的 token
	cachedToken, err := dao.Rdb.Get(context.Background(), redisKey).Result()
	if err == nil { // Redis 中存在 token
		utils.ResponseSuccess(c, cachedToken, http.StatusOK)
		return
	}
	// Redis 中不存在 token，进行用户验证和生成新 token
	token, err := service.LoginUser(username, password, email, info, "")
	if err != nil {
		utils.ResponseFail(c, err.Error(), http.StatusBadRequest)
		return
	}
	// 返回新生成的 token
	utils.ResponseSuccess(c, token, http.StatusOK)
}

// RequestPasswordReset 发请求,重置密码
func RequestPasswordReset(c *gin.Context) {
	email := c.DefaultPostForm("email", "")
	user, err := service.GetUserByPattern("email", email)
	if err != nil {
		utils.ResponseFail(c, err.Error(), http.StatusBadRequest)
		return
	}
	code := utils.GenerateCode()
	flag, err := service.ResetRequest(code, user)
	if err != nil {
		utils.ResponseFail(c, err.Error(), http.StatusBadRequest)
		return
	}
	if flag {
		err := mail.SendEmailCode(code, email, "重置密码")
		if err != nil {
			utils.ResponseFail(c, err.Error(), http.StatusBadRequest)
			return
		}
		utils.ResponseSuccess(c, "验证码已发送,3分钟内有效", http.StatusOK)
	} else {
		expireTime, err := service.GetExpireTime(user, "resetPassword")
		if err != nil {
			utils.ResponseFail(c, err.Error(), http.StatusBadRequest)
			return
		}
		utils.ResponseSuccess(c, "请等待"+strconv.FormatInt(int64(expireTime.Seconds()), 10)+"秒后重试", http.StatusOK)
	}
}

// ResetPassword 执行重置密码
func ResetPassword(c *gin.Context) {
	email := c.DefaultPostForm("email", "")
	newPassword := c.DefaultPostForm("newPassword", "")
	//正则判断
	err := utils.MatchStr(newPassword)
	if err != nil {
		utils.ResponseFail(c, err.Error(), http.StatusBadRequest)
		return
	}
	//验证码
	code := c.PostForm("code")
	if email == "" || code == "" {
		utils.ResponseFail(c, "邮箱或验证码为空", http.StatusBadRequest)
		return
	}
	flag, err := service.VerifyCode(email, code, "resetPassword")
	if err != nil {
		utils.ResponseFail(c, err.Error(), http.StatusBadRequest)
		return
	}
	if flag {
		err := service.ResetPassword(email, newPassword)
		if err != nil {
			utils.ResponseFail(c, err.Error(), http.StatusBadRequest)
			return
		}
		utils.ResponseSuccess(c, "密码修改成功", http.StatusOK)
		return
	} else {
		utils.ResponseFail(c, "密码修改失败", http.StatusBadRequest)
		return
	}
}

// LoginByCodeRequest 验证码登录请求
func LoginByCodeRequest(c *gin.Context) {
	email := c.PostForm("email")
	user, err := service.GetUserByPattern("email", email)
	if err != nil {
		utils.ResponseFail(c, err.Error(), http.StatusBadRequest)
		return
	}
	code := utils.GenerateCode()
	flag, err := service.LoginByCodeRequest(code, user)
	if err != nil {
		utils.ResponseFail(c, err.Error(), http.StatusBadRequest)
		return
	}
	if flag {
		err := mail.SendEmailCode(code, email, "登陆验证码")
		if err != nil {
			utils.ResponseFail(c, err.Error(), http.StatusBadRequest)
			return
		}
		utils.ResponseSuccess(c, "验证码已发送,3分钟内有效", http.StatusOK)
	} else {
		expireTime, err := service.GetExpireTime(user, "loginCode")
		if err != nil {
			utils.ResponseFail(c, err.Error(), http.StatusBadRequest)
			return
		}
		utils.ResponseSuccess(c, "请等待"+strconv.FormatInt(int64(expireTime.Seconds()), 10)+"秒后重试", http.StatusOK)
	}
}

// LoginByCode 执行验证码登录
func LoginByCode(c *gin.Context) {
	email := c.PostForm("email")
	code := c.PostForm("code")
	if email == "" || code == "" {
		utils.ResponseFail(c, "邮箱或验证码为空", http.StatusBadRequest)
		return
	}
	user, err := service.GetUserByPattern("email", email)
	flag, err := service.VerifyCode(email, code, "loginCode")
	if err != nil {
		utils.ResponseFail(c, err.Error(), http.StatusBadRequest)
		return
	}
	if flag {
		info := utils.GetUserAgent(c)
		// 构建 Redis 键
		redisKey := "session:" + user.Username + ":" + info
		// 尝试从 Redis 中获取现有的 token
		cachedToken, err := dao.Rdb.Get(context.Background(), redisKey).Result()
		if err == nil { // Redis 中存在 token
			utils.ResponseSuccess(c, cachedToken, http.StatusOK)
			return
		}
		// Redis 中不存在 token，进行用户验证和生成新 token
		token, err := service.LoginUser(user.Username, user.Password, email, info, "withoutVerifyPassword")
		if err != nil {
			utils.ResponseFail(c, err.Error(), http.StatusBadRequest)
			return
		}
		// 返回新生成的 token
		utils.ResponseSuccess(c, token, http.StatusOK)
	}
}
