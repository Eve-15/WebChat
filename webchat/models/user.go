package models

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"sync"
)

type User struct {
	ID       string          `json:"id"`
	Username string          `json:"username"`
	Password string          `json:"password"`
	Conn     *websocket.Conn `json:"-"`
	Channel  chan []byte     `json:"-"`
}

// 全局变量
var (
	OnlineMap = make(map[string]*User) // 在线用户映射
	UserMap   = make(map[string]*User) // 所有注册用户映射（简化示例，实际应使用数据库）
	Mu        = &sync.Mutex{}          // 互斥锁，用于保护OnlineMap和UserMap的并发访问
)

func NewUser(username, password string) *User {
	return &User{
		ID:       uuid.New().String(),
		Username: username,
		Password: password,
		Channel:  make(chan []byte),
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
	leaveMessage := &Message{
		Sender:  u.Username,
		Content: "left the chat",
		Type:    "leave",
	}

	Mu.Lock()
	for _, user := range OnlineMap {
		user.WriteMessage(leaveMessage)
	}
	delete(OnlineMap, u.ID)
	Mu.Unlock()

	u.CloseChannel()
	u.Conn.Close()
}
