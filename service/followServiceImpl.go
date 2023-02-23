package service

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"
	"x-tiktok/config"
	"x-tiktok/dao"
	"x-tiktok/middleware/rabbitmq"
	"x-tiktok/middleware/redis"
)

// FollowServiceImp 该结构体继承FollowService接口。
type FollowServiceImp struct {
	//MessageService
	FollowService
	UserService
}

var (
	followServiceImp  *FollowServiceImp //controller层通过该实例变量调用service的所有业务方法。
	followServiceOnce sync.Once         //限定该service对象为单例，节约内存。
)

func CacheTimeGenerator() time.Duration {
	// 先设置随机数 - 这里比较重要
	rand.Seed(time.Now().Unix())
	// 再设置缓存时间
	// 10 + [0~20) 分钟的随机时间
	return time.Duration((10 + rand.Int63n(20)) * int64(time.Minute))
}

func convertToInt64Array(strArr []string) ([]int64, error) {
	int64Arr := make([]int64, len(strArr))
	for i, str := range strArr {
		int64Val, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return nil, err
		}
		int64Arr[i] = int64Val
	}
	return int64Arr, nil
}

// NewFSIInstance 生成并返回FollowServiceImp结构体单例变量。
func NewFSIInstance() *FollowServiceImp {
	followServiceOnce.Do(
		func() {
			followServiceImp = &FollowServiceImp{
				UserService: &UserServiceImpl{},
			}
		})
	return followServiceImp
}

//-------------------------------------API IMPLEMENT--------------------------------------------

/*
	关注业务
*/

// FollowAction 关注操作的业务
func (followService *FollowServiceImp) FollowAction(userId int64, targetId int64) (bool, error) {

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
		followService.AddToRDBWhenFollow(userId, targetId)
		return true, nil

	}
	//发送消息队列
	err = followAddMQ.PublishSimpleFollow(fmt.Sprintf("%d-%d-%s", userId, targetId, "insert"))
	if err != nil {
		return false, err
	}
	followService.AddToRDBWhenFollow(userId, targetId)
	return true, nil

}

func (followService *FollowServiceImp) AddToRDBWhenFollow(userId int64, targetId int64) {
	followDao := dao.NewFollowDaoInstance()
	// 尝试给following数据库追加user关注target的记录
	keyCnt1, err1 := redis.UserFollowings.Exists(redis.Ctx, strconv.FormatInt(userId, 10)).Result()

	if err1 != nil {
		log.Println(err1.Error())
	}

	// 只判定键是否不存在，若不存在即从数据库导入
	if keyCnt1 <= 0 {
		userFollowingsId, _, err := followDao.GetFollowingsInfo(userId)
		if err != nil {
			log.Println(err.Error())
			return
		}
		ImportToRDBFollowing(userId, userFollowingsId)
	}
	// 数据库导入到redis结束后追加记录
	redis.UserFollowings.SAdd(redis.Ctx, strconv.FormatInt(userId, 10), targetId)

	// 尝试给follower数据库追加target的粉丝有user的记录
	keyCnt2, err2 := redis.UserFollowers.Exists(redis.Ctx, strconv.FormatInt(targetId, 10)).Result()

	if err2 != nil {
		log.Println(err2.Error())
	}

	if keyCnt2 <= 0 {
		//获取target的粉丝，直接刷新，关注时刷新target的粉丝
		userFollowersId, _, err := followDao.GetFollowersInfo(targetId)
		if err != nil {
			log.Println(err.Error())
			return
		}
		ImportToRDBFollower(targetId, userFollowersId)
	}

	redis.UserFollowers.SAdd(redis.Ctx, strconv.FormatInt(targetId, 10), userId)

	// 进行好友的判定，本接口实现user对target的关注，若此时target也关注了user，进行friend数据库的记录追加
	// user的好友有target，target的好友有user
	if flag, _ := followService.CheckIsFollowing(targetId, userId); flag {
		// 尝试给friend数据库追加user的好友有target的记录
		keyCnt3, err3 := redis.UserFriends.Exists(redis.Ctx, strconv.FormatInt(userId, 10)).Result()

		if err3 != nil {
			log.Println(err3.Error())
		}
		if keyCnt3 <= 0 {
			userFriendsId1, _, err := followDao.GetFriendsInfo(userId)
			if err != nil {
				log.Println(err)
				return
			}
			ImportToRDBFriend(userId, userFriendsId1)
		}

		redis.UserFriends.SAdd(redis.Ctx, strconv.FormatInt(userId, 10), targetId)

		// 尝试给friend数据库追加target的好友有user的记录
		keyCnt4, err4 := redis.UserFriends.Exists(redis.Ctx, strconv.FormatInt(targetId, 10)).Result()

		if err4 != nil {
			log.Println(err4.Error())
		}
		if keyCnt4 <= 0 {
			//获取target的粉丝，直接刷新，关注时刷新target的粉丝
			userFriendsId2, _, err := followDao.GetFriendsInfo(targetId)
			if err != nil {
				log.Println(err)
				return
			}
			ImportToRDBFriend(targetId, userFriendsId2)
		}

		redis.UserFriends.SAdd(redis.Ctx, strconv.FormatInt(targetId, 10), userId)
	}
}

