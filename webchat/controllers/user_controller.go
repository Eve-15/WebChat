package controllers

import (
	"github.com/Eve-15/GoProjects/webchat/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 用户上线
func UserOnline(c *gin.Context) {
	var json struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := services.Online(c, json.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to establish WebSocket connection"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id": user.ID,
		"name":    user.Name,
	})
}

// 用户下线
func UserOffline(c *gin.Context) {
	userID := c.Param("id")
	if err := services.Offline(userID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "user logged out"})
}

// 获取在线用户列表
func GetOnlineUsers(c *gin.Context) {
	users := services.GetOnlineUsers()
	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}
