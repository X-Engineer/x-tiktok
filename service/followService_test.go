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
