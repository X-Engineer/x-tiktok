package dao

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func TestInsertLikeInfo(t *testing.T) {
	err := InsertLikeInfo(Like{1001, 151, 1, 1, time.Now(), time.Now()})
	if err != nil {
		return
	}
}

func TestUpdateLikeInfo(t *testing.T) {
	err := UpdateLikeInfo(152, 5, 1)
	if err != nil {
		return
	}
}

func TestGetLikeListByUserId(t *testing.T) {
	list, _, err := GetLikeListByUserId(152)
	if err != nil {
		log.Print(err.Error())
	}
	for _, v := range list {
		fmt.Printf("%d\n", v)
	}
}

func TestIsVideoLikedByUser(t *testing.T) {
	islike, err := IsVideoLikedByUser(152, 5)
	if err != nil {
		log.Print(err.Error())
	}
	fmt.Printf("islike：%d\n", islike)
}

func TestVideoLikedCount(t *testing.T) {
	likeCnt, err := VideoLikedCount(9)
	if err != nil {
		log.Print(err.Error())
	}
	fmt.Printf("Like Count：%d\n", likeCnt)
}

func TestGetLikeCountByUser(t *testing.T) {
	likeCnt, err := GetLikeCountByUser(152)
	if err != nil {
		log.Print(err.Error())
	}
	fmt.Printf("Like Count：%d\n", likeCnt)
}
