package services

import (
	"errors"
	"github.com/Eve-15/WebChat/webchat/models"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidPassword = errors.New("invalid password")
)

// Register 注册用户
func Register(username, password string) (*models.User, error) {
	models.Mu.Lock()
	defer models.Mu.Unlock()

	if _, exists := models.UserMap[username]; exists {
		return nil, errors.New("username already taken")
	}

	user := models.NewUser(username, password)
	models.UserMap[username] = user
	return user, nil
}

// Authenticate 验证用户登录
func Authenticate(username, password string) (*models.User, error) {
	models.Mu.Lock()
	defer models.Mu.Unlock()

	user, exists := models.UserMap[username]
	if !exists {
		return nil, ErrUserNotFound
	}

	if user.Password != password {
		return nil, ErrInvalidPassword
	}

	return user, nil
}
