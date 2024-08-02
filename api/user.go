package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"qasociety/service"
	"qasociety/utils"
)

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
