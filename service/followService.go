package service

type FollowService interface {

	// AddFollowRelation 当前用户关注目标用户
	AddFollowRelation(userId int64, targetId int64) (bool, error)
	// DeleteFollowRelation 当前用户取消对目标用户的关注
	DeleteFollowRelation(userId int64, targetId int64) (bool, error)
	// GetFollowing 获取当前用户的关注列表
	//GetFollowing(userId int64) ([]User, error)
	// GetFollowers 获取当前用户的粉丝列表
	//GetFollowers(userId int64) ([]User, error)
}
