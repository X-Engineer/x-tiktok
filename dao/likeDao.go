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

// GetLikeListByUserId 获取当前用户点赞视频id列表
func GetLikeListByUserId(userId int64) ([]int64, int64, error) {
	var LikedList []int64
	result := Db.Model(&Like{}).Where("user_id=? and liked=?", userId, 1).Order("created_at desc").Pluck("video_id", &LikedList)
	likeCnt := result.RowsAffected
	if result.Error != nil {
		log.Println("LikedVideoIdList:", result.Error.Error())
		return nil, -1, result.Error
	}
	return LikedList, likeCnt, nil
}

// VideoLikedCount 统计视频点赞数量
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

// UsersOfLikeVideo 点赞该视频的用户 id 列表和数量
func UsersOfLikeVideo(videoId int64) ([]int64, int64, error) {
	var userIdList []int64
	result := Db.Model(&Like{}).Where("video_id=? and liked=?", videoId, 1).Pluck("user_id", &userIdList)
	likeCnt := result.RowsAffected
	if likeCnt == 0 {
		return nil, 0, result.Error
	}
	if result.Error != nil {
		log.Println("UsersOfLikeVideo:", result.Error.Error())
		return nil, 0, result.Error
	}
	return userIdList, likeCnt, nil
}

// UpdateLikeInfo 更新点赞数据
func UpdateLikeInfo(userId int64, videoId int64, liked int8) error {
	// Update即使更新不存在的记录也不会报错
	result := Db.Model(Like{}).Where(map[string]interface{}{"user_id": userId, "video_id": videoId}).Update("liked", liked)
	if result.RowsAffected == 0 {
		return errors.New("update like failed, record not exists")
	}
	log.Println("LikeDao-UpdateLikeInfo: return success") //函数执行成功，返回正确信息
	return nil
}

// InsertLikeInfo 插入点赞数据
func InsertLikeInfo(like Like) error {
	err := Db.Model(Like{}).Create(&like).Error
	if err != nil {
		log.Println(err.Error())
		return errors.New("insert likes failed")
	}
	return nil
}

// IsVideoLikedByUser 获取视频点赞信息（当前用户是否点赞）
func IsVideoLikedByUser(userId int64, videoId int64) (int8, error) {
	var isLiked int8
	result := Db.Model(Like{}).Select("liked").Where("user_id= ? and video_id= ?", userId, videoId).First(&isLiked)
	c := result.RowsAffected
	if c == 0 {
		return -1, errors.New("current user haven not liked current video")
	}
	if result.Error != nil {
		//如果查询数据库失败，返回获取likeInfo信息失败
		log.Println(result.Error)
	}
	return isLiked, nil
}

// GetLikeCountByUser 获取用户点赞数
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

// IsLikedByUser 重写用户是否点赞视频
func IsLikedByUser(userId int64, videoId int64) (bool, error) {
	var like Like
	result := Db.Model(Like{}).
		Where("user_id = ? and video_id = ? and liked = ?", userId, videoId, 1).
		First(&like)
	if result.RowsAffected == 0 {
		//这里不能返回err, 没关注也是返回成功的；如果返回err的话，后续代码就阻塞调用这个函数的地方
		return false, nil
	}
	return true, nil
}

// GetUserVideoLikedByOther 计算用户被他人点赞的视频列表id
func GetUserVideoLikedByOther(userId int64) ([]int64, error) {
	var likedList []int64
	result := Db.Model(Like{}).
		Joins("join video on like.video_id = video.id and author_id = ? and liked = ?", userId, 1).
		Distinct("video_id").
		Pluck("video_id", &likedList)
	if result.Error != nil {
		return nil, result.Error
	}
	return likedList, nil
}

// GetUserVideoLikedTotalCount 获取用户发布视频的总获赞数量
func GetUserVideoLikedTotalCount(userId int64) (int64, error) {
	var totalLikedCount int64
	result := Db.Model(Like{}).Joins("join video on like.video_id = video.id and author_id = ? and liked = ?", userId, 1).Count(&totalLikedCount)
	if result.Error != nil {
		return 0, result.Error
	}
	return totalLikedCount, nil
}
