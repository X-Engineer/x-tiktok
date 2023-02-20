package dao

import (
	"errors"
	"log"
	"strconv"
	"time"
	"x-tiktok/config"
	"x-tiktok/middleware/redis"
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
func GetLikeListByUserId(userId int64) ([]int64, error) {
	var LikedList []int64
	strUserId := strconv.FormatInt(userId, 10)
	result := Db.Model(&Like{}).Where("user_id=? and liked=?", userId, 1).Order("created_at desc").Pluck("video_id", &LikedList)

	if result.Error != nil {
		log.Println("LikedVideoIdList:", result.Error.Error())
		redis.RdbLikeUserId.Del(redis.Ctx, strUserId)
		return nil, result.Error
	}
	//遍历videoIdList,添加进key的集合中，若失败，删除key，并返回错误信息，这么做的原因是防止脏读，
	//保证redis与mysql数据一致性
	for _, likeVideoId := range LikedList {
		log.Printf("%d", likeVideoId)
		log.Printf("更新RdbLikeUserId缓存ing……")
		if _, err2 := redis.RdbLikeUserId.SAdd(redis.Ctx, strUserId, likeVideoId).Result(); err2 != nil {
			log.Printf("方法:GetFavouriteList RedisLikeUserId add value失败")
			redis.RdbLikeUserId.Del(redis.Ctx, strUserId)
			return nil, err2
		}
		log.Printf("更新RdbLikeVideoId缓存ing……")
		strVideoId := strconv.FormatInt(likeVideoId, 10)
		if _, err2 := redis.RdbLikeVideoId.SAdd(redis.Ctx, strVideoId, userId).Result(); err2 != nil {
			log.Printf("方法:GetFavouriteList RdbLikeVideoId add value失败")
			redis.RdbLikeVideoId.Del(redis.Ctx, strVideoId)
			return nil, err2
		}
	}
	return LikedList, nil
}

// VideoLikedCount 统计视频点赞数量
func VideoLikedCount(videoId int64) (int64, error) {
	var count int64
	strVideoId := strconv.FormatInt(videoId, 10)
	//数据库中查询点赞数量
	err := Db.Model(Like{}).Where(map[string]interface{}{"video_id": videoId, "liked": 1}).Count(&count).Error
	if err != nil {
		log.Println("LikeDao-Count: return count failed") //函数返回提示错误信息
		return -1, errors.New("find likes count failed")
	}
	log.Println("LikeDao-Count: return count success") //函数执行成功，返回正确信息

	//维护Redis中videoId的点赞用户信息
	//获取点赞当前视频的用户列表
	userIdList, err1 := GetLikeUserList(videoId)
	if err1 != nil {
		log.Printf(err1.Error())
		return 0, err1
	}
	//依次将用户id添加到键为strVideoId的set中
	for _, likeuserId := range userIdList {
		redis.RdbLikeVideoId.SAdd(redis.Ctx, strVideoId, likeuserId)
	}
	return count, nil
}

// UpdateLikeInfo 更新点赞数据
func UpdateLikeInfo(userId int64, videoId int64, liked int8) error {
	strVideoId := strconv.FormatInt(videoId, 10)
	strUserId := strconv.FormatInt(userId, 10)
	// Update即使更新不存在的记录也不会报错
	result := Db.Model(Like{}).Where(map[string]interface{}{"user_id": userId, "video_id": videoId}).Update("liked", liked)
	if result.RowsAffected == 0 {
		return errors.New("update like failed, record not exists")
	}
	log.Println("LikeDao-UpdateLikeInfo: return success") //函数执行成功，返回正确信息
	if liked == config.LIKE {
		if _, err1 := redis.RdbLikeVideoId.SAdd(redis.Ctx, strVideoId, userId).Result(); err1 != nil {
			log.Printf("方法:FavouriteAction RedisLikeVideoId add value失败：%v", err1)
			return err1
		}
		if _, err2 := redis.RdbLikeUserId.SAdd(redis.Ctx, strUserId, videoId).Result(); err2 != nil {
			log.Printf("方法:FavouriteAction RedisLikeUserId add value失败：%v", err2)
			return err2
		}
	} else {
		if _, err1 := redis.RdbLikeVideoId.SRem(redis.Ctx, strVideoId, userId).Result(); err1 != nil {
			log.Printf("方法:FavouriteAction RedisLikeVideoId add value失败：%v", err1)
			return err1
		}
		if _, err2 := redis.RdbLikeUserId.SRem(redis.Ctx, strUserId, videoId).Result(); err2 != nil {
			log.Printf("方法:FavouriteAction RedisLikeUserId add value失败：%v", err2)
			return err2
		}
	}

	return nil
}

// InsertLikeInfo 插入点赞数据
func InsertLikeInfo(like Like) error {
	strVideoId := strconv.FormatInt(like.VideoId, 10)
	strUserId := strconv.FormatInt(like.UserId, 10)
	err := Db.Model(Like{}).Create(&like).Error
	if err != nil {
		log.Println(err.Error())
		return errors.New("insert likes failed")
	}
	if _, err1 := redis.RdbLikeVideoId.SAdd(redis.Ctx, strVideoId, like.UserId).Result(); err1 != nil {
		log.Printf("方法:FavouriteAction RedisLikeVideoId add value失败：%v", err1)
		return err1
	}
	if _, err2 := redis.RdbLikeUserId.SAdd(redis.Ctx, strUserId, like.VideoId).Result(); err2 != nil {
		log.Printf("方法:FavouriteAction RedisLikeUserId add value失败：%v", err2)
		return err2
	}
	return nil
}

// IsVideoLikedByUser 获取视频点赞信息（当前用户是否点赞）
func IsVideoLikedByUser(userId int64, videoId int64) (int8, error) {
	var isLiked int8
	//将int64 videoId转换为 string strVideoId
	strUserId := strconv.FormatInt(userId, 10)
	//将int64 videoId转换为 string strVideoId
	strVideoId := strconv.FormatInt(videoId, 10)
	result := Db.Model(Like{}).Select("liked").Where("user_id= ? and video_id= ?", userId, videoId).First(&isLiked)
	c := result.RowsAffected
	if c == 0 {
		return -1, nil
	}
	if result.Error != nil {
		//如果查询数据库失败，返回获取likeInfo信息失败
		log.Println(result.Error)
	}
	if isLiked == config.LIKE {
		redis.RdbLikeUserId.SAdd(redis.Ctx, strUserId, videoId)
		redis.RdbLikeVideoId.SAdd(redis.Ctx, strVideoId, userId)
	}
	return isLiked, nil
}

// GetLikeCountByUser 获取用户点赞数
func GetLikeCountByUser(userId int64) (int64, error) {
	//var count int64
	//数据库中查询点赞数量
	likelist, err := GetLikeListByUserId(userId)
	count := len(likelist)
	log.Printf("count:%d", count)
	if err != nil {
		log.Println("LikeDao-Count: return count failed") //函数返回提示错误信息
		return -1, errors.New("find likes count failed")
	}
	strUserId := strconv.FormatInt(userId, 10)
	redis.RdbLikeUserCnt.Set(redis.Ctx, strUserId, count, time.Hour)
	log.Println("LikeDao-Count: return count success") //函数执行成功，返回正确信息

	return int64(count), nil
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

func GetLikeUserList(videoId int64) ([]int64, error) {
	var LikeUserList []int64
	result := Db.Model(&Like{}).Where("video_id=? and liked=?", videoId, 1).Order("created_at desc").Pluck("user_id", &LikeUserList)

	if result.Error != nil {
		log.Println("LikeUserIdList:", result.Error.Error())
		return nil, result.Error
	}
	strVideoId := strconv.FormatInt(videoId, 10)
	//遍历userIdList,添加进key的集合中，若失败，删除key，并返回错误信息，这么做的原因是防止脏读，
	//保证redis与mysql数据一致性
	for _, likeUserId := range LikeUserList {
		if _, err1 := redis.RdbLikeVideoId.SAdd(redis.Ctx, strVideoId, likeUserId).Result(); err1 != nil {
			log.Printf("方法:FavouriteAction RedisLikeVideoId add value失败")
			redis.RdbLikeVideoId.Del(redis.Ctx, strVideoId)
			return nil, err1
		}
	}
	return LikeUserList, nil
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
