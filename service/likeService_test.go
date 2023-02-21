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
	err := likeServiceImp.FavoriteAction(18, 15, 2)
	if err != nil {
		return
	}
}

func TestGetVideoLikedCount(t *testing.T) {
	redis.InitRedis()
	likeCnt, err := likeServiceImp.GetVideoLikedCount(25)
	if err != nil {
		log.Default()
	}
	fmt.Println(likeCnt)
}

func TestGetUserLikeCount(t *testing.T) {
	redis.InitRedis()
	likeCnt, err := likeServiceImp.GetUserLikeCount(1)
	if err != nil {
		log.Default()
	}
	log.Println(likeCnt)
}

func TestLikeServiceImpl_IsLikedByUser(t *testing.T) {
	redis.InitRedis()
	liked, err := likeServiceImp.IsLikedByUser(12, 23)
	if err != nil {
		log.Default()
	}
	log.Println(liked)
}

func TestLikeServiceImpl_GetUserLikedCnt(t *testing.T) {
	count, err := likeServiceImp.GetUserLikedCnt(1)
	if err != nil {
		log.Default()
	}
	log.Println(count)
}

func TestRdsGetUserLikedCnt(t *testing.T) {
	redis.InitRedis()
	count, err := likeServiceImp.RdsGetUserLikedCnt(1)
	if err != nil {
		log.Default()
	}
	log.Println(count)
}
