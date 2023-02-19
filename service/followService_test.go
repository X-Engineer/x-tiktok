package service

import (
	"fmt"
	"log"
	"testing"
)

func TestFollowServiceImp_GetFollowings(t *testing.T) {
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

	userIdCnt, err := followServiceImp.GetFollowingCnt(11)
	if err != nil {
		log.Default()
	}
	fmt.Println(userIdCnt)
}

func TestFollowServiceImp_GetFollowerCnt(t *testing.T) {

	userIdCnt, err := followServiceImp.GetFollowerCnt(11)
	if err != nil {
		log.Default()
	}
	fmt.Println(userIdCnt)
}

func TestFollowServiceImp_CheckIsFollowing(t *testing.T) {
	result, err := followServiceImp.CheckIsFollowing(11, 2)
	if err != nil {
		log.Default()
	}
	fmt.Println(result)
}
