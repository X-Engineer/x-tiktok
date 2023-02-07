package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"sync"
	"x-tiktok/dao"
)

type UserServiceImpl struct {
	// 关注服务
	FollowService
	// 点赞服务

}

var (
	userServiceImp  *UserServiceImpl
	userServiceOnce sync.Once
)

func GetUserServiceInstance() *UserServiceImpl {
	userServiceOnce.Do(func() {
		userServiceImp = &UserServiceImpl{
			FollowService: &FollowServiceImp{},
		}
	})
	return userServiceImp
}

func (usi *UserServiceImpl) GetUserBasicInfoById(id int64) dao.UserBasicInfo {
	user, err := dao.GetUserBasicInfoById(id)
	if err != nil {
		log.Println("Err:", err.Error())
		log.Println("User Not Found")
		return user
	}
	log.Println("Query User Success")
	return user
}

func (usi *UserServiceImpl) GetUserBasicInfoByName(name string) dao.UserBasicInfo {
	user, err := dao.GetUserBasicInfoByName(name)
	if err != nil {
		log.Println("Err:", err.Error())
		log.Println("User Not Found")
		return user
	}
	log.Println("Query User Success")
	return user
}

func (usi *UserServiceImpl) InsertUser(user *dao.UserBasicInfo) bool {
	flag := dao.InsertUser(user)
	if flag == false {
		log.Println("Insert Fail!")
		return false
	}
	return true
}

func (usi *UserServiceImpl) GetUserLoginInfoById(id int64) (User, error) {
	user := User{
		Id:            5,
		Name:          "qcj",
		FollowCount:   1,
		FollowerCount: 99999,
		IsFollow:      false,
	}
	u, err := dao.GetUserBasicInfoById(id)
	fmt.Println(u)
	if err != nil {
		log.Println("Err:", err.Error())
		log.Println("User Not Found")
	}
	userService := GetUserServiceInstance()
	// 计算关注数
	followCnt, _ := userService.GetFollowingCnt(id)
	// 计算粉丝数
	followerCnt, _ := userService.GetFollowerCnt(id)
	// 计算作品数

	user = User{
		Id:            u.Id,
		Name:          u.Name,
		FollowCount:   followCnt,
		FollowerCount: followerCnt,
		IsFollow:      false,
	}

	return user, nil
}

// 给密码加密

func EnCoder(password string) string {
	h := hmac.New(sha256.New, []byte(password))
	sha := hex.EncodeToString(h.Sum(nil))
	return sha
}
