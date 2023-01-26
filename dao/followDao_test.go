package dao

import (
	"fmt"
	"log"
	"testing"
)

func TestFollowDao_InsertFollowRelation(t *testing.T) {
	followDao.InsertFollowRelation(2, 3)
}

func TestFollowDao_FindRelation(t *testing.T) {
	follow, err := followDao.FindEverFollowing(2, 3)
	if err == nil {
		log.Default()
	}
	fmt.Print(follow)
}

func TestFollowDao_UpdateFollowRelation(t *testing.T) {
	followDao.UpdateFollowRelation(2, 3, 1)
}
