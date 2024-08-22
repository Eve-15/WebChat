package services

import (
	"encoding/json"
	"github.com/Eve-15/GoProjects/webchat/models"
	"strings"
)

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
		_, message, err := user.Conn.ReadMessage()
		if err != nil {
			break
		}

		var msg models.Message
		if err := json.Unmarshal(message, &msg); err != nil {
			continue
		}

		// 根据消息类型处理
		switch msg.Type {
		case "message":
			// 广播普通消息给所有用户
			BroadcastMessage(&msg, user)
		case "private":
			// 处理私聊消息
			handlePrivateMessage(&msg, user)
		case "system":
			// 处理系统消息（可扩展）
			handleSystemMessage(&msg, user)
		default:
			// 其他类型消息处理
		}
	}
}

// 处理私聊消息
func handlePrivateMessage(msg *models.Message, sender *models.User) {
	models.Mu.Lock()
	defer models.Mu.Unlock()

	// 假设msg.Content格式为"targetUserID:messageContent"
	parts := strings.SplitN(msg.Content, ":", 2)
	if len(parts) != 2 {
		return
	}

	targetUserID := parts[0]
	messageContent := parts[1]

	// 查找目标用户
	targetUser, exists := models.OnlineMap[targetUserID]
	if !exists {
		// 目标用户不存在，可能要发送反馈给发送者
		return
	}

	// 构建私聊消息并发送
	privateMessage := &models.Message{
		Sender:  sender.Username,
		Content: messageContent,
		Type:    "private",
	}
	targetUser.WriteMessage(privateMessage)
}

// 处理系统消息（可以根据业务需求扩展）
func handleSystemMessage(msg *models.Message, sender *models.User) {
	// 例如处理一些特殊的系统指令或广播系统消息
	// 这个部分可以根据你的具体业务逻辑进行扩展
}
