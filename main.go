package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	// init 操作
	// 启动 websocket 服务
	//go RunMessageServer()

	r := gin.Default()

	initRouter(r)

	r.Run()
}
