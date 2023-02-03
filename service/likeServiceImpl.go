package service

import (
	"log"
	"x-tiktok/config"
	"x-tiktok/dao"
)

//type User struct {
//	id int64   // 用户id
//	name string   // 用户名称
//	follow_count int64   // 关注总数
//	follower_count int64   // 粉丝总数
//	is_follow bool   // true-已关注，false-未关注
//}

type Video struct {
	id             int64  // 视频唯一标识
	author         User   // 视频作者信息
	play_url       string // 视频播放地址
	cover_url      string // 视频封面地址
	favorite_count int64  // 视频的点赞总数
	comment_count  int64  // 视频的评论总数
	is_favorite    bool   // true-已点赞，false-未点赞
	title          string // 视频标题
}

type LikeServiceImpl struct {
	LikeService
	FollowService
}

func (*LikeServiceImpl) FavoriteAction(userId int64, videoId int64, actionType int8) error {

	err := dao.UpdateLikeInfo(userId, videoId, actionType)
	if err != nil {
		log.Print("Favorite action failed!")
		return err
	}
	return nil
}

func GetVideo(videoId []int64,likeCnt int64) ([]Video, error) {
	var video []Video
	for i:=0;int64(i)<likeCnt;i++ {
		video[i].id=videoId[i]
		video[i].author=
	}

}

// 获取点赞信息
func (*LikeServiceImpl) GetLikesList(userId int64) ([]int64,int64, error) {
	likelist,likeCnt, err := dao.GetLikeListByUserId(userId)
	if err != nil {
		log.Print("Get like list failed!")
		return nil,-1, err
	}
	return likelist,likeCnt, nil
}

// 当前用户是否点赞该视频
func (*LikeServiceImpl) IsLikedByUser(userId int64, videoId int64) (bool, error) {
	liked, err := dao.GetVideoLikedByUser(userId, videoId)
	if err != nil {
		log.Print("Get like info failed!")
		return false, err
	}
	if liked == config.LIKE {
		return true, nil
	}
	return false, nil
}
