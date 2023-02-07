package service

import (
	"fmt"
	"log"
	"testing"
	"x-tiktok/dao"
)

func TestCommentServiceImpl_GetCommentCnt(t *testing.T) {
	count, err := commentServiceImpl.GetCommentCnt(14)
	if err != nil {
		log.Default()
	}
	fmt.Println(count)
}

func TestCommentServiceImpl_CommentAction(t *testing.T) {
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
	err := commentServiceImpl.DeleteCommentAction(1)
	if err != nil {
		log.Default()
	}
}

func TestCommentServiceImpl_GetCommentList(t *testing.T) {
	commentList, err := commentServiceImpl.GetCommentList(14, 5)
	if err != nil {
		log.Default()
	}
	fmt.Println(commentList)
}
