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
	likeCnt, err := likeServiceImp.GetUserLikeCount(1)
	if err != nil {
		log.Default()
	}
	fmt.Println(likeCnt)
}
