package service

type LikeService interface {
	// FavoriteAction 点赞操作
	FavoriteAction(userId int64, videoId int64, actionType int32) error

	//获取当前用户点赞列表
	GetLikesList(userId int64) ([]int64, error)
	////获取视频列表
	//GetVideo(videoId []int64, likeCnt int64) ([]Video, error)

	// IsLikedByUser 当前用户是否点赞该视频
	IsLikedByUser(userId int64, videoId int64) (bool, error)
	// GetUserLikeCount 获取用户点赞数量
	GetUserLikeCount(userId int64) (int64, error)
	// GetVideoLikeCount 获取视频点赞数
	GetVideoLikeCount(videoId int64) (int64, error)
}
