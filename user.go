package websocket

import (
	"errors"
	"fmt"
	"github.com/link1st/gowebsocket/v2/lib/cache"
	"github.com/link1st/gowebsocket/v2/servers/grpcclient"
	"time"
)

// UserList 查询所有用户
func UserList(appID uint32) (userList []string) {
	userList = make([]string, 0)
	currentTime := uint64(time.Now().Unix())
	servers, err := cache.GetServerAll(currentTime)
	if err != nil {
		fmt.Println("给全体用户发消息", err)
		return
	}
	for _, server := range servers {
		var (
			list []string
		)
		if IsLocal(server) {
			list = GetUserList(appID)
		} else {
			list, _ = grpcclient.GetUserList(server, appID)
		}
		userList = append(userList, list...)
	}
	return
}

// CheckUserOnline 查询用户是否在线
func CheckUserOnline(appID uint32, userID string) (online bool) {
	// 全平台查询
	if appID == 0 {
		for _, appID := range GetAppIDs() {
			online, _ = checkUserOnline(appID, userID)
			if online == true {
				break
			}
		}
	} else {
		online, _ = checkUserOnline(appID, userID)
	}
	return
}