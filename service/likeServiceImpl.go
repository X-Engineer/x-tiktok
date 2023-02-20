package service

import (
	"errors"
	"fmt"
	redis2 "github.com/redis/go-redis/v9"
	"log"
	"strconv"
	"sync"
	"time"
	"x-tiktok/config"
	"x-tiktok/dao"
	"x-tiktok/middleware/rabbitmq"
	"x-tiktok/middleware/redis"
)

type LikeServiceImpl struct {
	VideoService
}

var (
	likeServiceImp      *LikeServiceImpl
	likeServiceInstance sync.Once
)

func NewLikeServImpInstance() *LikeServiceImpl {
	likeServiceInstance.Do(func() {
		likeServiceImp = &LikeServiceImpl{
			VideoService: &VideoServiceImpl{},
		}
	})
	return likeServiceImp
}

func (*LikeServiceImpl) FavoriteAction(userId int64, videoId int64, actionType int32) error {
	//先看缓存里有无点赞数据
	liked, err := likeServiceImp.RdsIsLikedByUser(userId, videoId)
	if liked <= 1 {
		//缓存没有命中则去数据库查询
		islike, err1 := dao.IsVideoLikedByUser(userId, videoId)
		if err1 != nil {
			log.Print(err1.Error())
		}
		//将标记字段liked更新为数据库查询结果islike
		liked = islike
	}

	log.Print("islike:", liked)
	log.Println("actionType:", actionType)
	// 获取点赞和取消点赞的消息队列
	likeAddMQ := rabbitmq.SimpleLikeAddMQ
	likeDelMQ := rabbitmq.SimpleLikeDelMQ
	//根据是否有点赞记录来决定接下来的数据库操作，如果没有点赞记录则当前操作一定是点赞，则在like表中插入一条新的记录，执行insert命令
	if liked == -1 {
		//用户没有点赞过该视频
		//插入一条新记录
		// 消息队列
		err := likeAddMQ.PublishSimple(fmt.Sprintf("%d-%d-%s", userId, videoId, "insert"))
		return err
	}
	//数据库中有对应记录，说明该用户曾对此视频点过赞，则根据当前的点赞状态（actiontype为1：点赞，2：取消点赞）来更新数据库同时维护缓存数据
	if actionType == 1 {
		err = likeAddMQ.PublishSimple(fmt.Sprintf("%d-%d-%s", userId, videoId, "update"))
	} else {
		err = likeDelMQ.PublishSimple(fmt.Sprintf("%d-%d-%s", userId, videoId, "update"))
	}
	if err != nil {
		log.Print(err.Error() + "Favorite action failed!")
		return err
	} else {
		log.Print("Favorite action succeed!")
	}
	return nil
}

// GetLikesList 获取当前用户点赞的视频列表
func (*LikeServiceImpl) GetLikesList(userId int64) ([]Video, error) {
	//先从缓存中获取点赞的视频id列表
	videoIdList, err := likeServiceImp.RdsGetLikesList(userId)
	if videoIdList == nil {
		//将int64 videoId转换为 string strVideoId
		strUserId := strconv.FormatInt(userId, 10)
		//对键userId设置默认值，防止脏读
		if _, err := likeServiceImp.RedisSetDefaultValueAndExpireTime(redis.RdbLikeUserId, strUserId); err != nil {
			return nil, err
		}
		likedVideoIdList, err := dao.GetLikeListByUserId(userId)
		if err != nil {
			log.Print("Get like list failed!")
			return nil, err
		}
		videoIdList = likedVideoIdList
	}

	likeService := NewLikeServImpInstance()
	likedVideoInfoList, err := likeService.GetVideoListById(videoIdList, userId)
	if err != nil {
		log.Println("Get videoList failed")
	}
	return likedVideoInfoList, nil
}

// GetUserLikeCount 获取用户点赞数量
func (*LikeServiceImpl) GetUserLikeCount(userId int64) (int64, error) {
	cnt, err := likeServiceImp.RdsGetUserLikedCnt(userId)
	log.Printf("缓存查询成功！查询结果为：%d", cnt)
	if cnt < 0 {
		log.Printf("缓存查询失败，开始查数据库并更新缓存！")
		if err != nil {
			log.Print(err.Error())
			return -1, err
		}
		//将int64 userId转换为 string strUserId
		strUserId := strconv.FormatInt(userId, 10)
		//cnt值为-1时，说明strUserId无记录，则需加入默认值防脏读
		if cnt == -1 {
			log.Println("strUserId无记录，需加入默认值防脏读！")
			if _, err0 := likeServiceImp.RedisSetDefaultValueAndExpireTime(redis.RdbLikeUserId, strUserId); err0 != nil {
				return -1, err0
			}
		}

		//从数据库中查询点赞数并更新缓存信息
		likeCnt, err1 := dao.GetLikeCountByUser(userId)
		if err1 != nil {
			log.Print("Get like count failed!")
			return -1, err1
		}
		return likeCnt, nil
	} else {
		if err != nil {
			log.Print(err.Error())
			return 0, err
		} else {
			//命中缓存，成功查询则返回值即可
			return cnt, nil
		}
	}

}

