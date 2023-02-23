package service

import (
	"encoding/json"
	"log"
	"math/rand"
	"sort"
	"strconv"
	"sync"
	"time"
	"x-tiktok/config"
	"x-tiktok/dao"
	"x-tiktok/middleware/rabbitmq"
	"x-tiktok/middleware/redis"
)

type CommentServiceImpl struct {
	UserService
}

var (
	commentServiceImpl *CommentServiceImpl
	commentServiceOnce sync.Once
)

func GetCommentServiceInstance() *CommentServiceImpl {
	commentServiceOnce.Do(func() {
		commentServiceImpl = &CommentServiceImpl{
			&UserServiceImpl{},
		}
	})
	return commentServiceImpl
}

func (commentService *CommentServiceImpl) CommentAction(comment dao.Comment) (Comment, error) {
	csi := GetCommentServiceInstance()
	commentRes, err := dao.InsertComment(comment)
	if err != nil {
		return Comment{}, err
	}
	user, err := csi.GetUserLoginInfoById(comment.UserId)
	if err != nil {
		log.Println(err.Error())
	}
	// 随机数生成种子
	rand.Seed(time.Now().Unix())
	commentData := Comment{
		Id:         commentRes.Id,
		User:       user,
		Content:    commentRes.Content,
		CreateDate: commentRes.CreatedAt.Format(config.GO_STARTER_TIME),
		LikeCount:  int64(rand.Intn(100)),
		TeaseCount: int64(rand.Intn(100)),
	}
	// redis操作：将发表的评论id存入redis
	go func() {
		insertRedisVCId(strconv.FormatInt(comment.VideoId, 10), strconv.FormatInt(commentRes.Id, 10), commentData)
		log.Println("commentAction save in redis")
	}()

	return commentData, nil
}

func (commentService *CommentServiceImpl) DeleteCommentAction(commentId int64) error {
	// redis操作：先查redis，若有则更新redis. 若不在redis中，直接走数据库删除，返回客户端。
	commentIdToStr := strconv.FormatInt(commentId, 10)
	n, err := redis.RdbCVid.Exists(redis.Ctx, commentIdToStr).Result()
	if err != nil {
		log.Println(err)
	}
	// 删除评论的消息队列
	commentDelMQ := rabbitmq.SimpleCommentDelMQ
	// 缓存有此id
	if n > 0 {
		// 根据commentId查出对应的videoId
		vid, err := redis.RdbCVid.Get(redis.Ctx, commentIdToStr).Result()
		if err != nil {
			log.Println("redisCV not found", err)
		}
		// 删除缓存,CV直接把key删除，VC只需要移除下面的value(commentId)。
		n1, err := redis.RdbCVid.Del(redis.Ctx, commentIdToStr).Result()
		if err != nil {
			log.Println("redisCV delete failed", err)
		}
		n2, err := redis.RdbCIdComment.Del(redis.Ctx, commentIdToStr).Result()
		if err != nil {
			log.Println("redisCIdComment delete failed", err)
		}
		n3, err := redis.RdbVCid.SRem(redis.Ctx, vid, commentIdToStr).Result()
		if err != nil {
			log.Println("redisVc Remove failed", err)
		}
		log.Println("del comment in redis successfully:", n1, n2, n3)
		// 写数据库
		err = commentDelMQ.PublishSimple(commentIdToStr)
		return err
	}
	// 不在缓存中，直接走数据库，在消息队列中执行
	err = commentDelMQ.PublishSimple(commentIdToStr)
	return err
}

