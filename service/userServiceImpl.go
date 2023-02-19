package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"sync"
	"x-tiktok/config"
	"x-tiktok/dao"
)

type UserServiceImpl struct {
	// 关注服务
	FollowService
	// 点赞服务
	LikeService
	// 视频服务
	VideoService
}

var (
	userServiceImp  *UserServiceImpl
	userServiceOnce sync.Once
)

func GetUserServiceInstance() *UserServiceImpl {
	userServiceOnce.Do(func() {
		userServiceImp = &UserServiceImpl{
			FollowService: &FollowServiceImp{},
			LikeService:   &LikeServiceImpl{},
			VideoService:  &VideoServiceImpl{},
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

// GetUserLoginInfoById 未登录情况返回用户信息
func (usi *UserServiceImpl) GetUserLoginInfoById(id int64) (User, error) {
	user := User{
		Id:              5,
		Name:            "qcj",
		FollowCount:     1,
		FollowerCount:   99999,
		IsFollow:        false,
		Avatar:          config.CUSTOM_DOMAIN + config.OSS_USER_AVATAR_DIR,
		BackgroundImage: config.BG_IMAGE,
		Signature:       config.SIGNATURE,
		TotalFavorited:  10,
		FavoriteCount:   10,
		WorkCount:       8,
	}
	u, err := dao.GetUserBasicInfoById(id)
	fmt.Println(u)
	if err != nil {
		log.Println("Err:", err.Error())
		log.Println("User Not Found")
	}
	user.Id = u.Id
	user.Name = u.Name
	userService := GetUserServiceInstance()
	var wg sync.WaitGroup
	wg.Add(5)
	go func(id int64) {
		// 计算关注数
		followCnt, err := userService.GetFollowingCnt(id)
		if err != nil {
			return
		}
		user.FollowCount = followCnt
		wg.Done()
	}(id)

	go func(id int64) {
		// 计算粉丝数
		followerCnt, _ := userService.GetFollowerCnt(id)
		if err != nil {
			return
		}
		user.FollowerCount = followerCnt
		wg.Done()
	}(id)

	go func(id int64) {
		// 计算作品数
		workCount, err := userService.GetVideoCnt(id)
		if err != nil {
			return
		}
		user.WorkCount = workCount
		wg.Done()
	}(id)

	go func(id int64) {
		// 计算被点赞数, 找出用户被点赞的视频，循环求和:在likeservide实现
		totalFavorited, err := userService.GetUserLikedCnt(id)
		if err != nil {
			return
		}
		user.TotalFavorited = totalFavorited
		wg.Done()
	}(id)

	go func(id int64) {
		// 计算喜欢数量
		favoriteCount, err := userService.GetUserLikeCount(id)
		if err != nil {
			return
		}
		user.FavoriteCount = favoriteCount
		wg.Done()
	}(id)
	wg.Wait()
	return user, nil
}

// GetUserLoginInfoByIdWithCurId 登录情况下返回用户信息, 第一个id是视频作者的id，第二个id是我们用户的id
func (usi *UserServiceImpl) GetUserLoginInfoByIdWithCurId(id int64, curId int64) (User, error) {
	user := User{
		Id:              5,
		Name:            "qcj",
		FollowCount:     1,
		FollowerCount:   99999,
		IsFollow:        false,
		Avatar:          config.CUSTOM_DOMAIN + config.OSS_USER_AVATAR_DIR,
		BackgroundImage: config.BG_IMAGE,
		Signature:       config.SIGNATURE,
		TotalFavorited:  10,
		FavoriteCount:   10,
		WorkCount:       8,
	}
	u, err := dao.GetUserBasicInfoById(id)
	fmt.Println(u)
	if err != nil {
		log.Println("Err:", err.Error())
		log.Println("User Not Found")
	}
	user.Id = u.Id
	user.Name = u.Name
	userService := GetUserServiceInstance()

	var wg sync.WaitGroup
	wg.Add(6)
	go func(id int64) {
		// 计算关注数
		followCnt, err := userService.GetFollowingCnt(id)
		if err != nil {
			return
		}
		user.FollowCount = followCnt
		wg.Done()
	}(id)

	go func(id int64) {
		// 计算粉丝数
		followerCnt, _ := userService.GetFollowerCnt(id)
		if err != nil {
			return
		}
		user.FollowerCount = followerCnt
		wg.Done()
	}(id)

	go func(id int64, curId int64) {
		// 计算是否关注, 这个地方又有点奇怪？只有在当前登录的情况下关注作者，后面该作者的视频才会显示已关注；退出重新登录就没了！
		isFollow, err := userService.CheckIsFollowing(curId, id)
		if err != nil {
			return
		}
		user.IsFollow = isFollow
		wg.Done()

	}(id, curId)

	go func(id int64) {
		// 计算作品数
		workCount, err := userService.GetVideoCnt(id)
		if err != nil {
			return
		}
		user.WorkCount = workCount
		wg.Done()
	}(id)

	go func(id int64) {
		// 计算被点赞数, 找出用户被点赞的视频，循环求和:在likeservide实现
		totalFavorited, err := userService.GetUserLikedCnt(id)
		if err != nil {
			return
		}
		user.TotalFavorited = totalFavorited
		wg.Done()
	}(id)

	go func(id int64) {
		// 计算喜欢数量
		favoriteCount, err := userService.GetUserLikeCount(id)
		if err != nil {
			return
		}
		user.FavoriteCount = favoriteCount
		wg.Done()
	}(id)
	wg.Wait()
	return user, nil
}

// 给密码加密

func EnCoder(password string) string {
	h := hmac.New(sha256.New, []byte(password))
	sha := hex.EncodeToString(h.Sum(nil))
	return sha
}
