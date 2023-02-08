package dao

import (
	"log"
	"testing"
)

func TestInsertComment(t *testing.T) {
	var comment Comment = Comment{
		UserId:  5,
		VideoId: 14,
		Content: "这条评论来自单元测试TestInsertComment",
	}
	commentRes, err := InsertComment(comment)
	if err != nil {
		log.Println(err)
	}
	log.Println("返回的评论是", commentRes)
}

func TestDeleteComment(t *testing.T) {
	err := DeleteComment(1)
	if err == nil {
		log.Println("delete comment success")
	}
}

func TestGetCommentList(t *testing.T) {
	commentList, err := GetCommentList(21)
	if err == nil {
		log.Println(commentList)
	}
}

func TestGetCommentCnt(t *testing.T) {
	count, err := GetCommentCnt(21)
	if err == nil {
		log.Println(count)
	}
}
