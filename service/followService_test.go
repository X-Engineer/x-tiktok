package service

import (
	"fmt"
	"log"
	"testing"
	"x-tiktok/middleware/redis"
)

func TestFollowServiceImp_GetFollowings(t *testing.T) {

	redis.InitRedis()
	followings, err := followServiceImp.GetFollowings(1)

	if err != nil {
		log.Default()
	}
	fmt.Println(followings)
}

func TestFollowServiceImp_GetFollowers(t *testing.T) {
	followers, err := followServiceImp.GetFollowers(2)

	if err != nil {
		log.Default()
	}
	fmt.Println(followers)
}

func TestFollowServiceImp_GetFollowingCnt(t *testing.T) {

	redis.InitRedis()
	userIdCnt, err := followServiceImp.GetFollowingCnt(7)
	if err != nil {
		log.Default()
	}
	fmt.Println(userIdCnt)
}

func TestFollowServiceImp_GetFollowerCnt(t *testing.T) {

	redis.InitRedis()
	userIdCnt, err := followServiceImp.GetFollowerCnt(1)
	if err != nil {
		log.Default()
	}
	fmt.Println(userIdCnt)
}

func TestFollowServiceImp_CheckIsFollowing(t *testing.T) {
	redis.InitRedis()
	result, err := followServiceImp.CheckIsFollowing(1, 5)
	if err != nil {
		log.Default()
	}
	fmt.Println(result)
}

func TestFollowServiceImp_FollowAction(t *testing.T) {
	redis.InitRedis()
	result, err := followServiceImp.FollowAction(1, 16)
	if err != nil {
		log.Default()
	}
	fmt.Println(result)
}
