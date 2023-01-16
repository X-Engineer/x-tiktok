package main

import (
	"github.com/gin-gonic/gin"
	"x-tiktok/service"
)

func main() {
	// init 操作
	// 启动 websocket 服务
	go service.RunMessageServer()

	r := gin.Default()

	initRouter(r)

	r.Run()
}
