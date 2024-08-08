package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"qasociety/mail"
	"qasociety/service"
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

	err := service.RegisterUser(username, password, email)
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
	token, err := service.LoginUser(username, password, email)
	if err != nil {
		utils.ResponseFail(c, err.Error(), http.StatusBadRequest)
		return
	}
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
		err := mail.SendEmailCode(code, email)
		if err != nil {
			utils.ResponseFail(c, err.Error(), http.StatusBadRequest)
			return
		}
		utils.ResponseSuccess(c, "验证码已发送,3分钟内有效", http.StatusOK)
	} else {
		expireTime, err := service.GetExpireTime(user)
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
	//验证码
	code := c.PostForm("code")
	if email == "" || code == "" {
		utils.ResponseFail(c, "邮箱或验证码为空", http.StatusBadRequest)
		return
	}
	flag, err := service.VerifyCode(email, code)
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