/*
	取关业务
*/

// CancelFollowAction 取关操作的业务
func (followService *FollowServiceImp) CancelFollowAction(userId int64, targetId int64) (bool, error) {

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
		DelToRDBWhenCancelFollow(userId, targetId)
		return true, nil

	}
	// 没有关注关系
	return false, nil
}
func DelToRDBWhenCancelFollow(userId int64, targetId int64) {
	// 当a取关b时，redis的三个关注数据库会有以下操作
	redis.UserFollowings.SRem(redis.Ctx, strconv.FormatInt(userId, 10), targetId)

	redis.UserFollowers.SRem(redis.Ctx, strconv.FormatInt(targetId, 10), userId)

	// a取关b，如果a和b属于互关的用户，则两者的互关记录都会删除
	redis.UserFriends.SRem(redis.Ctx, strconv.FormatInt(userId, 10), targetId)
	redis.UserFriends.SRem(redis.Ctx, strconv.FormatInt(targetId, 10), userId)
}

/*
	获取关注列表业务
*/

// GetFollowingsByRedis 从redis获取登陆用户关注列表
func GetFollowingsByRedis(userId int64) ([]int64, int64, error) {
	followDao := dao.NewFollowDaoInstance()
	// 判定键是否存在
	keyCnt, err := redis.UserFollowings.Exists(redis.Ctx, strconv.FormatInt(userId, 10)).Result()

	if err != nil {
		log.Println(err.Error())
	}

	// 若键存在，获取缓存数据后返回
	if keyCnt > 0 {
		ids := redis.UserFollowings.SMembers(redis.Ctx, strconv.FormatInt(userId, 10)).Val()
		idsInt64, _ := convertToInt64Array(ids)

		return idsInt64, int64(len(idsInt64)), nil
	} else {
		// 键不存在，获取数据库数据，更新缓存并返回
		userFollowingsId, userFollowingsCnt, err1 := followDao.GetFollowingsInfo(userId)
		if err1 != nil {
			log.Println(err1.Error())
		}
		ImportToRDBFollowing(userId, userFollowingsId)
		return userFollowingsId, userFollowingsCnt, nil
	}

}

// GetFollowings 获取正在关注的用户详情列表业务
func (followService *FollowServiceImp) GetFollowings(userId int64) ([]User, error) {
	// 调用集成redis的关注用户获取接口获取关注用户id和关注用户数量
	userFollowingsId, userFollowingsCnt, err := GetFollowingsByRedis(userId)
	if nil != err {
		log.Println(err.Error())
	}

	// 根据关注用户数量创建空用户结构体数组
	userFollowings := make([]User, userFollowingsCnt)

	// 传入buildtype调用用户构建函数构建关注用户数组
	err1 := followService.BuildUser(userId, userFollowings, userFollowingsId, 0)

	if nil != err1 {
		log.Println(err1.Error())
	}

	return userFollowings, nil
}

/*
	获取粉丝列表业务
*/

