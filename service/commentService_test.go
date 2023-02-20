package service

import (
	"fmt"
	"log"
	"testing"
	"x-tiktok/dao"
	"x-tiktok/middleware/redis"
)

func TestCommentServiceImpl_GetCommentCnt(t *testing.T) {
	redis.InitRedis()
	count, err := commentServiceImpl.GetCommentCnt(25)
	if err != nil {
		log.Default()
	}
	fmt.Println(count)
}

func TestCommentServiceImpl_CommentAction(t *testing.T) {
	redis.InitRedis()
	var comment dao.Comment = dao.Comment{
		UserId:  5,
		VideoId: 14,
		Content: "这条评论来自单元测试TestInsertComment",
	}
	commentRes, err := commentServiceImpl.CommentAction(comment)
	if err != nil {
		log.Default()
	}
	fmt.Println(commentRes)
}

func TestCommentServiceImpl_DeleteCommentAction(t *testing.T) {
	redis.InitRedis()
	err := commentServiceImpl.DeleteCommentAction(1)
	if err != nil {
		log.Default()
	}
}

func TestCommentServiceImpl_GetCommentList(t *testing.T) {
	redis.InitRedis()
	commentList, err := commentServiceImpl.GetCommentList(24, 1)
	if err != nil {
		log.Default()
	}
	fmt.Println(commentList)
}
