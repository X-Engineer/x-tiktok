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
	err := UpdateLikeInfo(5, 14, 1)
	if err != nil {
		log.Println(err)
	}
}

func TestGetLikeListByUserId(t *testing.T) {
	list, cnt, err := GetLikeListByUserId(5)
	if err != nil {
		log.Print(err.Error())
	}
	log.Println(cnt)
	for _, v := range list {
		fmt.Printf("%d\n", v)
	}
}

func TestIsVideoLikedByUser(t *testing.T) {
	islike, err := IsVideoLikedByUser(5, 200)
	if err != nil {
		log.Print(err.Error())
	}
	fmt.Printf("islike：%d\n", islike)
}

func TestVideoLikedCount(t *testing.T) {
	likeCnt, err := VideoLikedCount(20)
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

func TestIsLikedByUser(t *testing.T) {
	flag, err := IsLikedByUser(5, 23)
	if err != nil {
		log.Default()
	}
	log.Println(flag)
}

func TestGetUserVideoLikedByOther(t *testing.T) {
	likedList, err := GetUserVideoLikedByOther(5)
	if err != nil {
		log.Default()
	}
	log.Println(likedList)
}
