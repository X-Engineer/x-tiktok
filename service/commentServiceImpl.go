package service

import (
	"log"
	"math/rand"
	"sort"
	"strconv"
	"sync"
	"time"
	"x-tiktok/config"
	"x-tiktok/dao"
	"x-tiktok/middleware/rabbitmq"
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
		LikeCount:  int64(rand.Intn(10000)),
		TeaseCount: int64(rand.Intn(100)),
	}
	return commentData, nil
}

func (commentService *CommentServiceImpl) DeleteCommentAction(commentId int64) error {
	commentDelMQ := rabbitmq.SimpleCommentDelMQ
	err := commentDelMQ.PublishSimple(strconv.FormatInt(commentId, 10))
	return err
	//return dao.DeleteComment(commentId)
}

func (commentService *CommentServiceImpl) GetCommentList(videoId int64, userId int64) ([]Comment, error) {
	// 先根据videoId查评论id，再查用户信息
	plainCommentList, err := dao.GetCommentList(videoId)
	// 拿评论出错
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	n := len(plainCommentList)
	//fmt.Println("视频评论的数量：", n)
	commentInfoList := make([]Comment, 0, n)
	var wg sync.WaitGroup
	wg.Add(n)
	for _, comment := range plainCommentList {
		var commentData Comment
		go func(comment dao.Comment) {
			commentService.CombineComment(&commentData, &comment)
			commentInfoList = append(commentInfoList, commentData)
			wg.Done()
		}(comment)
	}
	wg.Wait()
	// 按照评论的先后时间降序排列
	sort.Sort(CommentSlice(commentInfoList))
	//fmt.Println(commentInfoList)
	return commentInfoList, nil
}

func (commentService *CommentServiceImpl) CombineComment(comment *Comment, plainComment *dao.Comment) error {
	commentServiceNew := GetCommentServiceInstance()
	user, err := commentServiceNew.GetUserLoginInfoById(plainComment.UserId)
	if err != nil {
		comment.User = user
	}
	comment.Id = plainComment.Id
	comment.Content = plainComment.Content
	comment.CreateDate = plainComment.CreatedAt.Format(config.GO_STARTER_TIME)
	return nil
}

func (commentService *CommentServiceImpl) GetCommentCnt(videoId int64) (int64, error) {
	return dao.GetCommentCnt(videoId)
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
