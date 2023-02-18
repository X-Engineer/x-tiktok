package service

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"x-tiktok/config"
	"x-tiktok/dao"
	"x-tiktok/middleware/rabbitmq"
	"x-tiktok/middleware/redis"
)

// FollowServiceImp 该结构体继承FollowService接口。
type FollowServiceImp struct {
	//MessageService
	FollowService
}

var (
	followServiceImp  *FollowServiceImp //controller层通过该实例变量调用service的所有业务方法。
	followServiceOnce sync.Once         //限定该service对象为单例，节约内存。
)

// RedisFollowPrefix 前缀
var RedisFollowPrefix = "follow:"

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
func (followService *FollowServiceImp) FollowAction(userId int64, targetId int64) (bool, error) {

	AddToRDBWhenFollow(int(userId), int(targetId))

	followDao := dao.NewFollowDaoInstance()
	follow, err := followDao.FindEverFollowing(userId, targetId)
	// 获取关注的消息队列
	followAddMQ := rabbitmq.SimpleFollowAddMQ
	// 寻找SQL 出错。
	if nil != err {
		return false, err
	}
	// 曾经关注过，只需要update一下followed即可。
	if nil != follow {
		//发送消息队列
		err := followAddMQ.PublishSimpleFollow(fmt.Sprintf("%d-%d-%s", userId, targetId, "update"))
		if err != nil {
			return false, err
		}
		return true, nil
		//_, err := followDao.UpdateFollowRelation(userId, targetId, 1)
		//// update 出错。
		//if nil != err {
		//	return false, err
		//}
		//// update 成功。
		//return true, nil
	}
	//发送消息队列
	err = followAddMQ.PublishSimpleFollow(fmt.Sprintf("%d-%d-%s", userId, targetId, "insert"))
	if err != nil {
		return false, err
	}
	return true, nil
	// 曾经没有关注过，需要插入一条关注关系。
	//_, err = followDao.InsertFollowRelation(userId, targetId)
	//if nil != err {
	//	// insert 出错
	//	return false, err
	//}
	//// insert 成功。
	//return true, nil
}

func AddToRDBWhenFollow(userId int, targetId int) {
	// 当a关注b时，redis的三个关注数据库会有以下操作
	redis.UserFollowings.SAdd(redis.Ctx, strconv.Itoa(userId), -1)
	redis.UserFollowings.SAdd(redis.Ctx, strconv.Itoa(userId), targetId)

	redis.UserFollowers.SAdd(redis.Ctx, strconv.Itoa(targetId), -1)
	redis.UserFollowers.SAdd(redis.Ctx, strconv.Itoa(targetId), userId)

	// 如果此时b也关注了a,说明b的互关用户有a,然后a的互关对象也有b
	if flag, _ := redis.UserFollowings.SIsMember(redis.Ctx, strconv.Itoa(targetId), userId).Result(); flag {
		redis.UserFriends.SAdd(redis.Ctx, strconv.Itoa(userId), targetId)
		redis.UserFriends.SAdd(redis.Ctx, strconv.Itoa(targetId), userId)
	}
}

// CancelFollowAction 取关操作的业务
func (followService *FollowServiceImp) CancelFollowAction(userId int64, targetId int64) (bool, error) {

	DelToRDBWhenCancelFollow(int(userId), int(targetId))

	// 获取取关的消息队列
	followDelMQ := rabbitmq.SimpleFollowDelMQ
	followDao := dao.NewFollowDaoInstance()
	follow, err := followDao.FindEverFollowing(userId, targetId)
	// 寻找 SQL 出错。
	if nil != err {
		return false, err
	}
	// 曾经关注过，只需要update一下cancel即可。
	if nil != follow {
		err := followDelMQ.PublishSimpleFollow(fmt.Sprintf("%d-%d-%s", userId, targetId, "update"))
		if err != nil {
			return false, err
		}
		return true, nil
		//_, err := followDao.UpdateFollowRelation(userId, targetId, 0)
		//// update 出错。
		//if nil != err {
		//	return false, err
		//}
		//// update 成功。
		//return true, nil
	}
	// 没有关注关系
	return false, nil
}
func DelToRDBWhenCancelFollow(userId int, targetId int) {
	// 当a取关b时，redis的三个关注数据库会有以下操作
	redis.UserFollowings.SRem(redis.Ctx, strconv.Itoa(userId), targetId)

	redis.UserFollowers.SRem(redis.Ctx, strconv.Itoa(targetId), userId)

	// a取关b，如果a和b属于互关的用户，则两者的互关记录都会删除
	redis.UserFriends.SRem(redis.Ctx, strconv.Itoa(userId), targetId)
	redis.UserFriends.SRem(redis.Ctx, strconv.Itoa(targetId), userId)
}