func (commentService *CommentServiceImpl) GetCommentList(videoId int64, userId int64) ([]Comment, error) {
	//redis操作：先查缓存是否命中，若命中，取缓存中的之；否则去读数据库并更新缓存
	videoIdToStr := strconv.FormatInt(videoId, 10)
	cnt, err := redis.RdbVCid.SCard(redis.Ctx, videoIdToStr).Result()
	if err != nil {
		log.Println("SCard failed", err)
	}
	// 缓存中存在评论列表
	if cnt > 0 {
		var commentInfoList []Comment
		log.Println("videoId", videoIdToStr)
		commentIdStringList, err := redis.RdbVCid.SMembers(redis.Ctx, videoIdToStr).Result()
		if err != nil {
			log.Println("read redis vId failed", err)
			//return nil, err
		}
		for _, commentIdString := range commentIdStringList {
			var commentData Comment
			commentString, err := redis.RdbCIdComment.Get(redis.Ctx, commentIdString).Result()
			b := []byte(commentString)
			err = json.Unmarshal(b, &commentData)
			if err != nil {
				log.Println("unmarshal failed", err)
			}
			commentInfoList = append(commentInfoList, commentData)
		}
		log.Println("从redis读取的评论列表")
		sort.Sort(CommentSlice(commentInfoList))
		return commentInfoList, nil
	}
	// 评论不在缓存中，评论既不在缓存也不在数据库中
	// 先根据videoId查评论id，再查用户信息
	plainCommentList, err := dao.GetCommentList(videoId)
	// 拿评论出错
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	n := len(plainCommentList)
	//fmt.Println("视频评论的数量：", n)
	// 如果没有评论, 即评论不在数据库也不在缓存
	if n == 0 {
		return nil, nil
	}
	// 评论在数据库不在缓存中，查数据库并更新到缓存
	commentInfoList := make([]Comment, 0, n)
	var wg sync.WaitGroup
	wg.Add(n)
	for _, comment := range plainCommentList {
		var commentData Comment
		go func(comment dao.Comment) {
			commentService.CombineComment(&commentData, &comment)
			commentInfoList = append(commentInfoList, commentData)
			commentIdToStr := strconv.FormatInt(comment.Id, 10)
			insertRedisVCId(videoIdToStr, commentIdToStr, commentData)
			wg.Done()
		}(comment)

	}
	wg.Wait()
	// 按照评论的先后时间降序排列
	sort.Sort(CommentSlice(commentInfoList))
	// 防止脏读，-1在这里我体会不到用处，反而会产生脏数据
	//_, err = redis.RdbVCid.SAdd(redis.Ctx, videoIdToStr, -1).Result()
	//if err != nil {
	//	log.Println("redis save fail:vId-cId")
	//}
	//// 设置key的过期时间
	//_, err = redis.RdbVCid.Expire(redis.Ctx, videoIdToStr, time.Minute*1).Result()
	//if err != nil {
	//	log.Println("set expire time failed")
	//}
	// 组装好commentList每一个comment序列化后将其存入redis里
	//for _, comment := range commentInfoList {
	//	commentIdToStr := strconv.FormatInt(comment.Id, 10)
	//	insertRedisVCId(videoIdToStr, commentIdToStr, comment)
	//
	//}
	log.Println("get commentList success")
	return commentInfoList, nil
}

func (commentService *CommentServiceImpl) CombineComment(comment *Comment, plainComment *dao.Comment) error {
	commentServiceNew := GetCommentServiceInstance()
	user, err := commentServiceNew.GetUserLoginInfoById(plainComment.UserId)
	if err == nil {
		comment.User = user
	}
	// 随机数生成种子
	rand.Seed(time.Now().Unix())
	comment.Id = plainComment.Id
	comment.Content = plainComment.Content
	comment.CreateDate = plainComment.CreatedAt.Format(config.GO_STARTER_TIME)
	comment.LikeCount = int64(rand.Intn(100))
	comment.TeaseCount = int64(rand.Intn(100))
	return nil
}

func (commentService *CommentServiceImpl) GetCommentCnt(videoId int64) (int64, error) {
	videoIdToStr := strconv.FormatInt(videoId, 10)
	cnt, err := redis.RdbVCid.SCard(redis.Ctx, videoIdToStr).Result()
	if err != nil {
		log.Println("SCard failed", err)
	}
	// 如果在缓存中直接返回
	if cnt > 0 {
		log.Println("从redis读取的评论数量")
		return cnt, nil
	}
	return dao.GetCommentCnt(videoId)
}

// redis中存储videId与commentId对应关系
func insertRedisVCId(videoId string, commentId string, comment Comment) {
	_, err := redis.RdbVCid.SAdd(redis.Ctx, videoId, commentId).Result()
	if err != nil {
		log.Println("redis save fail:vId-cId")
		redis.RdbVCid.Del(redis.Ctx, videoId)
		return
	}
	// 设置键的有效期，为数据不一致情况兜底
	redis.RdbVCid.Expire(redis.Ctx, videoId, config.ExpireTime)
	// 设置键的有效期，为数据不一致情况兜底
	_, err = redis.RdbCVid.Set(redis.Ctx, commentId, videoId, config.ExpireTime).Result()
	if err != nil {
		log.Println("redis save fail:cId-vId")
		return
	}
	b, err := json.Marshal(comment)
	if err != nil {
		log.Println("serialize failed in redis save", err)
	}
	// 设置键的有效期，为数据不一致情况兜底
	_, err = redis.RdbCIdComment.Set(redis.Ctx, commentId, string(b), config.ExpireTime).Result()
	if err != nil {
		log.Println("redis save fail:cId-comment")
		return
	}
}

// CommentSlice Golang实现任意类型sort函数的流程
type CommentSlice []Comment

func (commentSlice CommentSlice) Len() int {
	return len(commentSlice)
}

func (commentSlice CommentSlice) Less(i, j int) bool {
	return commentSlice[i].CreateDate > commentSlice[j].CreateDate
}

func (commentSlice CommentSlice) Swap(i, j int) {
	commentSlice[i], commentSlice[j] = commentSlice[j], commentSlice[i]
}
