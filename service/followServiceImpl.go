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
	UserService
}

var (
	followServiceImp  *FollowServiceImp //controller层通过该实例变量调用service的所有业务方法。
	followServiceOnce sync.Once         //限定该service对象为单例，节约内存。
)

// RedisFollowPrefix 前缀
//var RedisFollowPrefix = "follow:"

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
		AddToRDBWhenFollow(userId, targetId)
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
	AddToRDBWhenFollow(userId, targetId)
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

func AddToRDBWhenFollow(userId int64, targetId int64) {
	followDao := dao.NewFollowDaoInstance()
	res, err := redis.UserFollowings.Get(redis.Ctx, strconv.FormatInt(userId, 10)).Result()
	fmt.Println(res)

	// 当前逻辑只判断了following数据库里是否存在userId键，不存在的话是把所有需要导入的数据库键全部导入，不知道这样做结果是否正确
	if err == redis.NilError {
		userFollowingsId, _, err1 := followDao.GetFollowingsInfo(userId)
		//获取target的粉丝，直接刷新，关注时刷新target的粉丝
		userFollowersId, _, err2 := followDao.GetFollowersInfo(targetId)

		userFriendsId1, _, err3 := followDao.GetFriendsInfo(userId)
		userFriendsId2, _, err4 := followDao.GetFriendsInfo(targetId)

		if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
			return
		}
		ImportToRDBFollowing(userId, userFollowingsId)
		ImportToRDBFollower(targetId, userFollowersId)
		ImportToRDBFriend(userId, userFriendsId1)
		ImportToRDBFriend(targetId, userFriendsId2)
	}

	// 当a关注b时，redis的三个关注数据库会有以下操作z
	redis.UserFollowings.SAdd(redis.Ctx, strconv.FormatInt(userId, 10), -1)
	redis.UserFollowings.SAdd(redis.Ctx, strconv.FormatInt(userId, 10), targetId)

	redis.UserFollowers.SAdd(redis.Ctx, strconv.FormatInt(targetId, 10), -1)
	redis.UserFollowers.SAdd(redis.Ctx, strconv.FormatInt(targetId, 10), userId)

	// 如果此时b也关注了a,说明b的互关用户有a,然后a的互关对象也有b
	if flag, _ := redis.UserFollowings.SIsMember(redis.Ctx, strconv.FormatInt(targetId, 10), userId).Result(); flag {
		redis.UserFriends.SAdd(redis.Ctx, strconv.FormatInt(userId, 10), -1)
		redis.UserFriends.SAdd(redis.Ctx, strconv.FormatInt(userId, 10), targetId)

		redis.UserFriends.SAdd(redis.Ctx, strconv.FormatInt(targetId, 10), -1)
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
	_, err := redis.UserFollowings.Get(redis.Ctx, strconv.FormatInt(userId, 10)).Result()
	if err == redis.NilError {
		userFollowingsId, userFollowingsCnt, err := followDao.GetFollowingsInfo(userId)
		if err != nil {
			log.Println(err.Error())
		}
		ImportToRDBFollowing(userId, userFollowingsId)
		return userFollowingsId, userFollowingsCnt, nil
	}
	redis.UserFollowings.SRem(redis.Ctx, strconv.FormatInt(userId, 10), -1)

	ids := redis.UserFollowings.SMembers(redis.Ctx, strconv.FormatInt(userId, 10)).Val()
	fmt.Println(ids)
	idsInt64, _ := convertToInt64Array(ids)

	redis.UserFollowings.SAdd(redis.Ctx, strconv.FormatInt(userId, 10), -1)
	return idsInt64, int64(len(idsInt64)), nil
}

// GetFollowings 获取正在关注的用户详情列表业务
func (followService *FollowServiceImp) GetFollowings(userId int64) ([]User, error) {
	//followDao := dao.NewFollowDaoInstance()
	//userFollowingsId, userFollowingsCnt, err := followDao.GetFollowingsInfo(userId)

	//  这里我想用上缓存去获取，但是下面得用list去声明，不能把长度写死 make([]User, userFollowingsCnt)。因为会有id 为-1的情况
	userFollowingsId, userFollowingsCnt, err := GetFollowingsByRedis(userId)
	if nil != err {
		log.Println(err.Error())
		return nil, err
	}

	userFollowings := make([]User, userFollowingsCnt)

	err = followService.BuildUser(userId, userFollowings, userFollowingsId, 0)

	if nil != err {
		log.Println(err.Error())
		return nil, err
	}

	return userFollowings, nil
}

/*
	获取粉丝列表业务
*/

