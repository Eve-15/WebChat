package main

import (
	"github.com/Eve-15/GoProjects/webchat/routers" // 导入路由
	"github.com/gin-gonic/gin"                     // 导入gin框架
)

func main() {
	r := gin.Default()

	// 设置路由
	routes.SetupRoutes(r)

	// 启动服务器
	r.Run(":8080")
}
