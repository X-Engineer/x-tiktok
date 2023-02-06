package service

import (
	"log"
	"sync"
	"time"
	"x-tiktok/config"
	"x-tiktok/dao"
)

type LikeServiceImpl struct {
	LikeService
}

var (
	likeServiceImp    *LikeServiceImpl
	likeServiceInstan sync.Once
)

func NewLikeServImpInstance() *LikeServiceImpl {
	likeServiceInstan.Do(
		func() {
			likeServiceImp = &LikeServiceImpl{}
		})
	return likeServiceImp
}

func (*LikeServiceImpl) FavoriteAction(userId int64, videoId int64, actionType int32) error {
	islike, err := dao.IsVideoLikedByUser(userId, videoId)
	log.Print("islike:", islike)
	log.Println("actionType:", actionType)
	if err != nil {
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
			return err
		} else {
			//查询失败
			log.Print(err.Error())
			return err
		}
	}
	//该用户曾对此视频点过赞
	err = dao.UpdateLikeInfo(userId, videoId, int8(actionType))
	if err != nil {
		log.Print(err.Error() + "Favorite action failed!")
		return err
	} else {
		log.Print("Favorite action succed!")
	}
	return nil
}

// 获取点赞信息
func (*LikeServiceImpl) GetLikesList(userId int64) ([]int64, error) {
	//likelist, likeCnt, err := dao.GetLikeListByUserId(userId)
	likelist, _, err := dao.GetLikeListByUserId(userId)
	if err != nil {
		log.Print("Get like list failed!")
		return nil, err
	}
	//for i := 0; i < int(likeCnt); i++ {
	//	//获取Video信息列表
	//	//VideoService(.)
	//}
	return likelist, nil
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
func (*LikeServiceImpl) GetVideoLikeCount(videoId int64) (int64, error) {
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
		log.Print("Get like info failed!")
		return false, err
	}
	if liked == config.LIKE {
		return true, nil
	}
	return false, nil
}
