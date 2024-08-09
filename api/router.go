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
		// 请求重置密码
		userGroup.POST("/request_password_reset", RequestPasswordReset)
		//执行重置密码
		userGroup.POST("/reset_password", ResetPassword)
	}
	// 问题相关路由
	questionGroup := r.Group("/question")
	{
		// 需要JWT认证的路由
		questionGroup.Use(middleware.JWTAuthMiddleware())
		// 创建问题
		questionGroup.POST("/create", CreateQuestion)
		// 分页查询获取热门问题列表
		questionGroup.GET("/topic10", TopQuestions)
		// 获取指定问题
		questionGroup.GET("/:id", GetQuestionByID)
		// 更新问题
		questionGroup.PUT("/:id", UpdateQuestion)
		// 删除问题
		questionGroup.DELETE("/:id", DeleteQuestion)
		//添加正常的获取问题列表功能
		questionGroup.GET("/", ListQuestions)
		// 回答相关路由
		answerGroup := questionGroup.Group("/:id/answer")
		{
			answerGroup.Use(middleware.JWTAuthMiddleware())
			// 创建回答
			answerGroup.POST("/create", CreateAnswer)
			//	// 获取回答列表
			//直接用SearchAnswers但是不传入参数pattern就可以实现Answer列表
			//	answerGroup.GET("/", ListAnswers)
			// 分页搜索回复
			answerGroup.GET("/search", SearchAnswers)
			// 删除回答
			answerGroup.DELETE("/:answer_id", DeleteAnswer)
			// 评论相关路由
			commentGroup := answerGroup.Group("/:answer_id/comment")
			{
				// 创建评论
				commentGroup.POST("/create", CreateComment)
				// 分页查询获取评论列表
				commentGroup.GET("/", ListComments)
				// 删除评论
				commentGroup.DELETE("/:comment_id", DeleteComment)
			}
		}
	}
	r.Run(":8080")
}