// GetFollowersByRedis 从redis中获取用户粉丝列表
func GetFollowersByRedis(userId int64) ([]int64, int64, error) {
	followDao := dao.NewFollowDaoInstance()
	keyCnt, err := redis.UserFollowers.Exists(redis.Ctx, strconv.FormatInt(userId, 10)).Result()

	if err != nil {
		log.Println(err.Error())
	}

	if keyCnt > 0 {
		// 键存在，获取键中集合元素
		ids := redis.UserFollowers.SMembers(redis.Ctx, strconv.FormatInt(userId, 10)).Val()
		idsInt64, _ := convertToInt64Array(ids)

		return idsInt64, int64(len(idsInt64)), nil
	} else {
		// 键不存在，获取数据库数据更新至redis，返回数据库所获取数据
		userFollowersId, userFollowersCnt, err1 := followDao.GetFollowersInfo(userId)
		if err1 != nil {
			log.Println(err1.Error())
		}
		ImportToRDBFollower(userId, userFollowersId)
		return userFollowersId, userFollowersCnt, nil
	}

}

// GetFollowers 获取粉丝详情列表业务
func (followService *FollowServiceImp) GetFollowers(userId int64) ([]User, error) {
	// 调用集成redis的粉丝获取接口获取粉丝id和粉丝数量
	userFollowersId, userFollowersCnt, err := GetFollowersByRedis(userId)

	if nil != err {
		log.Println(err.Error())
	}

	// 根据粉丝数量创建空用户结构体数组
	userFollowers := make([]User, userFollowersCnt)

	// 传入buildtype调用用户构建函数构建粉丝数组
	err1 := followService.BuildUser(userId, userFollowers, userFollowersId, 1)

	if nil != err1 {
		log.Println(err1.Error())
	}

	return userFollowers, nil

}

/*
	获取用户好友列表业务
*/

// 从redis中获取好友信息
func GetFriendsByRedis(userId int64) ([]int64, int64, error) {
	followDao := dao.NewFollowDaoInstance()
	keyCnt, err := redis.UserFriends.Exists(redis.Ctx, strconv.FormatInt(userId, 10)).Result()

	if err != nil {
		log.Println(err.Error())
	}

	if keyCnt > 0 {
		// 键存在，获取键中集合元素
		ids := redis.UserFriends.SMembers(redis.Ctx, strconv.FormatInt(userId, 10)).Val()
		idsInt64, _ := convertToInt64Array(ids)

		return idsInt64, int64(len(idsInt64)), nil

	} else {
		// 键不存在，获取数据库数据更新至redis，返回数据库所获取数据
		userFriendsId, userFriendsCnt, err1 := followDao.GetFriendsInfo(userId)
		if err1 != nil {
			log.Println(err1.Error())
		}
		ImportToRDBFriend(userId, userFriendsId)

		return userFriendsId, userFriendsCnt, nil
	}

}

// GetFriends 获取用户好友列表（附带与其最新聊天记录）
func (followService *FollowServiceImp) GetFriends(userId int64) ([]FriendUser, error) {
	// 调用集成redis的好友获取接口获取好友id和好友数量
	userFriendId, userFriendCnt, err := GetFriendsByRedis(userId)

	if nil != err {
		log.Println(err.Error())
	}

	// 使用好友数量创建空好友结构体数组
	userFriends := make([]FriendUser, userFriendCnt)

	// 调用好友构建函数构建好友数组
	err1 := followService.BuildFriendUser(userId, userFriends, userFriendId)

	if err1 != nil {
		log.Println(err1.Error())
	}

	return userFriends, nil
}

/*
	对外提供服务之返回登陆用户的关注用户数量
*/

