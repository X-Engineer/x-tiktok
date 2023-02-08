package service

import (
	"fmt"
	"log"
	"testing"
)

func TestLikeServiceImpl_FavoriteAction(t *testing.T) {
	err := likeServiceImp.FavoriteAction(157, 5, 1)
	if err != nil {
		return
	}
}

func TestGetVideoLikedCount(t *testing.T) {
	likeCnt, err := likeServiceImp.GetVideoLikedCount(20)
	if err != nil {
		log.Default()
	}
	fmt.Println(likeCnt)
}

func TestGetUserLikeCount(t *testing.T) {
	likeCnt, err := likeServiceImp.GetUserLikeCount(5)
	if err != nil {
		log.Default()
	}
	fmt.Println(likeCnt)
}

func TestLikeServiceImpl_IsLikedByUser(t *testing.T) {
	liked, err := likeServiceImp.IsLikedByUser(5, 23)
	if err != nil {
		log.Default()
	}
	log.Println(liked)
}

func TestLikeServiceImpl_GetUserLikedCnt(t *testing.T) {
	count, err := likeServiceImp.GetUserLikedCnt(5)
	if err != nil {
		log.Default()
	}
	log.Println(count)
}
