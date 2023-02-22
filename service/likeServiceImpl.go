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
	"x-tiktok/util"
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
	islike, err := dao.IsVideoLikedByUser(userId, videoId)
	log.Print("islike:", islike)
	log.Println("actionType:", actionType)
	// 获取点赞和取消点赞的消息队列
	likeAddMQ := rabbitmq.SimpleLikeAddMQ
	likeDelMQ := rabbitmq.SimpleLikeDelMQ
	if islike == -1 {
		//用户没有点赞过该视频
		//插入一条新记录
		//var likeinfo dao.Like
		//likeinfo.UserId = userId
		//likeinfo.VideoId = videoId
		//likeinfo.Liked = int8(actionType)
		//likeinfo.CreatedAt = time.Now()
		//likeinfo.UpdatedAt = time.Now()
		//err = dao.InsertLikeInfo(likeinfo)
		//return nil
		//  更新 redis
		syncLikeRedis(userId, videoId, 1)
		// 消息队列
		err := likeAddMQ.PublishSimple(fmt.Sprintf("%d-%d-%s", userId, videoId, "insert"))
		return err
	}
	//该用户曾对此视频点过赞
	//err = dao.UpdateLikeInfo(userId, videoId, int8(actionType))
	if actionType == 1 {
		syncLikeRedis(userId, videoId, 1)
		err = likeAddMQ.PublishSimple(fmt.Sprintf("%d-%d-%s", userId, videoId, "update"))
	} else {
		syncLikeRedis(userId, videoId, 2)
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

// GetLikesList 获取用户点赞视频列表
func (*LikeServiceImpl) GetLikesList(userId int64) ([]Video, error) {
	likedVideoIdList, _, err := GetLikeVideoIdListByRedis(userId)
	//likedVideoIdList, _, err := dao.GetLikeListByUserId(userId)
	if err != nil {
		log.Print("Get like list failed!")
		return nil, err
	}
	likeService := NewLikeServImpInstance()
	likedVideoInfoList, err := likeService.GetVideoListById(likedVideoIdList, userId)
	if err != nil {
		log.Println("Get videoList failed")
	}
	return likedVideoInfoList, nil
}

// GetUserLikeCount 获取用户点赞视频的数量
func (*LikeServiceImpl) GetUserLikeCount(userId int64) (int64, error) {
	likeCnt, err := GetUserLikeCountByRedis(userId)
	//_, likeCnt, err := dao.GetLikeListByUserId(userId)
	if err != nil {
		log.Print("Get like count failed!")
		return -1, err
	}
	return likeCnt, nil
}

// GetVideoLikedCount 获取视频点赞数
func (*LikeServiceImpl) GetVideoLikedCount(videoId int64) (int64, error) {
	likeCnt, err := GetVideoLikedCountByRedis(videoId)
	//likeCnt, err := dao.VideoLikedCount(videoId)
	if err != nil {
		log.Print("Get like count failed!")
		return -1, err
	}
	return likeCnt, nil
}

// IsLikedByUser 当前用户是否点赞该视频
func (*LikeServiceImpl) IsLikedByUser(userId int64, videoId int64) (bool, error) {
	liked, err := IsLikedByUserByRedis(userId, videoId)
	//liked, err := dao.IsLikedByUser(userId, videoId)
	if err != nil {
		return false, err
	}
	return liked, nil
}

// GetUserLikedCnt 计算用户被点赞的视频获赞总数
func (*LikeServiceImpl) GetUserLikedCnt(userId int64) (int64, error) {
	//likedVideoIdList, err := dao.GetUserVideoLikedByOther(userId)
	//if err != nil {
	//	return -1, nil
	//}
	//var count int64 = 0
	//for _, videoId := range likedVideoIdList {
	//	tmp, _ := dao.VideoLikedCount(videoId)
	//	count += tmp
	//}
	// 通过一条 sql 语句计算该用户的视频获赞总数，减少数据库访问次数
	count, err := dao.GetUserVideoLikedTotalCount(userId)
	if err != nil {
		return -1, err
	}
	return count, nil
}

// GetLikeVideoIdListByRedis 查询 Redis 中 userId 点赞过的视频 id 列表，存在缓存直接返回，不存在则从数据库中查询后，再返回
func GetLikeVideoIdListByRedis(userId int64) ([]int64, int64, error) {
	userIdStr := strconv.FormatInt(userId, 10)
	var likedVideoIdList = make([]int64, 0, config.VIDEO_INIT_NUM_PER_AUTHOR)
	keyCnt, err := redis.RdbUVid.Exists(redis.Ctx, userIdStr).Result()
	if err != nil {
		log.Println(err)
	}
	if keyCnt > 0 {
		// RdbUVid 存在 userId
		vIds, _ := redis.RdbUVid.SMembers(redis.Ctx, userIdStr).Result()
		likedVideoIdList, _ = util.StrArrToInt64Arr(vIds)
	} else {
		// 不存在这个 key，从数据库中导入用户最新点赞的视频数据到 redis
		videoIds, _, err := dao.GetLikeListByUserId(userId)
		if err != nil {
			log.Println(err)
		}
		ImportVideoIdsFromDb(userId, videoIds)
		likedVideoIdList = videoIds
	}
	return likedVideoIdList, int64(len(likedVideoIdList)), nil
}

// ImportVideoIdsFromDb 更新最新的用户点赞视频到 redis
func ImportVideoIdsFromDb(userId int64, videoIds []int64) error {
	userIdStr := strconv.FormatInt(userId, 10)
	for _, videoId := range videoIds {
		redis.RdbUVid.SAdd(redis.Ctx, userIdStr, videoId)
	}
	// 设置过期时间，为数据不一致情况兜底
	redis.RdbUVid.Expire(redis.Ctx, userIdStr, config.ExpireTime)
	return nil
}

// ImportUserIdsFromDb 更新点赞 videoId 视频的用户 id 到 Redis
func ImportUserIdsFromDb(videoId int64, userIds []int64) error {
	videoIdStr := strconv.FormatInt(videoId, 10)
	for _, userId := range userIds {
		redis.RdbVUid.SAdd(redis.Ctx, videoIdStr, userId)
	}
	// 设置过期时间，为数据不一致情况兜底
	redis.RdbVUid.Expire(redis.Ctx, videoIdStr, config.ExpireTime)
	return nil
}

// GetUserLikeCountByRedis 通过 Redis 获取用户点赞视频的数量
func GetUserLikeCountByRedis(userId int64) (int64, error) {
	userIdStr := strconv.FormatInt(userId, 10)
	keyCnt, err := redis.RdbUVid.Exists(redis.Ctx, userIdStr).Result()
	if err != nil {
		log.Println(err)
	}
	if keyCnt > 0 {
		// RdbUVid 存在 userId
		cnt, _ := redis.RdbUVid.SCard(redis.Ctx, userIdStr).Result()
		return cnt, nil
	} else {
		// 不存在这个 key，从数据库中导入用户最新点赞的视频数据到 redis
		videoIds, _, err := dao.GetLikeListByUserId(userId)
		if err != nil {
			log.Println(err)
		}
		ImportVideoIdsFromDb(userId, videoIds)
		return int64(len(videoIds)), nil
	}
}

// GetVideoLikedCountByRedis 通过 Redis 获取视频被点赞的数量
func GetVideoLikedCountByRedis(videoId int64) (int64, error) {
	videoIdStr := strconv.FormatInt(videoId, 10)
	keyCnt, err := redis.RdbVUid.Exists(redis.Ctx, videoIdStr).Result()
	if err != nil {
		log.Println(err)
	}
	if keyCnt > 0 {
		// RdbVUid 存在 videoId
		cnt, _ := redis.RdbVUid.SCard(redis.Ctx, videoIdStr).Result()
		return cnt, nil
	} else {
		// 不存在这个 key，从数据库中导入视频最新获赞数据到 redis
		userIdList, likeCnt, err := dao.UsersOfLikeVideo(videoId)
		if err != nil {
			log.Println(err)
			return -1, err
		}
		ImportUserIdsFromDb(videoId, userIdList)
		if err != nil {
			log.Println(err)
		}
		return likeCnt, nil
	}
}

// IsLikedByUserByRedis 通过 Redis 判断当前 videoId 是否被 userId 点赞
func IsLikedByUserByRedis(userId int64, videoId int64) (bool, error) {
	userIdStr := strconv.FormatInt(userId, 10)
	keyCnt, err := redis.RdbUVid.Exists(redis.Ctx, userIdStr).Result()
	if err != nil {
		log.Println(err)
	}
	if keyCnt > 0 {
		// RdbUVid 存在 userId
		result, err := redis.RdbUVid.SIsMember(redis.Ctx, userIdStr, videoId).Result()
		if err != nil {
			return false, err
		}
		return result, nil
	} else {
		// 不存在这个 key
		liked, _ := dao.IsLikedByUser(userId, videoId)
		if liked {
			redis.RdbUVid.SAdd(redis.Ctx, userIdStr, videoId)
		}
		return liked, nil
	}
}

// 点赞/取消时同步更新 redis 中的数据
func syncLikeRedis(userId int64, videoId int64, actionType int32) error {
	userIdStr := strconv.FormatInt(userId, 10)
	videoIdStr := strconv.FormatInt(videoId, 10)
	switch actionType {
	case 1:
		// 点赞
		redis.RdbUVid.SAdd(redis.Ctx, userIdStr, videoId)
		redis.RdbVUid.SAdd(redis.Ctx, videoIdStr, userId)
	case 2:
		// 取消点赞
		redis.RdbUVid.SRem(redis.Ctx, userIdStr, videoId)
		redis.RdbVUid.SRem(redis.Ctx, videoIdStr, userId)
	default:
		log.Println("syncLikeRedis 传入的 ActionType 参数错误")
	}
	return nil
}
