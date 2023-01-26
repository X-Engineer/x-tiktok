package service

import (
	"sync"
	"x-tiktok/dao"
)

// FollowServiceImp 该结构体继承FollowService接口。
type FollowServiceImp struct {
	FollowService
}

var (
	followServiceImp  *FollowServiceImp //controller层通过该实例变量调用service的所有业务方法。
	followServiceOnce sync.Once         //限定该service对象为单例，节约内存。
)

// NewFSIInstance 生成并返回FollowServiceImp结构体单例变量。
func NewFSIInstance() *FollowServiceImp {
	followServiceOnce.Do(
		func() {
			followServiceImp = &FollowServiceImp{
				//todo 这块暂时不考虑
				//UserService: &UserServiceImpl{
				//	// 存在我调userService中，userService又要调我。
				//	FollowService: &FollowServiceImp{},
				//},
			}
		})
	return followServiceImp
}

// AddFollowRelation 给定当前用户和目标对象id，添加他们之间的关注关系。
func (*FollowServiceImp) AddFollowRelation(userId int64, targetId int64) (bool, error) {
	followDao := dao.NewFollowDaoInstance()
	follow, err := followDao.FindEverFollowing(userId, targetId)
	// 寻找SQL 出错。
	if nil != err {
		return false, err
	}
	// 曾经关注过，只需要update一下followed即可。
	if nil != follow {
		_, err := followDao.UpdateFollowRelation(userId, targetId, 1)
		// update 出错。
		if nil != err {
			return false, err
		}
		// update 成功。
		return true, nil
	}
	// 曾经没有关注过，需要插入一条关注关系。
	_, err = followDao.InsertFollowRelation(userId, targetId)
	if nil != err {
		// insert 出错
		return false, err
	}
	// insert 成功。
	return true, nil
}

// DeleteFollowRelation 给定当前用户和目标用户id，删除其关注关系。
func (*FollowServiceImp) DeleteFollowRelation(userId int64, targetId int64) (bool, error) {
	followDao := dao.NewFollowDaoInstance()
	follow, err := followDao.FindEverFollowing(userId, targetId)
	// 寻找 SQL 出错。
	if nil != err {
		return false, err
	}
	// 曾经关注过，只需要update一下cancel即可。
	if nil != follow {
		_, err := followDao.UpdateFollowRelation(userId, targetId, 0)
		// update 出错。
		if nil != err {
			return false, err
		}
		// update 成功。
		return true, nil
	}
	// 没有关注关系
	return false, nil
}
