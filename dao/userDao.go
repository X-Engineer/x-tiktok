package dao

import (
	"log"
	"time"
)

type UserBasicInfo struct {
	Id        int64
	Name      string
	Password  string
	CreatedAt time.Time `gorm:"CreatedAt" column:"CreatedAt"`
	UpdatedAt time.Time `gorm:"CreatedAt" column:"UpdatedAt"`
}

func (user UserBasicInfo) TableName() string {
	return "user"
}

func InsertUser(user *UserBasicInfo) bool {
	if err := Db.Create(&user).Error; err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}

func GetUserBasicInfoByName(name string) (UserBasicInfo, error) {
	user := UserBasicInfo{}
	if err := Db.Where("name = ?", name).First(&user).Error; err != nil {
		log.Println("获取用户信息读库失败", err.Error())
		return user, err
	}
	return user, nil
}

func GetUserBasicInfoById(id int64) (UserBasicInfo, error) {
	user := UserBasicInfo{}
	if err := Db.Where("id = ?", id).First(&user).Error; err != nil {
		log.Println(err.Error())
		return user, err
	}
	return user, nil
}
