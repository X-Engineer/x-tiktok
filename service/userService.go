package service

import (
	"x-tiktok/dao"
)


type UserService interface {

	GetUserBasicInfoById(id int64) dao.UserBasicInfo

	GetUserBasicInfoByName(name string) dao.UserBasicInfo

	GetUserLoginInfoById(id int64) (User, error)

	InsertUser(user *dao.UserBasicInfo) bool


}

type User struct {
	Id             int64  `json:"id,omitempty"`
	Name           string `json:"name,omitempty"`
	FollowCount    int64  `json:"follow_count,omitempty"`
	FollowerCount  int64  `json:"follower_count,omitempty"`
	IsFollow       bool   `json:"is_follow,omitempty"`
}
