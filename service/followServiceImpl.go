package service

import (
	"log"
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

// FollowAction 关注操作的业务
func (*FollowServiceImp) FollowAction(userId int64, targetId int64) (bool, error) {
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

// CancelFollowAction 取关操作的业务
func (*FollowServiceImp) CancelFollowAction(userId int64, targetId int64) (bool, error) {
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

// GetFollowings 获取正在关注的用户详情列表业务
func (*FollowServiceImp) GetFollowings(userId int64) ([]User, error) {
	followDao := dao.NewFollowDaoInstance()

	userFollowingsId, userFollowingsCnt, err := followDao.GetFollowingsInfo(userId)

	if nil != err {
		log.Println(err.Error())
	}

	userFollowings := make([]User, userFollowingsCnt)

	for i := 0; int64(i) < userFollowingsCnt; i++ {
		userFollowings[i].Id = userFollowingsId[i]

		var err1 error
		userFollowings[i].Name, err1 = followDao.GetUserName(userFollowingsId[i])
		if nil != err1 {
			log.Println(err1.Error())
			return nil, err1
		}

		var err2 error
		userFollowings[i].FollowCount, err2 = followDao.GetFollowingCnt(userFollowingsId[i])
		if nil != err2 {
			log.Println(err2.Error())
			return nil, err2
		}

		var err3 error
		userFollowings[i].FollowerCount, err3 = followDao.GetFollowerCnt(userFollowingsId[i])
		if nil != err3 {
			log.Println(err3.Error())
			return nil, err3
		}

		userFollowings[i].IsFollow = true
	}

	return userFollowings, nil
}

// GetFollowers 获取粉丝详情列表业务
func (*FollowServiceImp) GetFollowers(userId int64) ([]User, error) {
	followDao := dao.NewFollowDaoInstance()

	userFollowersId, userFollowersCnt, err := followDao.GetFollowersInfo(userId)

	if nil != err {
		log.Println(err.Error())
	}

	userFollowers := make([]User, userFollowersCnt)

	for i := 0; int64(i) < userFollowersCnt; i++ {
		userFollowers[i].Id = userFollowersId[i]

		var err1 error
		userFollowers[i].Name, err1 = followDao.GetUserName(userFollowersId[i])
		if nil != err1 {
			log.Println(err1.Error())
			return nil, err1
		}

		var err2 error
		userFollowers[i].FollowCount, err2 = followDao.GetFollowingCnt(userFollowersId[i])
		if nil != err2 {
			log.Println(err2.Error())
			return nil, err2
		}

		var err3 error
		userFollowers[i].FollowerCount, err3 = followDao.GetFollowerCnt(userFollowersId[i])
		if nil != err3 {
			log.Println(err3.Error())
			return nil, err3
		}

		isFollowResult, err4 := followDao.FindEverFollowing(userId, userFollowersId[i])
		if nil != err4 {
			log.Println(err4.Error())
			return nil, err4
		}

		if nil != isFollowResult && isFollowResult.Followed == 1 {
			userFollowers[i].IsFollow = true
		} else {
			userFollowers[i].IsFollow = false
		}

	}

	return userFollowers, nil

}

func (*FollowServiceImp) GetFriends(userId int64) ([]User, error) {
	followDao := dao.NewFollowDaoInstance()

	userFollowersId, userFollowersCnt, err := followDao.GetFollowersInfo(userId)

	if nil != err {
		log.Println(err.Error())
	}

	userFollowers := make([]User, userFollowersCnt)

	for i := 0; int64(i) < userFollowersCnt; i++ {
		userFollowers[i].Id = userFollowersId[i]

		var err1 error
		userFollowers[i].Name, err1 = followDao.GetUserName(userFollowersId[i])
		if nil != err1 {
			log.Println(err1.Error())
			return nil, err1
		}

		var err2 error
		userFollowers[i].FollowCount, err2 = followDao.GetFollowingCnt(userFollowersId[i])
		if nil != err2 {
			log.Println(err2.Error())
			return nil, err2
		}

		var err3 error
		userFollowers[i].FollowerCount, err3 = followDao.GetFollowerCnt(userFollowersId[i])
		if nil != err3 {
			log.Println(err3.Error())
			return nil, err3
		}

		isFollowResult, err4 := followDao.FindEverFollowing(userId, userFollowersId[i])
		if nil != err4 {
			log.Println(err4.Error())
			return nil, err4
		}

		if nil != isFollowResult && isFollowResult.Followed == 1 {
			userFollowers[i].IsFollow = true
		} else {
			userFollowers[i].IsFollow = false
		}

	}

	return userFollowers, nil
}
