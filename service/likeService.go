package service

type LikeService interface {
	//点赞操作
	FavoriteAction(userId int64, videoId int64, actionType int8) error

	//获取点赞列表
	GetLikesList(userId int64) ([]int64, error)
	//获取视频列表
	GetVideo(videoId int64) ([]Video, error)

	//当前用户是否点赞该视频
	IsLikedByUser(userId int64, videoId int64) (bool, error)
}