// GetFollowersByRedis 从redis中获取用户粉丝列表
func GetFollowersByRedis(userId int64) ([]int64, int64, error) {
	followDao := dao.NewFollowDaoInstance()
	_, err := redis.UserFollowers.Get(redis.Ctx, strconv.FormatInt(userId, 10)).Result()
	if err == redis.NilError {
		userFollowersId, userFollowersCnt, err := followDao.GetFollowersInfo(userId)
		if err != nil {
			log.Println(err.Error())
		}
		ImportToRDBFollower(userId, userFollowersId)
		return userFollowersId, userFollowersCnt, nil
	}
	redis.UserFollowers.SRem(redis.Ctx, strconv.FormatInt(userId, 10), -1)

	ids := redis.UserFollowers.SMembers(redis.Ctx, strconv.FormatInt(userId, 10)).Val()
	idsInt64, _ := convertToInt64Array(ids)

	redis.UserFollowings.SAdd(redis.Ctx, strconv.FormatInt(userId, 10), -1)
	return idsInt64, int64(len(idsInt64)), nil
}

// GetFollowers 获取粉丝详情列表业务
func (followService *FollowServiceImp) GetFollowers(userId int64) ([]User, error) {
	//followDao := dao.NewFollowDaoInstance()
	//
	//userFollowersId, userFollowersCnt, err := followDao.GetFollowersInfo(userId)

	userFollowersId, userFollowersCnt, err := GetFollowersByRedis(userId)

	if nil != err {
		log.Println(err.Error())
		return nil, err
	}

	userFollowers := make([]User, userFollowersCnt)

	err = followService.BuildUser(userId, userFollowers, userFollowersId, 1)

	if nil != err {
		log.Println(err.Error())
		return nil, err
	}

	return userFollowers, nil

}

/*
	获取用户好友列表业务
*/

// 从redis中获取好友信息
func GetFriendsByRedis(userId int64) ([]int64, int64, error) {
	followDao := dao.NewFollowDaoInstance()
	_, err := redis.UserFriends.Get(redis.Ctx, strconv.FormatInt(userId, 10)).Result()
	if err == redis.NilError {
		userFriendsId, userFriendsCnt, err := followDao.GetFriendsInfo(userId)
		if err != nil {
			log.Println(err.Error())
		}
		ImportToRDBFriend(userId, userFriendsId)

		// 从mysql读的数据不会有脏读的数据直接返回
		return userFriendsId, userFriendsCnt, nil
	}

	// redis读的数据有脏读数据，注意剔除-1
	redis.UserFriends.SRem(redis.Ctx, strconv.FormatInt(userId, 10), -1)

	ids := redis.UserFriends.SMembers(redis.Ctx, strconv.FormatInt(userId, 10)).Val()
	idsInt64, _ := convertToInt64Array(ids)

	redis.UserFriends.SAdd(redis.Ctx, strconv.FormatInt(userId, 10), -1)
	return idsInt64, int64(len(idsInt64)), nil
}

// GetFriends 获取用户好友列表（附带与其最新聊天记录）
func (followService *FollowServiceImp) GetFriends(userId int64) ([]FriendUser, error) {
	//followDao := dao.NewFollowDaoInstance()
	//// 关注用户的id List和count
	//userFriendId, userFriendCnt, err := followDao.GetFollowingsInfo(userId)

	userFriendId, userFriendCnt, err := GetFriendsByRedis(userId)

	if nil != err {
		log.Println(err.Error())
		return nil, err
	}

	// fmt.Print(userFriendId)
	userFriends := make([]FriendUser, userFriendCnt)

	// userFollowings, err1 := followService.GetFollowings(userId)

	//if err1 != nil {
	//	log.Println(err1.Error())
	//	return nil, err1
	//}

	err = followService.BuildFriendUser(userId, userFriends, userFriendId)

	return userFriends, nil
}

/*
	对外提供服务之返回登陆用户的关注用户数量
*/

// GetFollowingCnt 加入redis 根据用户id查询关注数
func (followService *FollowServiceImp) GetFollowingCnt(userId int64) (int64, error) {
	//followDao := dao.NewFollowDaoInstance()
	//return followDao.GetFollowingCnt(userId)
	//redis.InitRedis()
	redis.UserFollowings.SAdd(redis.Ctx, strconv.FormatInt(userId, 10), -1)
	cnt, err := redis.UserFollowings.SCard(redis.Ctx, strconv.FormatInt(userId, 10)).Result()
	if err != nil {
		log.Println(err.Error())
	}
	if cnt > 0 {
		redis.UserFollowings.Expire(redis.Ctx, strconv.Itoa(int(userId)), config.ExpireTime)
		return cnt - 1, err
	}

	followDao := dao.NewFollowDaoInstance()
	ids, _, err := followDao.GetFollowingsInfo(userId)

	if err != nil {
		return 0, err
	}

	go ImportToRDBFollowing(userId, ids)

	return int64(len(ids)), nil
}

/*
	对外提供服务之返回登陆用户的粉丝用户数量
*/

// GetFollowerCnt 根据用户id查询粉丝数
func (followService *FollowServiceImp) GetFollowerCnt(userId int64) (int64, error) {
	//followDao := dao.NewFollowDaoInstance()
	//return followDao.GetFollowerCnt(userId)
	//redis.InitRedis()
	redis.UserFollowers.SAdd(redis.Ctx, strconv.FormatInt(userId, 10), -1)
	if cnt, err := redis.UserFollowers.SCard(redis.Ctx, strconv.Itoa(int(userId))).Result(); cnt > 0 {
		redis.UserFollowers.Expire(redis.Ctx, strconv.Itoa(int(userId)), config.ExpireTime)
		return cnt - 1, err
	}

	followDao := dao.NewFollowDaoInstance()
	ids, _, err := followDao.GetFollowersInfo(userId)

	if err != nil {
		return 0, err
	}

	go ImportToRDBFollower(userId, ids)

	return int64(len(ids)), nil
}

