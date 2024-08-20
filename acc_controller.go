// Package websocket 处理
package websocket

import (
	"encoding/json"
	"fmt"
	"github.com/link1st/gowebsocket/v2/common"
	"github.com/link1st/gowebsocket/v2/models"
	"time"
)

// PingController 简化版 ping
func PingController(client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {
	code = common.OK
	fmt.Println("webSocket_request ping接口", client.Addr, seq, message)
	data = "pong"
	return
}

// LoginController 用户登录
func LoginController(client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {
	code = common.OK
	fmt.Println("webSocket_request 用户登录", seq, string(message))

	// 直接处理登录请求，不做用户ID验证
	request := &models.Login{}
	if err := json.Unmarshal(message, request); err != nil {
		code = common.ParameterIllegal
		fmt.Println("用户登录 解析数据失败", seq, err)
		return
	}

	// 模拟登录成功，不做实际验证
	client.Login(request.AppID, request.UserID, uint64(time.Now().Unix()))
	fmt.Println("用户登录 成功", seq, client.Addr, request.UserID)

	return
}

// HeartbeatController 心跳接口
func HeartbeatController(client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {
	code = common.OK
	fmt.Println("webSocket_request 心跳接口", client.AppID, client.UserID)

	// 直接处理心跳请求，不做登录检查
	request := &models.HeartBeat{}
	if err := json.Unmarshal(message, request); err != nil {
		code = common.ParameterIllegal
		fmt.Println("心跳接口 解析数据失败", seq, err)
		return
	}

	// 更新心跳时间
	client.Heartbeat(uint64(time.Now().Unix()))
	fmt.Println("心跳接口 成功（简化版）", seq, client.AppID, client.UserID)

	return
}