// GetFollowings 获取正在关注的用户详情列表业务
func (followService *FollowServiceImp) GetFollowings(userId int64) ([]User, error) {
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
func (followService *FollowServiceImp) GetFollowers(userId int64) ([]User, error) {
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

// GetFriends 获取用户好友列表（附带与其最新聊天记录）
func (followService *FollowServiceImp) GetFriends(userId int64) ([]FriendUser, error) {
	followDao := dao.NewFollowDaoInstance()
	// 关注用户的id List和count
	userFriendId, userFriendCnt, err := followDao.GetFollowingsInfo(userId)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	fmt.Print(userFriendId)
	var userFriends []FriendUser

	userFollowings, err1 := followService.GetFollowings(userId)
	if err1 != nil {
		log.Println(err1.Error())
		return nil, err1
	}

	for i := 0; int64(i) < userFriendCnt; i++ {
		var friendUserTemp FriendUser
		//使用消息模块服务
		msi := messageServiceImpl
		messageInfo, err := msi.LatestMessage(userId, userFriendId[i])
		//没有发生过聊天，不返回
		if err != nil {
			continue
		}

		friendUserTemp.Id = userFollowings[i].Id
		friendUserTemp.Name = userFollowings[i].Name
		friendUserTemp.FollowerCount = userFollowings[i].FollowerCount
		friendUserTemp.FollowCount = userFollowings[i].FollowCount
		friendUserTemp.Avatar = config.CUSTOM_DOMAIN + config.OSS_USER_AVATAR_DIR
		// 传入当前登陆用户id-userId和好友id-userFriendsId 得到最新聊天消息及其类型
		friendUserTemp.Message = messageInfo.message
		friendUserTemp.MsgType = messageInfo.msgType
		userFriends = append(userFriends, friendUserTemp)
	}
	return userFriends, nil
}

// GetFollowingCnt 加入redis 根据用户id查询关注数
func (followService *FollowServiceImp) GetFollowingCnt(userId int64) (int64, error) {
	//followDao := dao.NewFollowDaoInstance()
	//return followDao.GetFollowingCnt(userId)

	if cnt, err := redis.UserFollowings.SCard(redis.Ctx, strconv.Itoa(int(userId))).Result(); cnt > 0 {
		redis.UserFollowings.Expire(redis.Ctx, strconv.Itoa(int(userId)), config.ExpireTime)
		return cnt - 1, err
	}

	followDao := dao.NewFollowDaoInstance()
	ids, _, err := followDao.GetFollowingsInfo(userId)

	if err != nil {
		return 0, err
	}

	go ImportToRDBFollowing(int(userId), ids)

	return int64(len(ids)), nil
}

// ImportToRDBFollowing 将登陆用户的关注id列表导入到following数据库中
func ImportToRDBFollowing(userId int, ids []int64) {
	redis.UserFollowings.SAdd(redis.Ctx, strconv.Itoa(userId), -1)

	for id := range ids {
		redis.UserFollowings.SAdd(redis.Ctx, strconv.Itoa(userId), int(id))
	}

	redis.UserFollowings.Expire(redis.Ctx, strconv.Itoa(userId), config.ExpireTime)
}

// GetFollowerCnt 根据用户id查询粉丝数
func (followService *FollowServiceImp) GetFollowerCnt(userId int64) (int64, error) {
	//followDao := dao.NewFollowDaoInstance()
	//return followDao.GetFollowerCnt(userId)

	if cnt, err := redis.UserFollowers.SCard(redis.Ctx, strconv.Itoa(int(userId))).Result(); cnt > 0 {
		redis.UserFollowers.Expire(redis.Ctx, strconv.Itoa(int(userId)), config.ExpireTime)
		return cnt - 1, err
	}

	followDao := dao.NewFollowDaoInstance()
	ids, _, err := followDao.GetFollowersInfo(userId)

	if err != nil {
		return 0, err
	}

	go ImportToRDBFollower(int(userId), ids)

	return int64(len(ids)), nil
}

// ImportToRDBFollower 将登陆用户的关注id列表导入到follower数据库中
func ImportToRDBFollower(userId int, ids []int64) {
	redis.UserFollowings.SAdd(redis.Ctx, strconv.Itoa(userId), -1)

	for id := range ids {
		redis.UserFollowings.SAdd(redis.Ctx, strconv.Itoa(userId), int(id))
	}

	redis.UserFollowings.Expire(redis.Ctx, strconv.Itoa(userId), config.ExpireTime)
}

// CheckIsFollowing 判断当前登录用户是否关注了目标用户
func (followService *FollowServiceImp) CheckIsFollowing(userId int64, targetId int64) (bool, error) {
	//followDao := dao.NewFollowDaoInstance()
	//return followDao.FindFollowRelation(userId, targetId)

	if flag, err := redis.UserFollowings.SIsMember(redis.Ctx, strconv.Itoa(int(userId)), targetId).Result(); flag {
		if err != nil {
			return false, err
		} else {
			return true, nil
		}
	}

	// 该键有效说明是没有关注
	if cnt, err := redis.UserFollowings.SCard(redis.Ctx, strconv.Itoa(int(userId))).Result(); cnt > 0 {
		if err != nil {
			return false, err
		}

		redis.UserFollowings.Expire(redis.Ctx, strconv.Itoa(int(userId)), config.ExpireTime)
		return false, nil
	}

	// 该键无效，导入
	followDao := dao.NewFollowDaoInstance()
	ids, _, err := followDao.GetFollowingsInfo(userId)

	if err != nil {
		return false, err
	}

	go ImportToRDBFollowing(int(userId), ids)

	return followDao.FindFollowRelation(userId, targetId)

}
