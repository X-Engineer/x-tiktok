package dao

import (
	"errors"
	"log"
)

type Like struct {
	Id        int64
	UserId    int64
	VideoId   int64
	Liked     int8
	CreatedAt string
	UpdatedAt string
}

func (Like) TableName() string {
	return "likes"
}

// 增删改查
// 获取当前用户点赞视频列表
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
	err := Db.Model(Like{}).Where("user_id=?,video_id=?", userId, videoId).Update("liked", liked)
	if err != nil {
		log.Println("LikeDao-UpdateLikeInfo: return update likes failed") //函数返回提示错误信息
		return errors.New("update likes failed")
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
func GetVideoLikedByUser(userId int64, videoId int64) (int8, error) {
	var isLiked int8
	err := Db.Model(Like{}).Where(map[string]interface{}{"user_id": userId, "video_id": videoId}).Pluck("liked", isLiked).Error
	if err != nil {
		if "record not found" == err.Error() {
			log.Println("can't find data")
			return -1, nil
		} else {
			//如果查询数据库失败，返回获取likeInfo信息失败
			log.Println(err.Error())
			return isLiked, errors.New("get likeInfo failed")
		}
	}
	return isLiked, nil
}

//
//func GetLikeInfo(userId int64) ([]int64, int64, error) {
//
//	var likeCnt int64
//	var videoId []int64
//
//	// following_id -> user_id
//	result := Db.Model(&Like{}).Where("user_id = ?", userId).Where("liked = ?", 1).Pluck("video_id", &videoId)
//	likeCnt = result.RowsAffected
//
//	if nil != result.Error {
//		log.Println(result.Error.Error())
//		return nil, 0, result.Error
//	}
//
//	return videoId, likeCnt, nil
//}
