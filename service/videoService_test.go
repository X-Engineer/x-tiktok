package service

import (
	"fmt"
	"log"
	"testing"
	"time"
	"x-tiktok/middleware/redis"
)

func TestVideoServiceImpl_PublishList(t *testing.T) {
	redis.InitRedis()
	videoList, err := videoServiceImp.PublishList(1)
	if err != nil {
		log.Default()
	}
	fmt.Println(videoList)
}

func TestVideoServiceImpl_Feed(t *testing.T) {
	redis.InitRedis()
	videoList, nextTime, err := videoServiceImp.Feed(time.Now(), 1)
	if err != nil {
		log.Default()
	}
	fmt.Println(nextTime)
	fmt.Println(videoList)
}
