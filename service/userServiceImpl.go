package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"x-tiktok/dao"
)

type UserServiceImpl struct {

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
	if flag == false{
		log.Println("Insert Fail!")
		return false
	}
	return true
}

func (usi *UserServiceImpl) GetUserLoginInfoById(id int64) (User,error) {
	user := User{
		Id:             5,
		Name:           "qcj",
		FollowCount:    1,
		FollowerCount:  99999,
		IsFollow:       false,
	}
	u, err := dao.GetUserBasicInfoById(id)
	if err != nil {
		log.Println("Err:", err.Error())
		log.Println("User Not Found")
	}
	// 计算关注数

	// 计算粉丝数

	// 计算作品数

	// 计算喜欢数
	user = User{
		Id:             u.Id,
		Name:           u.Name,
		FollowCount:    1,
		FollowerCount:  99999,
		IsFollow:       false,
	}

	return user, nil
}

// 给密码加密

func EnCoder(password string) string {
	h := hmac.New(sha256.New,[]byte(password))
	sha := hex.EncodeToString(h.Sum(nil))
	return sha
}

func GenerateToken(username string) string {
	return "1234567"
}

