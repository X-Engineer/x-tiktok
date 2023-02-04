package dao

import (
	"errors"
	"log"
	"time"
)

type Like struct {
	Id        int64
	UserId    int64
	VideoId   int64
	Liked     int8
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Like) TableName() string {
	return "like"
}

// 增删改查
// 获取当前用户点赞视频id列表
func GetLikeListByUserId(userId int64) ([]int64, int64, error) {
	var LikedList []int64
	result := Db.Model(&Like{}).Where("user_id=? and liked=?", userId, 1).Pluck("video_id", &LikedList)
	likeCnt := result.RowsAffected
	if result.Error != nil {
		log.Println("LikedVideoIdList:", result.Error.Error())
		return nil, -1, result.Error
	}
	return LikedList, likeCnt, nil
}

// 统计视频点赞数量
func VideoLikedCount(videoId int64) (int64, error) {
	var count int64
	//数据库中查询点赞数量
	err := Db.Model(Like{}).Where(map[string]interface{}{"video_id": videoId, "liked": 1}).Count(&count).Error
	if err != nil {
		log.Println("LikeDao-Count: return count failed") //函数返回提示错误信息
		return -1, errors.New("find likes count failed")
	}
	log.Println("LikeDao-Count: return count success") //函数执行成功，返回正确信息
	return count, nil
}

// 更新点赞数据
func UpdateLikeInfo(userId int64, videoId int64, liked int8) error {
	err := Db.Model(Like{}).Where(map[string]interface{}{"user_id": userId, "video_id": videoId}).Update("liked", liked).Error

	if err != nil {
		log.Println(err.Error) //函数返回提示错误信息
		return errors.New("update like failed")
	}
	log.Println("LikeDao-UpdateLikeInfo: return success") //函数执行成功，返回正确信息
	return nil
}

// 插入点赞数据
func InsertLikeInfo(like Like) error {
	err := Db.Model(Like{}).Create(&like).Error
	if err != nil {
		log.Println(err.Error())
		return errors.New("insert likes failed")
	}
	return nil
}

// 获取视频点赞信息（当前用户是否点赞）
func IsVideoLikedByUser(userId int64, videoId int64) (int8, error) {
	var isLiked int8
	result := Db.Model(Like{}).Where("user_id= ?", userId).Where("video_id= ?", videoId).Pluck("liked", &isLiked)
	c := result.RowsAffected
	if result.Error != nil {

		//如果查询数据库失败，返回获取likeInfo信息失败
		log.Println(result.Error.Error())
		return isLiked, errors.New("get likeInfo failed")
		//}
	} else if c == 0 {
		return -1, errors.New("record not found")
	}
	return isLiked, nil
}

// 获取用户点赞数
func GetLikeCountByUser(userId int64) (int64, error) {
	var count int64
	//数据库中查询点赞数量
	err := Db.Model(Like{}).Where(map[string]interface{}{"user_id": userId, "liked": 1}).Count(&count).Error
	if err != nil {
		log.Println("LikeDao-Count: return count failed") //函数返回提示错误信息
		return -1, errors.New("find likes count failed")
	}
	log.Println("LikeDao-Count: return count success") //函数执行成功，返回正确信息
	return count, nil
}
