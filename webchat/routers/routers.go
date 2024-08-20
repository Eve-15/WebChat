package routes

import (
	"github.com/Eve-15/GoProjects/webchat/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	chatGroup := r.Group("/chat")
	{
		chatGroup.GET("/ws", controllers.ChatWebSocket) // WebSocket 连接接口
	}

	apiGroup := r.Group("/api")
	{
		apiGroup.GET("/users", controllers.GetOnlineUsers)     // 获取在线用户列表
		apiGroup.POST("/users", controllers.UserOnline)        // 用户上线接口
		apiGroup.DELETE("/users/:id", controllers.UserOffline) // 用户下线接口
	}
}
