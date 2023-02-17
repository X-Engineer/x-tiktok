package service

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func TestMessageServiceImpl_SendMessage(t *testing.T) {
	err := messageServiceImpl.SendMessage(2, 8, "2 号用户发送消息给 8 号用户", 1)
	if err == nil {
		log.Println("SendMessage Service 正常")
	}
}

func TestMessageServiceImpl_MessageChat(t *testing.T) {
	chat, _ := messageServiceImpl.MessageChat(8, 2, time.Now())
	for _, msg := range chat {
		log.Println(fmt.Sprintf("%+v", msg))
	}
}

func TestMessageServiceImpl_LatestMessage(t *testing.T) {
	message, _ := messageServiceImpl.LatestMessage(7, 1)
	log.Println(fmt.Sprintf("%+v", message))
}
