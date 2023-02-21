package dao

import (
	"log"
	"testing"
	"time"
	"x-tiktok/middleware/redis"
)

func TestInsertLikeInfo(t *testing.T) {
	err := InsertLikeInfo(Like{76, 18, 17, 1, time.Now(), time.Now()})
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
	list, err := GetLikeListByUserId(5)
	if err != nil {
		log.Print(err.Error())
	}
	//log.Println(cnt)
	for _, v := range list {
		log.Printf("%d\n", v)
	}
}

func TestIsVideoLikedByUser(t *testing.T) {
	islike, err := IsVideoLikedByUser(5, 200)
	if err != nil {
		log.Print(err.Error())
	}
	log.Printf("islike：%d\n", islike)
}

func TestVideoLikedCount(t *testing.T) {
	redis.InitRedis()
	likeCnt, err := VideoLikedCount(20)
	if err != nil {
		log.Print(err.Error())
	}
	log.Printf("Like Count：%d\n", likeCnt)
}

func TestGetLikeCountByUser(t *testing.T) {
	redis.InitRedis()
	likeCnt, err := GetLikeCountByUser(5)
	if err != nil {
		log.Print(err.Error())
	}
	log.Printf("Like Count：%d\n", likeCnt)
}

func TestIsLikedByUser(t *testing.T) {
	flag, err := IsLikedByUser(5, 23)
	if err != nil {
		log.Default()
	}
	log.Println(flag)
}

//func TestGetUserVideoLikedByOther(t *testing.T) {
//	likedList, err := GetUserVideoLikedByOther(5)
//	if err != nil {
//		log.Default()
//	}
//	log.Println(likedList)
//}
