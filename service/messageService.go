package service

import "time"

type Message struct {
	Id         int64  `json:"id"`
	UserId     int64  `json:"from_user_id"`
	ReceiverId int64  `json:"to_user_id"`
	MsgContent string `json:"content"`
	CreatedAt  int64  `json:"create_time"`
}

// LatestMessage 提供给用户好友列表接口的最新一条聊天信息, msgType 消息类型，0 => 当前请求用户接收的消息， 1 => 当前请求用户发送的消息
type LatestMessage struct {
	message string `json:"message"`
	msgType int64  `json:"msg_type"`
}

type MessageService interface {
	// SendMessage 发送消息服务
	SendMessage(fromUserId int64, toUserId int64, content string, actionType int64) error

	// MessageChat 聊天记录服务，注意返回的 Message 结构体字段与 Dao 层的不完全相同
	MessageChat(loginUserId int64, targetUserId int64, latestTime time.Time) ([]Message, error)

	// LatestMessage 返回两个 loginUserId 和好友 targetUserId 最近的一条聊天记录
	LatestMessage(loginUserId int64, targetUserId int64) (LatestMessage, error)
}
