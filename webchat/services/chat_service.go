package services

import (
	"github.com/Eve-15/GoProjects/webchat/models"
	"github.com/gorilla/websocket"
	"sync"
)

var (
	OnlineMap = make(map[string]*models.User)
	mu        sync.Mutex
)

func CreateUser(conn *websocket.Conn) *models.User {
	user := models.NewUser(conn.RemoteAddr().String(), conn)
	mu.Lock()
	OnlineMap[user.Name] = user
	mu.Unlock()

	joinMessage := &models.Message{
		Sender:  user.Name,
		Content: "joined the chat",
		Type:    "join",
	}
	BroadcastMessage(joinMessage, user)

	go user.ReadPump()
	go user.WritePump()

	return user
}
