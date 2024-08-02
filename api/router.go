package api

import (
	"github.com/gin-gonic/gin"
	"qasociety/api/middleware"
)

// InitRouter 初始化路由
func InitRouter() {
	r := gin.Default()
	// 中间件
	r.Use(middleware.CORS())
	// 用户路由
	userGroup := r.Group("/user")
	{
		// 用户注册
		userGroup.POST("/register", Register)
		// 用户登录
		userGroup.POST("/login", Login)
	}
	// 问题相关路由
	questionGroup := r.Group("/question")
	{
		// 需要JWT认证的路由
		questionGroup.Use(middleware.JWTAuthMiddleware())
		// 创建问题
		questionGroup.POST("/create", CreateQuestion)
		//// 获取问题列表
		//questionGroup.GET("/", ListQuestions)
		//// 获取指定问题
		//questionGroup.GET("/:id", GetQuestion)
		//// 更新问题
		//questionGroup.PUT("/:id", UpdateQuestion)
		//// 删除问题
		//questionGroup.DELETE("/:id", DeleteQuestion)
		//
		//// 回答相关路由
		//answerGroup := questionGroup.Group("/:question_id/answer")
		//{
		//	// 创建回答
		//	answerGroup.POST("/create", CreateAnswer)
		//	// 获取回答列表
		//	answerGroup.GET("/", ListAnswers)
		//	// 获取指定回答
		//	answerGroup.GET("/:id", GetAnswer)
		//	// 更新回答
		//	answerGroup.PUT("/:id", UpdateAnswer)
		//	// 删除回答
		//	answerGroup.DELETE("/:id", DeleteAnswer)
		//	// 评论相关路由
		//	commentGroup := answerGroup.Group("/:answer_id/comment")
		//	{
		//		// 创建评论
		//		commentGroup.POST("/create", CreateComment)
		//		// 获取评论列表
		//		commentGroup.GET("/", ListComments)
		//		// 获取指定评论
		//		commentGroup.GET("/:id", GetComment)
		//		// 更新评论
		//		commentGroup.PUT("/:id", UpdateComment)
		//		// 删除评论
		//		commentGroup.DELETE("/:id", DeleteComment)
		//	}
		//}
	}

	r.Run(":8080")
}
