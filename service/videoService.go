package service

import (
	"mime/multipart"
	"x-tiktok/dao"
)

// Video 返回给 Controller 层的 Video 结构体
type Video struct {
	dao.Video
	Author        User  `json:"author"`
	FavoriteCount int64 `json:"favorite_count"`
	CommentCount  int64 `json:"comment_count"`
	IsFavorite    int64 `json:"is_favorite"`
}

type VideoService interface {
	// Publish 将传入的视频流保存到 OSS 中，并在数据库中添加记录
	Publish(data *multipart.FileHeader, title string, userId int64) error
}
