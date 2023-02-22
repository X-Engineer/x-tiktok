package service

import (
	"fmt"
	"log"
	"testing"
	"x-tiktok/middleware/rabbitmq"
	"x-tiktok/middleware/redis"
)

func TestLikeServiceImpl_FavoriteAction(t *testing.T) {
	redis.InitRedis()
	rabbitmq.InitRabbitMQ()
	rabbitmq.InitLikeRabbitMQ()
	err := likeServiceImp.FavoriteAction(6, 21, 2)
	if err != nil {
		return
	}
}

func TestGetVideoLikedCount(t *testing.T) {
	redis.InitRedis()
	likeCnt, err := likeServiceImp.GetVideoLikedCount(20)
	if err != nil {
		log.Default()
	}
	fmt.Println(likeCnt)
}

func TestGetUserLikeCount(t *testing.T) {
	redis.InitRedis()
	likeCnt, err := likeServiceImp.GetUserLikeCount(5)
	if err != nil {
		log.Default()
	}
	fmt.Println(likeCnt)
}

func TestLikeServiceImpl_IsLikedByUser(t *testing.T) {
	redis.InitRedis()
	liked, err := likeServiceImp.IsLikedByUser(2, 23)
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

func TestLikeServiceImpl_GetLikesList(t *testing.T) {
	redis.InitRedis()
	list, _ := likeServiceImp.GetLikesList(7)
	log.Println(list)
}
