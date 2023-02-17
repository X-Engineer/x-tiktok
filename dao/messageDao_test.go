package dao

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func TestSendMessage(t *testing.T) {
	fromUserId := 7
	toUserId := 2
	err := SendMessage(int64(fromUserId), int64(toUserId), fmt.Sprintf("我是 userId=%d,发送消息给 userId=%d", fromUserId, toUserId), 1)
	if err == nil {
		log.Println("SendMessage 测试成功！")
	}
}

func TestMessageChat(t *testing.T) {
	loginUserId := 7
	targetUserId := 1
	messages, err := MessageChat(int64(loginUserId), int64(targetUserId), time.Now())
	if err != nil {
		log.Println("MessageChat 测试失败")
	}
	for _, msg := range messages {
		log.Println(fmt.Sprintf("%d -> %d: %s (sendTime:%v)", msg.UserId, msg.ReceiverId, msg.MsgContent, msg.CreatedAt))
	}
}

func TestLatestMessage(t *testing.T) {
	loginUserId := 2
	targetUserId := 7
	message, err := LatestMessage(int64(loginUserId), int64(targetUserId))
	if err != nil {
		log.Println("LatestMessage 测试失败")
	}
	log.Println(fmt.Sprintf("%d -> %d 的最新一条消息记录：%s", message.UserId, message.ReceiverId, message.MsgContent))
}
