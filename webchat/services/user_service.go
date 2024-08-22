package services

import (
	"errors"
	"github.com/Eve-15/GoProjects/webchat/models"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

// WebSocket upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 用户上线（需要通过身份验证）
func Online(c *gin.Context, user *models.User) (*models.User, error) {
	// 升级 HTTP 连接为 WebSocket 连接
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return nil, err
	}

	// 设置 WebSocket 连接
	user.Conn = conn

	models.Mu.Lock()
	models.OnlineMap[user.ID] = user
	models.Mu.Unlock()

	// 广播用户加入消息
	joinMessage := &models.Message{
		Sender:  user.Username,
		Content: "joined the chat",
		Type:    "join",
	}
	BroadcastMessage(joinMessage, user)

	// 启动消息处理
	go HandleMessages(user)

	return user, nil
}

// 用户下线
func Offline(userID string) error {
	models.Mu.Lock()
	defer models.Mu.Unlock()

	user, exists := models.OnlineMap[userID]
	if !exists {
		return errors.New("user not found")
	}

	user.Disconnect()
	return nil
}

// 获取在线用户列表
func GetOnlineUsers() []string {
	models.Mu.Lock()
	defer models.Mu.Unlock()

	var users []string
	for _, user := range models.OnlineMap {
		users = append(users, user.Username)
	}
	return users
}
