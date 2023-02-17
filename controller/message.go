package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
	"x-tiktok/config"
	"x-tiktok/service"
)

type ChatResponse struct {
	Response
	MessageList []service.Message `json:"message_list"`
}

// MessageAction 发送消息
func MessageAction(c *gin.Context) {
	toUserId := c.Query("to_user_id")
	content := c.Query("content")
	actionType := c.Query("action_type")
	loginUserId := c.GetInt64("userId")
	targetUserId, err := strconv.ParseInt(toUserId, 10, 64)
	targetActionType, err1 := strconv.ParseInt(actionType, 10, 64)
	if err != nil || err1 != nil {
		log.Println("toUserId/actionType 参数错误")
		return
	}
	messageService := service.GetMessageServiceInstance()
	err = messageService.SendMessage(loginUserId, targetUserId, content, targetActionType)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Send Message 接口错误"})
	}
	c.JSON(http.StatusOK, Response{StatusCode: 0})
}

// MessageChat 消息列表
func MessageChat(c *gin.Context) {
	loginUserId := c.GetInt64("userId")
	toUserId := c.Query("to_user_id")
	// 首先判断 key 是否存在，不存在则设置初始 latestTime 为 1970
	latestTime, exist := config.LatestRequestTime[fmt.Sprintf("%d-%s", loginUserId, toUserId)]
	if exist != true {
		latestTime = time.Unix(0, 0)
	}
	targetUserId, err := strconv.ParseInt(toUserId, 10, 64)
	if err != nil {
		log.Println("toUserId 参数错误")
		return
	}
	//log.Println("loginUserId", loginUserId)
	//log.Println("to_user_id:", toUserId)
	messageService := service.GetMessageServiceInstance()
	messages, err := messageService.MessageChat(loginUserId, targetUserId, latestTime)
	// 如果聊天记录不为空，则将最新的一条聊天记录时间作为下次请求的时间
	if len(messages) != 0 {
		config.LatestRequestTime[fmt.Sprintf("%d-%s", loginUserId, toUserId)] = time.Unix(messages[len(messages)-1].CreatedAt, 0)
	}
	log.Println(messages)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
	} else {
		c.JSON(http.StatusOK, ChatResponse{Response: Response{StatusCode: 0, StatusMsg: "获取消息成功"}, MessageList: messages})
	}
}