/*
	对外提供服务之返回登陆用户是否关注目标用户的布尔值
*/

// CheckIsFollowing 判断当前登录用户是否关注了目标用户
func (followService *FollowServiceImp) CheckIsFollowing(userId int64, targetId int64) (bool, error) {
	//followDao := dao.NewFollowDaoInstance()
	//return followDao.FindFollowRelation(userId, targetId)

	if flag, err := redis.UserFollowings.SIsMember(redis.Ctx, strconv.Itoa(int(userId)), targetId).Result(); flag {
		redis.UserFollowings.SAdd(redis.Ctx, strconv.FormatInt(userId, 10), -1)
		if err != nil {
			return false, err
		} else {
			return true, nil
		}
	}

	// 该键有效说明是没有关注
	if cnt, err := redis.UserFollowings.SCard(redis.Ctx, strconv.Itoa(int(userId))).Result(); cnt > 0 {
		redis.UserFollowings.SAdd(redis.Ctx, strconv.FormatInt(userId, 10), -1)
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

	go ImportToRDBFollowing(userId, ids)

	return followDao.FindFollowRelation(userId, targetId)

}

/*
	提供目标用户id和对应的id列表导入到redis中的方法，一般用在更新失效键的逻辑中
*/

// ImportToRDBFollowing 将登陆用户的关注id列表导入到following数据库中
func ImportToRDBFollowing(userId int64, ids []int64) {
	redis.UserFollowings.SAdd(redis.Ctx, strconv.FormatInt(userId, 10), -1)

	for _, id := range ids {
		redis.UserFollowings.SAdd(redis.Ctx, strconv.FormatInt(userId, 10), int(id))
	}

	redis.UserFollowings.Expire(redis.Ctx, strconv.FormatInt(userId, 10), config.ExpireTime)
}

// ImportToRDBFollower 将登陆用户的关注id列表导入到follower数据库中
func ImportToRDBFollower(userId int64, ids []int64) {
	redis.UserFollowers.SAdd(redis.Ctx, strconv.FormatInt(userId, 10), -1)

	for _, id := range ids {
		redis.UserFollowers.SAdd(redis.Ctx, strconv.FormatInt(userId, 10), int(id))
	}

	redis.UserFollowers.Expire(redis.Ctx, strconv.FormatInt(userId, 10), config.ExpireTime)
}

func ImportToRDBFriend(userId int64, ids []int64) {
	redis.UserFriends.SAdd(redis.Ctx, strconv.FormatInt(userId, 10), -1)

	for _, id := range ids {
		redis.UserFriends.SAdd(redis.Ctx, strconv.FormatInt(userId, 10), int(id))
	}

	redis.UserFriends.Expire(redis.Ctx, strconv.FormatInt(userId, 10), config.ExpireTime)
}

/*
	将返回关注用户、返回粉丝用户、返回好友用户中的构建用户的逻辑独立出来
	注： builduser方法根据传入的buildtype决定是构建关注用户还是粉丝用户
*/

// BuildUser 根据传入的id列表和空user数组，构建业务所需user数组并返回
func (followService *FollowServiceImp) BuildUser(userId int64, users []User, ids []int64, buildtype int) error {

	for i := 0; i < len(ids); i++ {

		users[i].Id = ids[i]

		// 这里非要使用user那边的接口！
		user, err1 := followService.GetUserLoginInfoById(ids[i])
		if err1 != nil {
			log.Fatal(err1)
		}
		users[i].Name = user.Name

		var err2 error
		users[i].FollowCount, err2 = followService.GetFollowingCnt(ids[i])
		if nil != err2 {
			log.Println(err2.Error())
			return err2
		}

		var err3 error
		users[i].FollowerCount, err3 = followService.GetFollowerCnt(ids[i])
		if nil != err3 {
			log.Println(err3.Error())
			return err3
		}

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

	for i := 0; i < len(ids); i++ {

		friendUsers[i].Id = ids[i]

		user, err1 := followService.GetUserLoginInfoById(ids[i])
		if err1 != nil {
			log.Fatal(err1)
		}
		friendUsers[i].Name = user.Name

		var err2 error
		friendUsers[i].FollowCount, err2 = followService.GetFollowingCnt(ids[i])
		if nil != err2 {
			log.Println(err2.Error())
			return err2
		}

		var err3 error
		friendUsers[i].FollowerCount, err3 = followService.GetFollowerCnt(ids[i])
		if nil != err3 {
			log.Println(err3.Error())
			return err3
		}

		friendUsers[i].IsFollow = true
		friendUsers[i].Avatar = config.CUSTOM_DOMAIN + config.OSS_USER_AVATAR_DIR

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