// GetVideoLikedCount 获取视频点赞数
func (*LikeServiceImpl) GetVideoLikedCount(videoId int64) (int64, error) {
	//step1：先查询缓存中的videoId键值
	n, err := likeServiceImp.RdsGetVideoLikedCount(videoId)
	//没有命中缓存时n<0
	if n < 0 {
		//将int64 videoId转换为 string strVideoId
		strVideoId := strconv.FormatInt(videoId, 10)
		//step2：先对键videoId设置默认值，防止脏读
		likeServiceImp.RedisSetDefaultValueAndExpireTime(redis.RdbLikeVideoId, strVideoId)
		//step3：从数据库中查询点赞数并更新缓存信息
		likeCnt, err1 := dao.VideoLikedCount(videoId)
		if err1 != nil {
			log.Print("Get like count failed!")
			return -1, err1
		}
		return likeCnt, nil
	} else {
		if err != nil {
			log.Print(err.Error())
			return 0, err
		} else {
			//命中缓存，成功查询则返回值即可
			return n, nil
		}
	}
}

// IsLikedByUser 当前用户是否点赞该视频
func (*LikeServiceImpl) IsLikedByUser(userId int64, videoId int64) (bool, error) {
	liked, err := likeServiceImp.RdsIsLikedByUser(userId, videoId)
	if liked <= 1 {
		//缓存没有命中则去数据库查询
		islike, err1 := dao.IsVideoLikedByUser(userId, videoId)
		if err1 != nil {
			log.Print(err1.Error())
		}
		liked = islike
	}
	//liked, err := dao.IsVideoLikedByUser(userId, videoId)
	if err != nil {
		return false, err
	}
	if liked == config.LIKE {
		return true, nil
	}
	return false, nil
}

// GetUserLikedCnt逻辑错误(已修改）
// GetUserLikedCnt 计算用户的视频获赞总数
func (*LikeServiceImpl) GetUserLikedCnt(userId int64) (int64, error) {
	likeService := NewLikeServImpInstance()
	//这里应该是用户发表的视频列表而不是点赞的视频列表
	//likedVideoList, err := likeService.GetLikesList(userId)
	videoList, err := likeService.PublishList(userId)
	if err != nil {
		return -1, errors.New("获取用户发布的视频列表失败")
	}
	var count int64 = 0
	for _, video := range videoList {
		count += video.FavoriteCount
	}
	return count, nil
}

/***********************************************************************************************************************/
//Redis操作
//获取用户点赞数
func (i *LikeServiceImpl) RdsGetUserLikedCnt(userId int64) (int64, error) {
	//redis中存储类型是string，需要先把videoId转成string类型在进行查询
	strUserId := strconv.FormatInt(userId, 10)
	var LikeCnt int64
	//先从Redis缓存中获取用户点赞数
	if n, err := redis.RdbLikeUserCnt.Exists(redis.Ctx, strUserId).Result(); n > 0 {
		if err != nil {
			log.Printf("Redis：获取用户点赞数失败" + err.Error())
			return -1, err
		}
		cnt := redis.RdbLikeUserCnt.Get(redis.Ctx, strUserId).String()
		log.Printf("RdbLikeUserCnt: " + cnt)
		LikeCnt, _ = strconv.ParseInt(cnt, 10, 64)
		log.Printf("RdbLikeUserCnt缓存查询成功！")
	} else {
		log.Printf("RdbLikeUserCnt缓存查询失败！")
		count, err0 := dao.GetLikeCountByUser(userId)
		log.Printf("数据库获取用户点赞数: %d", count)

		if err0 != nil {
			log.Printf("Redis：RdbLikeUserCnt从数据库中获取用户点赞数失败" + err0.Error())
			return -1, err
		}
		//cnt := redis.RdbLikeUserCnt.Get(redis.Ctx, strUserId).String()
		//log.Printf("RdbLikeUserCnt new: " + cnt)
		LikeCnt = count
	}

	//如果key:strVideoId存在 则计算集合中userId个数
	n, err := redis.RdbLikeVideoId.Exists(redis.Ctx, strUserId).Result()
	if n > 0 {
		//如果有问题，说明查询redis失败,返回默认false,返回错误信息
		if err != nil {
			log.Printf("方法:FavouriteCount RedisLikeVideoId query key失败：%v", err)
			return 0, err
		}
		/*******************************************************************************/
		//数据一致性校验
		rdsLikeCnt, err0 := redis.RdbLikeUserId.SCard(redis.Ctx, strUserId).Result()
		if err0 != nil {
			log.Printf("Redis:RdbLikeUserId获取用户点赞数失败：%v", err0)
			return -1, err0
		}
		//将RdbLikeUserId中的用户点赞数和RdbLikeUserCnt的进行对比来保证数据一致性
		if rdsLikeCnt-1 != LikeCnt {
			log.Printf("用户点赞信息缓存数据不一致，需要从数据库中查询并更新数据！！")
			return -2, nil //返回值为-2说明数据不一致，则不需添加默认值，只要将缓存数据更新即可
		}
		/*******************************************************************************/

	} else {
		//n<=0说明缓存中没有userId，返回-1，执行数据库查询
		log.Printf("RdbLikeUserId缓存查询失败！")
		return -1, nil
	}
	//校验通过
	log.Printf("方法:FavouriteCount RedisLikeVideoId query count succeed")
	return LikeCnt, nil //去掉DefaultRedisValue
}

