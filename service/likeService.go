package service

type LikeService interface {
	// FavoriteAction 点赞操作
	FavoriteAction(userId int64, videoId int64, actionType int32) error

	// GetLikesList 获取当前用户点赞列表
	GetLikesList(userId int64) ([]Video, error)
	////获取视频列表
	//GetVideo(videoId []int64, likeCnt int64) ([]Video, error)

	// IsLikedByUser 当前用户是否点赞该视频
	IsLikedByUser(userId int64, videoId int64) (bool, error)
	// GetUserLikeCount 获取用户点赞数量
	GetUserLikeCount(userId int64) (int64, error)
	// GetVideoLikedCount 获取视频点赞数
	GetVideoLikedCount(videoId int64) (int64, error)

	// GetUserLikedCnt 计算用户被点赞的视频获赞总数
	GetUserLikedCnt(userId int64) (int64, error)
}
