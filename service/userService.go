package service

import (
	"x-tiktok/dao"
)

type UserService interface {

	// GetUserBasicInfoById 根据Id获取用户的用户名和密码
	GetUserBasicInfoById(id int64) dao.UserBasicInfo

	// GetUserBasicInfoByName 根据用户名获取用户的用户名和密码
	GetUserBasicInfoByName(name string) dao.UserBasicInfo

	// GetUserLoginInfoById 根据用户id获取用户的详细信息（未登录）
	GetUserLoginInfoById(id int64) (User, error)

	// GetUserLoginInfoByIdWithCurId 根据用户id获取用户的详细信息 (登录)
	GetUserLoginInfoByIdWithCurId(id int64, curId int64) (User, error)

	// InsertUser 添加一个用户
	InsertUser(user *dao.UserBasicInfo) bool
}

type User struct {
	Id              int64  `json:"id"`
	Name            string `json:"name"`
	FollowCount     int64  `json:"follow_count"`
	FollowerCount   int64  `json:"follower_count"`
	IsFollow        bool   `json:"is_follow"`
	Avatar          string `json:"avatar"`
	BackgroundImage string `json:"background_image"`
	Signature       string `json:"signature"`
	TotalFavorited  int64  `json:"total_favorited"`
	WorkCount       int64  `json:"work_count"`
	FavoriteCount   int64  `json:"favorite_count"`
}