// GetFollowingCnt 加入redis 根据用户id查询关注数
func (followService *FollowServiceImp) GetFollowingCnt(userId int64) (int64, error) {
	followDao := dao.NewFollowDaoInstance()

	keyCnt, err := redis.UserFollowings.Exists(redis.Ctx, strconv.FormatInt(userId, 10)).Result()

	if err != nil {
		log.Println(err.Error())
	}

	if keyCnt > 0 {
		// 键存在，获取键中集合元素个数
		cnt, err2 := redis.UserFollowings.SCard(redis.Ctx, strconv.FormatInt(userId, 10)).Result()
		if err2 != nil {
			log.Println(err2.Error())
		}
		redis.UserFollowings.Expire(redis.Ctx, strconv.Itoa(int(userId)), CacheTimeGenerator())
		return cnt, nil

	} else {
		// 键不存在，获取数据库数据更新至redis，返回数据库所获取数据
		ids, _, err1 := followDao.GetFollowingsInfo(userId)

		if err1 != nil {
			log.Println(err1.Error())
		}

		ImportToRDBFollowing(userId, ids)

		return int64(len(ids)), nil
	}

}

/*
	对外提供服务之返回登陆用户的粉丝用户数量
*/

// GetFollowerCnt 根据用户id查询粉丝数
func (followService *FollowServiceImp) GetFollowerCnt(userId int64) (int64, error) {
	followDao := dao.NewFollowDaoInstance()

	keyCnt, err := redis.UserFollowers.Exists(redis.Ctx, strconv.FormatInt(userId, 10)).Result()

	if err != nil {
		log.Println(err.Error())
	}

	if keyCnt > 0 {
		// 键存在，获取键中集合元素个数
		cnt, err2 := redis.UserFollowers.SCard(redis.Ctx, strconv.Itoa(int(userId))).Result()

		if err2 != nil {
			log.Println(err2.Error())
		}

		redis.UserFollowers.Expire(redis.Ctx, strconv.Itoa(int(userId)), CacheTimeGenerator())
		return cnt, nil

	} else {
		// 键不存在，获取数据库数据更新至redis，返回数据库所获取数据
		ids, _, err1 := followDao.GetFollowersInfo(userId)

		if err1 != nil {
			log.Println(err1.Error())
		}

		ImportToRDBFollower(userId, ids)

		return int64(len(ids)), nil
	}

}

/*
	对外提供服务之返回登陆用户是否关注目标用户的布尔值
*/

// CheckIsFollowing 判断当前登录用户是否关注了目标用户
func (followService *FollowServiceImp) CheckIsFollowing(userId int64, targetId int64) (bool, error) {
	followDao := dao.NewFollowDaoInstance()

	keyCnt, err := redis.UserFollowings.Exists(redis.Ctx, strconv.FormatInt(userId, 10)).Result()

	if err != nil {
		log.Println(err.Error())
	}

	if keyCnt > 0 {
		// 键存在判断是否存在userId和targetId键值对
		flag, err3 := redis.UserFollowings.SIsMember(redis.Ctx, strconv.Itoa(int(userId)), targetId).Result()

		if err3 != nil {
			log.Println(err3)
		}

		if flag {
			return true, nil
		} else {
			return false, nil
		}
	} else {
		// 键不存在，获取数据库数据更新至redis中，使用dao层方法判断是否有关注关系
		ids, _, err1 := followDao.GetFollowingsInfo(userId)

		if err1 != nil {
			log.Println(err1)
		}

		ImportToRDBFollowing(userId, ids)

		isFollow, err2 := followDao.FindFollowRelation(userId, targetId)

		if err2 != nil {
			log.Println(err2)
		}

		return isFollow, nil
	}

}

/*
	提供目标用户id和对应的id列表导入到redis中的方法，一般用在更新失效键的逻辑中
*/

// ImportToRDBFollowing 将登陆用户的关注id列表导入到following数据库中
func ImportToRDBFollowing(userId int64, ids []int64) {
	// 将传入的userId及其关注用户id更新至redis中
	for _, id := range ids {
		redis.UserFollowings.SAdd(redis.Ctx, strconv.FormatInt(userId, 10), int(id))
	}

	redis.UserFollowings.Expire(redis.Ctx, strconv.FormatInt(userId, 10), CacheTimeGenerator())
}

