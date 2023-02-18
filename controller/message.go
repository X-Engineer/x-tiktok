package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
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
	preMsgTime := c.Query("pre_msg_time")
	log.Println("preMsgTime", preMsgTime)
	covPreMsgTime, err := strconv.ParseInt(preMsgTime, 10, 64)
	if err != nil {
		log.Println("preMsgTime 参数错误")
		return
	}
	latestTime := time.Unix(covPreMsgTime, 0)
	targetUserId, err := strconv.ParseInt(toUserId, 10, 64)
	if err != nil {
		log.Println("toUserId 参数错误")
		return
	}
	messageService := service.GetMessageServiceInstance()
	messages, err := messageService.MessageChat(loginUserId, targetUserId, latestTime)
	log.Println(messages)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
	} else {
		c.JSON(http.StatusOK, ChatResponse{Response: Response{StatusCode: 0, StatusMsg: "获取消息成功"}, MessageList: messages})
	}
}
