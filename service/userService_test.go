package service

import (
	"fmt"
	"log"
	"testing"
	"x-tiktok/dao"
)

func TestUserServiceImpl_GetUserBasicInfoById(t *testing.T) {
	userBasicInfo := userServiceImp.GetUserBasicInfoById(1)
	fmt.Println(userBasicInfo)
}

func TestUserServiceImpl_GetUserBasicInfoByName(t *testing.T) {
	userBasicInfo := userServiceImp.GetUserBasicInfoByName("qcj")
	fmt.Println(userBasicInfo)
}

func TestUserServiceImpl_GetUserLoginInfoById(t *testing.T) {
	userLoginInfo, err := userServiceImp.GetUserLoginInfoById(1)
	if err != nil {
		log.Default()
	}
	fmt.Println(userLoginInfo)
}

func TestUserServiceImpl_InsertUser(t *testing.T) {
	user := dao.UserBasicInfo{Name: "unit test in service", Password: "unit test in service"}
	flag := userServiceImp.InsertUser(&user)
	fmt.Println(flag)
}