func (*LikeServiceImpl) RdsGetVideoLikedCount(videoId int64) (int64, error) {
	//redis中存储类型是string，需要先把videoId转成string类型在进行查询
	strVideoId := strconv.FormatInt(videoId, 10)
	var LikeCnt int64
	//先从Redis缓存中获取视频获赞数
	if n, err := redis.RdbLikeVideoCnt.Exists(redis.Ctx, strVideoId).Result(); n > 0 {
		if err != nil {
			log.Printf("Redis：获取用户点赞数失败" + err.Error())
			return -1, err
		}
		cnt := redis.RdbLikeVideoCnt.Get(redis.Ctx, strVideoId).String()
		LikeCnt, _ = strconv.ParseInt(cnt, 10, 64)
	}
	//如果key:strVideoId存在 则计算集合中userId个数
	n, err := redis.RdbLikeVideoId.Exists(redis.Ctx, strVideoId).Result()
	if n > 0 {
		//如果有问题，说明查询redis失败,返回默认false,返回错误信息
		if err != nil {
			log.Printf("方法:FavouriteCount RedisLikeVideoId query key失败：%v", err)
			return 0, err
		}
		/*******************************************************************************/
		//数据一致性校验
		rdsLikeCnt, err0 := redis.RdbLikeVideoId.SCard(redis.Ctx, strVideoId).Result()
		if err0 != nil {
			log.Printf("Redis:RdbLikeVideoId获取用户点赞数失败：%v", err0)
			return -1, err0
		}
		//将RdbLikeVideoId中的用户点赞数和RdbLikeVideoCnt的进行对比来保证数据一致性
		if rdsLikeCnt-1 != LikeCnt {
			log.Printf("视频获赞信息缓存数据不一致，需要从数据库中查询并更新数据！！")
			return -1, nil
		}
		/*******************************************************************************/
		//校验通过，返回值
		log.Printf("方法:FavouriteCount RedisLikeVideoId query count succeed")
		return rdsLikeCnt - 1, nil //去掉DefaultRedisValue
	}
	//n<=0说明缓存中没有videoId，返回-1，执行数据库查询
	return -1, err
}

func (i *LikeServiceImpl) RdsIsLikedByUser(userId int64, videoId int64) (int8, error) {
	//将int64 videoId转换为 string strVideoId
	strUserId := strconv.FormatInt(userId, 10)
	//将int64 videoId转换为 string strVideoId
	strVideoId := strconv.FormatInt(videoId, 10)

	//查询Redis缓存中 LikeUserId(key:strUserId)是否已经加载过此信息
	if n, err := redis.RdbLikeUserId.Exists(redis.Ctx, strUserId).Result(); n > 0 {
		//如果有问题，说明查询redis失败,返回错误信息
		if err != nil {
			log.Printf("方法:FavouriteAction RedisLikeUserId query key失败：%v", err)
			return -1, err
		} //如果加载过此信息key:strUserId，则加入value:videoId
		//如果redis LikeUserId 添加失败，数据库操作成功，会有脏数据，所以只有redis操作成功才执行数据库likes表操作
		if _, err1 := redis.RdbLikeUserId.SAdd(redis.Ctx, strUserId, videoId).Result(); err1 != nil {
			log.Printf("方法:FavouriteAction RedisLikeUserId add value失败：%v", err1)
			return -1, err1
		}
		return 1, nil
	} else {
		likeServiceImp.RedisSetDefaultValueAndExpireTime(redis.RdbLikeUserId, strUserId)
		return 0, nil
	}

	//查询Redis缓存中 LikeVideoId(key:strVideoId)是否已经加载过此信息
	if n, err := redis.RdbLikeVideoId.Exists(redis.Ctx, strVideoId).Result(); n > 0 {
		//如果有问题，说明查询redis失败,返回错误信息
		if err != nil {
			log.Printf("方法:FavouriteAction RedisLikeVideoId query key失败：%v", err)
			return -1, err
		} //如果加载过此信息key:strVideoId，则加入value:userId
		//如果redis LikeVideoId 添加失败，返回错误信息
		if _, err1 := redis.RdbLikeVideoId.SAdd(redis.Ctx, strVideoId, userId).Result(); err1 != nil {
			log.Printf("方法:FavouriteAction RedisLikeVideoId add value失败：%v", err1)
			return -1, err1
		}
		return 1, nil
	} else {
		likeServiceImp.RedisSetDefaultValueAndExpireTime(redis.RdbLikeVideoId, strVideoId)
		return 0, nil
	}
	//return 1, nil
}

