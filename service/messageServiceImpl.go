package service

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
	"x-tiktok/config"
	"x-tiktok/dao"
)

type MessageServiceImpl struct {
}

var (
	messageServiceImpl *MessageServiceImpl
	messageServiceOnce sync.Once
)

// GetMessageServiceInstance Go 单例模式
func GetMessageServiceInstance() *MessageServiceImpl {
	messageServiceOnce.Do(func() {
		messageServiceImpl = &MessageServiceImpl{}
	})
	return messageServiceImpl
}

func (messageService *MessageServiceImpl) SendMessage(fromUserId int64, toUserId int64, content string, actionType int64) error {
	var err error
	switch actionType {
	// actionType = 1 发送消息
	case 1:
		err = dao.SendMessage(fromUserId, toUserId, content, actionType)
	default:
		log.Println(fmt.Sprintf("未定义 actionType=%d", actionType))
		return errors.New(fmt.Sprintf("未定义 actionType=%d", actionType))
	}
	return err
}

func (messageService *MessageServiceImpl) MessageChat(loginUserId int64, targetUserId int64, latestTime time.Time) ([]Message, error) {
	messages := make([]Message, 0, config.VIDEO_INIT_NUM_PER_AUTHOR)
	plainMessages, err := dao.MessageChat(loginUserId, targetUserId, latestTime)
	if err != nil {
		log.Println("MessageChat Service:", err)
		return nil, err
	}
	err = messageService.getRespMessage(&messages, &plainMessages)
	if err != nil {
		log.Println("getRespMessage:", err)
		return nil, err
	}
	return messages, nil
}

func (messageService *MessageServiceImpl) LatestMessage(loginUserId int64, targetUserId int64) (LatestMessage, error) {
	plainMessage, err := dao.LatestMessage(loginUserId, targetUserId)
	if err != nil {
		log.Println("LatestMessage Service:", err)
		return LatestMessage{}, err
	}
	var latestMessage LatestMessage
	latestMessage.message = plainMessage.MsgContent
	if plainMessage.UserId == loginUserId {
		// 最新一条消息是当前登录用户发送的
		latestMessage.msgType = 1
	} else {
		// 最新一条消息是当前好友发送的
		latestMessage.msgType = 0
	}
	return latestMessage, nil
}

// 返回 message list 接口所需的 Message 结构体
func (messageService *MessageServiceImpl) getRespMessage(messages *[]Message, plainMessages *[]dao.Message) error {
	for _, tmpMessage := range *plainMessages {
		var message Message
		messageService.combineMessage(&message, &tmpMessage)
		*messages = append(*messages, message)
	}
	return nil
}

func (messageService *MessageServiceImpl) combineMessage(message *Message, plainMessage *dao.Message) error {
	message.Id = plainMessage.Id
	message.UserId = plainMessage.UserId
	message.ReceiverId = plainMessage.ReceiverId
	message.MsgContent = plainMessage.MsgContent
	message.CreatedAt = plainMessage.CreatedAt.Unix()
	return nil
}
