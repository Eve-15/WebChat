package models

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"sync"
)

// 在线用户的全局映射
var (
	OnlineMap = make(map[string]*User)
	Mu        sync.Mutex
)

type User struct {
	ID      string
	Name    string
	Conn    *websocket.Conn
	Channel chan []byte
}

func NewUser(name string, conn *websocket.Conn) *User {
	return &User{
		ID:      uuid.New().String(), // 生成唯一ID
		Name:    name,
		Conn:    conn,
		Channel: make(chan []byte),
	}
}

func (u *User) WriteMessage(message *Message) {
	jsonMessage, _ := json.Marshal(message)
	u.Channel <- jsonMessage
}

func (u *User) ReadPump() {
	defer u.Conn.Close()
	for {
		_, message, err := u.Conn.ReadMessage()
		if err != nil {
			break
		}
		u.Channel <- message
	}
}

func (u *User) WritePump() {
	defer u.Conn.Close()
	for message := range u.Channel {
		u.Conn.WriteMessage(websocket.TextMessage, message)
	}
}

func (u *User) CloseChannel() {
	close(u.Channel)
}

// Disconnect 方法，处理用户断开连接的逻辑
func (u *User) Disconnect() {
	// 创建离开消息广播给其他用户
	leaveMessage := &Message{
		Sender:  u.Name,
		Content: "left the chat",
		Type:    "leave",
	}

	// 广播离开消息
	Mu.Lock()
	for _, user := range OnlineMap {
		user.WriteMessage(leaveMessage)
	}
	// 从在线用户列表中删除自己
	delete(OnlineMap, u.ID)
	Mu.Unlock()

	// 关闭用户的消息通道
	u.CloseChannel()

	// 关闭 WebSocket 连接
	u.Conn.Close()
}
