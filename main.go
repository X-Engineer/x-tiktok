package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	// init 操作
	r := gin.Default()

	initMiddleware()
	initRouter(r)

	r.Run()
}