// ImportToRDBFollower 将登陆用户的关注id列表导入到follower数据库中
func ImportToRDBFollower(userId int64, ids []int64) {
	// 将传入的userId及其粉丝id更新至redis中
	for _, id := range ids {
		redis.UserFollowers.SAdd(redis.Ctx, strconv.FormatInt(userId, 10), int(id))
	}

	redis.UserFollowers.Expire(redis.Ctx, strconv.FormatInt(userId, 10), CacheTimeGenerator())
}

func ImportToRDBFriend(userId int64, ids []int64) {
	// 将传入的userId及其好友id更新至redis中
	for _, id := range ids {
		redis.UserFriends.SAdd(redis.Ctx, strconv.FormatInt(userId, 10), int(id))
	}

	redis.UserFriends.Expire(redis.Ctx, strconv.FormatInt(userId, 10), CacheTimeGenerator())
}

/*
	将返回关注用户、返回粉丝用户、返回好友用户中的构建用户的逻辑独立出来
	注： builduser方法根据传入的buildtype决定是构建关注用户还是粉丝用户
*/

// BuildUser 根据传入的id列表和空user数组，构建业务所需user数组并返回
func (followService *FollowServiceImp) BuildUser(userId int64, users []User, ids []int64, buildtype int) error {
	folowDao := dao.NewFollowDaoInstance()

	// 遍历传入的用户id，组成user结构体
	for i := 0; i < len(ids); i++ {

		// 用户id赋值
		users[i].Id = ids[i]

		// 用户name赋值
		var err1 error
		users[i].Name, err1 = folowDao.GetUserName(ids[i])
		if nil != err1 {
			log.Println(err1)
			return err1
		}

		// 用户关注数赋值
		var err2 error
		users[i].FollowCount, err2 = followService.GetFollowingCnt(ids[i])
		if nil != err2 {
			log.Println(err2.Error())
			return err2
		}

		// 用户粉丝数赋值
		var err3 error
		users[i].FollowerCount, err3 = followService.GetFollowerCnt(ids[i])
		if nil != err3 {
			log.Println(err3.Error())
			return err3
		}

		// 根据传入的buildtype决定是哪种业务的user构建
		if buildtype == 1 {
			// 粉丝用户的isfollow属性需要调用接口再确认一下
			users[i].IsFollow, _ = followService.CheckIsFollowing(userId, ids[i])
		} else {
			// 关注用户的isfollow属性确定是true
			users[i].IsFollow = true
		}

	}
	return nil
}

// BuildFriendUser 根据传入的id列表和空frienduser数组，构建业务所需frienduser数组并返回
func (followService *FollowServiceImp) BuildFriendUser(userId int64, friendUsers []FriendUser, ids []int64) error {

	msi := messageServiceImpl
	followDao := dao.NewFollowDaoInstance()

	// 遍历传入的好友id，组装好友user结构体
	for i := 0; i < len(ids); i++ {

		// 好友id赋值
		friendUsers[i].Id = ids[i]

		// 好友name赋值
		var err1 error
		friendUsers[i].Name, err1 = followDao.GetUserName(ids[i])
		if nil != err1 {
			log.Println(err1)
			return err1
		}

		// 好友关注数赋值
		var err2 error
		friendUsers[i].FollowCount, err2 = followService.GetFollowingCnt(ids[i])
		if nil != err2 {
			log.Println(err2.Error())
			return err2
		}

		// 好友粉丝数赋值
		var err3 error
		friendUsers[i].FollowerCount, err3 = followService.GetFollowerCnt(ids[i])
		if nil != err3 {
			log.Println(err3.Error())
			return err3
		}

		// 好友其他属性赋值
		friendUsers[i].IsFollow = true
		friendUsers[i].Avatar = config.CUSTOM_DOMAIN + config.OSS_USER_AVATAR_DIR

		// 调用message模块获取聊天记录
		messageInfo, err := msi.LatestMessage(userId, ids[i])

		//在根据id获取不到最新一条消息时，需要返回对应的id
		if err != nil {

			continue
		}

		friendUsers[i].Message = messageInfo.message
		friendUsers[i].MsgType = messageInfo.msgType
	}

	// 将空数组内属性构建完成即可，不用特意返回数组
	return nil
}
