package services

import (
	"encoding/json"
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

// 用户上线
func Online(c *gin.Context, name string) (*models.User, error) {
	// 升级 HTTP 连接为 WebSocket 连接
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return nil, err
	}

	// 创建新用户
	user := models.NewUser(name, conn)

	models.Mu.Lock()
	models.OnlineMap[user.ID] = user
	models.Mu.Unlock()

	// 广播用户加入消息
	joinMessage := &models.Message{
		Sender:  user.Name,
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
		users = append(users, user.Name)
	}
	return users
}

// 广播消息给所有在线用户
func BroadcastMessage(message *models.Message, sender *models.User) {
	models.Mu.Lock()
	defer models.Mu.Unlock()

	for _, u := range models.OnlineMap {
		if u.ID != sender.ID { // 不要给自己发送消息
			u.WriteMessage(message)
		}
	}
}

// 处理用户消息
func HandleMessages(user *models.User) {
	defer user.Disconnect()

	for {
		select {
		case message, ok := <-user.Channel:
			if !ok {
				return
			}
			var msg models.Message
			if err := json.Unmarshal(message, &msg); err != nil {
				continue
			}
			BroadcastMessage(&msg, user)
		}
	}
}
