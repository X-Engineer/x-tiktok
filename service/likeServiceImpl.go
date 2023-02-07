package service

import (
	"log"
	"sync"
	"time"
	"x-tiktok/config"
	"x-tiktok/dao"
)

type LikeServiceImpl struct {
	VideoService
}

var (
	likeServiceImp      *LikeServiceImpl
	likeServiceInstance sync.Once
)

func NewLikeServImpInstance() *LikeServiceImpl {
	likeServiceInstance.Do(func() {
		likeServiceImp = &LikeServiceImpl{
			VideoService: &VideoServiceImpl{},
		}
	})
	return likeServiceImp
}

func (*LikeServiceImpl) FavoriteAction(userId int64, videoId int64, actionType int32) error {
	islike, err := dao.IsVideoLikedByUser(userId, videoId)
	log.Print("islike:", islike)
	log.Println("actionType:", actionType)
	if islike == -1 {
		//用户没有点赞过该视频
		//插入一条新记录
		var likeinfo dao.Like
		likeinfo.UserId = userId
		likeinfo.VideoId = videoId
		likeinfo.Liked = int8(actionType)
		likeinfo.CreatedAt = time.Now()
		likeinfo.UpdatedAt = time.Now()
		err = dao.InsertLikeInfo(likeinfo)
		return nil
	}
	//该用户曾对此视频点过赞
	err = dao.UpdateLikeInfo(userId, videoId, int8(actionType))
	if err != nil {
		log.Print(err.Error() + "Favorite action failed!")
		return err
	} else {
		log.Print("Favorite action succeed!")
	}
	return nil
}

// 获取点赞信息
func (*LikeServiceImpl) GetLikesList(userId int64) ([]Video, error) {
	likedVideoIdList, _, err := dao.GetLikeListByUserId(userId)
	if err != nil {
		log.Print("Get like list failed!")
		return nil, err
	}
	likeService := NewLikeServImpInstance()
	likedVideoInfoList, err := likeService.GetVideoListById(likedVideoIdList, userId)
	if err != nil {
		log.Println("Get videoList failed")
	}
	return likedVideoInfoList, nil
}

// 获取用户点赞数量
func (*LikeServiceImpl) GetUserLikeCount(userId int64) (int64, error) {
	_, likeCnt, err := dao.GetLikeListByUserId(userId)
	if err != nil {
		log.Print("Get like count failed!")
		return -1, err
	}
	return likeCnt, nil
}

// 获取视频点赞数
func (*LikeServiceImpl) GetVideoLikedCount(videoId int64) (int64, error) {
	//校验是否存在该视频
	//？？？
	likeCnt, err := dao.VideoLikedCount(videoId)
	if err != nil {
		log.Print("Get like count failed!")
		return -1, err
	}
	return likeCnt, nil
}

// 当前用户是否点赞该视频
func (*LikeServiceImpl) IsLikedByUser(userId int64, videoId int64) (bool, error) {
	//先校验是否存在该用户
	//？？？
	liked, err := dao.IsVideoLikedByUser(userId, videoId)
	if err != nil {
		return false, err
	}
	if liked == config.LIKE {
		return true, nil
	}
	return false, nil
}

//func GetLikeService() LikeServiceImpl {
//	var videoService VideoServiceImpl
//	var likeService LikeServiceImpl
//	videoService.LikeService = &likeService
//	likeService.VideoService = &videoService
//	return likeService
//}
