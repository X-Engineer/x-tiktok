package service

import (
	"fmt"
	"log"
	"sync"
	"x-tiktok/dao"
	"x-tiktok/middleware/rabbitmq"
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
	// 获取点赞和取消点赞的消息队列
	likeAddMQ := rabbitmq.SimpleLikeAddMQ
	likeDelMQ := rabbitmq.SimpleLikeDelMQ
	if islike == -1 {
		//用户没有点赞过该视频
		//插入一条新记录
		//var likeinfo dao.Like
		//likeinfo.UserId = userId
		//likeinfo.VideoId = videoId
		//likeinfo.Liked = int8(actionType)
		//likeinfo.CreatedAt = time.Now()
		//likeinfo.UpdatedAt = time.Now()
		//err = dao.InsertLikeInfo(likeinfo)
		//return nil
		// 消息队列
		err := likeAddMQ.PublishSimple(fmt.Sprintf("%d-%d-%s", userId, videoId, "insert"))
		return err
	}
	//该用户曾对此视频点过赞
	//err = dao.UpdateLikeInfo(userId, videoId, int8(actionType))
	if actionType == 1 {
		err = likeAddMQ.PublishSimple(fmt.Sprintf("%d-%d-%s", userId, videoId, "update"))
	} else {
		err = likeDelMQ.PublishSimple(fmt.Sprintf("%d-%d-%s", userId, videoId, "update"))
	}
	if err != nil {
		log.Print(err.Error() + "Favorite action failed!")
		return err
	} else {
		log.Print("Favorite action succeed!")
	}
	return nil
}

// GetLikesList 获取点赞信息
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

// GetUserLikeCount 获取用户点赞数量
func (*LikeServiceImpl) GetUserLikeCount(userId int64) (int64, error) {
	_, likeCnt, err := dao.GetLikeListByUserId(userId)
	if err != nil {
		log.Print("Get like count failed!")
		return -1, err
	}
	return likeCnt, nil
}

// GetVideoLikedCount 获取视频点赞数
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

// IsLikedByUser 当前用户是否点赞该视频
func (*LikeServiceImpl) IsLikedByUser(userId int64, videoId int64) (bool, error) {
	liked, err := dao.IsLikedByUser(userId, videoId)
	if err != nil {
		return false, err
	}
	return liked, nil
}

// GetUserLikedCnt 计算用户被点赞的视频获赞总数
func (*LikeServiceImpl) GetUserLikedCnt(userId int64) (int64, error) {
	//likedVideoIdList, err := dao.GetUserVideoLikedByOther(userId)
	//if err != nil {
	//	return -1, nil
	//}
	//var count int64 = 0
	//for _, videoId := range likedVideoIdList {
	//	tmp, _ := dao.VideoLikedCount(videoId)
	//	count += tmp
	//}
	// 通过一条 sql 语句计算该用户的视频获赞总数，减少数据库访问次数
	count, err := dao.GetUserVideoLikedTotalCount(userId)
	if err != nil {
		return -1, err
	}
	return count, nil
}
