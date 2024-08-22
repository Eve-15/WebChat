package routes

import (
	"github.com/Eve-15/GoProjects/webchat/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	apiGroup := r.Group("/api")
	{
		apiGroup.POST("/register", controllers.Register)       // 用户注册
		apiGroup.POST("/login", controllers.Login)             // 用户登录
		apiGroup.POST("/online", controllers.UserOnline)       // 用户上线接口
		apiGroup.DELETE("/users/:id", controllers.UserOffline) // 用户下线接口
		apiGroup.GET("/users", controllers.GetOnlineUsers)     // 获取在线用户列表
	}
}
