package dao

import (
	"fmt"
	"log"
	"testing"
)

func TestFollowDao_InsertFollowRelation(t *testing.T) {
	//followDao.InsertFollowRelation(2, 3)
}

func TestFollowDao_FindRelation(t *testing.T) {
	follow, err := followDao.FindEverFollowing(2, 3)
	if err == nil {
		log.Default()
	}
	fmt.Print(follow)
}

func TestFollowDao_UpdateFollowRelation(t *testing.T) {
	// followDao.UpdateFollowRelation(2, 3, 1)

}

func TestFollowDao_GetFollowingsInfo(t *testing.T) {
	followingsID, followingsCnt, err := followDao.GetFollowingsInfo(1)

	if err != nil {
		log.Default()
	}

	fmt.Println(followingsID)
	fmt.Println(followingsCnt)

}

func TestFollowDao_GetUserName(t *testing.T) {
	name, err := followDao.GetUserName(2)
	if err != nil {
		log.Default()
	}
	fmt.Println(name)
}

func TestFollowDao_GetFriendsInfo(t *testing.T) {
	friendId, friendCnt, _ := followDao.GetFriendsInfo(6)

	fmt.Println(friendId)
	fmt.Println(friendCnt)

}