func (i *LikeServiceImpl) RdsGetLikesList(userId int64) ([]int64, error) {
	//Redis中键类型为string，所以要先将userId转成string类型
	strUserId := strconv.FormatInt(userId, 10)
	var LikeCnt int64
	//先从Redis缓存中获取用户点赞数
	if n, err := redis.RdbLikeUserCnt.Exists(redis.Ctx, strUserId).Result(); n > 0 {
		if err != nil {
			log.Printf("Redis：获取用户点赞数失败" + err.Error())
			return nil, err
		}
		cnt := redis.RdbLikeUserCnt.Get(redis.Ctx, strUserId).String()
		LikeCnt, _ = strconv.ParseInt(cnt, 10, 64)
	}

	//查询缓存数据库LikeUserId,如果key：strUserId存在,则获取Set集合中全部videoId
	if n, err := redis.RdbLikeUserId.Exists(redis.Ctx, strUserId).Result(); n > 0 {
		if err != nil {
			log.Printf("方法:GetFavouriteList RedisLikeVideoId query key失败：%v", err)
			return nil, err
		}
		/*******************************************************************************/
		//数据一致性校验
		rdsLikeCnt, err0 := redis.RdbLikeUserId.SCard(redis.Ctx, strUserId).Result()
		if err0 != nil {
			log.Printf("Redis:RdbLikeUserId获取用户点赞数失败：%v", err0)
			return nil, err0
		}
		//将RdbLikeUserId中的用户点赞数和RdbLikeUserCnt的进行对比来保证数据一致性
		if rdsLikeCnt-1 != LikeCnt {
			log.Printf("用户点赞信息缓存数据不一致，需要从数据库中查询并更新数据！！")
			return nil, nil
		}
		/*******************************************************************************/

		//校验通过，从缓存中获取用户点赞列表信息
		//获取集合中全部videoId
		strVideoIdList, err1 := redis.RdbLikeUserId.SMembers(redis.Ctx, strUserId).Result()
		//如果有问题，说明查询redis失败,返回默认nil,返回错误信息
		if err1 != nil {
			log.Printf("方法:GetFavouriteList RedisLikeVideoId get values失败：%v", err1)
			return nil, err1
		}

		i := len(strVideoIdList) - 1
		var videoIdList []int64
		//videoIdList := make(int64, i)
		if i == 0 {
			return videoIdList, nil
		}
		for j := 0; j <= i; j++ {
			//将string videoId转换为 int64 VideoId
			videoId, _ := strconv.ParseInt(strVideoIdList[j], 10, 64)
			//如果是默认值则跳过
			if videoId == config.DefaultRedisValue {
				continue
			}
			videoIdList[j] = videoId
		}
		return videoIdList, nil
	}
	log.Printf("Redis：获取用户点赞信息失败，缓存未命中！")
	return nil, nil
}

func (i *LikeServiceImpl) RedisSetDefaultValueAndExpireTime(c *redis2.Client, key string) (bool, error) {
	//key:strVideoId，加入value:DefaultRedisValue,过期才会删，防止删最后一个数据的时候数据库还没更新完出现脏读，或者数据库操作失败造成的脏读
	if _, err := c.SAdd(redis.Ctx, key, config.DefaultRedisValue).Result(); err != nil {
		log.Printf(c.ClientGetName(redis.Ctx).String() + key + "add value失败")
		c.Del(redis.Ctx, key)
		return false, err
	}
	//给键值设置有效期，类似于gc机制
	_, err := c.Expire(redis.Ctx, key,
		time.Duration(config.OneMonth)*time.Second).Result()
	if err != nil {
		log.Printf(c.ClientGetName(redis.Ctx).String() + key + "设置有效期失败")
		c.Del(redis.Ctx, key)
		return false, err
	}
	return true, nil
}
