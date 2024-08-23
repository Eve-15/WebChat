package controllers

import (
	"github.com/Eve-15/GoProjects/webchat/models"
	"github.com/Eve-15/GoProjects/webchat/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

// ChatWebSocket 处理 WebSocket 连接
func ChatWebSocket(c *gin.Context) {
	// 获取请求中的 user_id 参数
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	// 查找用户
	user, exists := models.UserMap[userID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// 将用户上线并建立 WebSocket 连接
	_, err := services.Online(c, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to establish WebSocket connection"})
		return
	}

	// 成功上线并连接 WebSocket，返回成功响应
	c.JSON(http.StatusOK, gin.H{"status": "user online and WebSocket connected"})
}
