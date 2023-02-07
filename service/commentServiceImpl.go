package service

import (
	"log"
	"sync"
	"x-tiktok/dao"
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
	commentData := Comment{
		Id:         commentRes.Id,
		User:       user,
		Content:    commentRes.Content,
		CreateDate: commentRes.CreatedAt.Format("2006-01-02 15:04:05"),
	}
	return commentData, nil
}

func (commentService *CommentServiceImpl) DeleteCommentAction(commentId int64) error {
	return dao.DeleteComment(commentId)
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
	for _, comment := range plainCommentList {
		var commentData Comment
		// 组合评论的用户信息
		commentService.CombineComment(&commentData, &comment)
		commentInfoList = append(commentInfoList, commentData)
	}
	return commentInfoList, nil
}

func (commentService *CommentServiceImpl) CombineComment(comment *Comment, plainComment *dao.Comment) error {
	var wg sync.WaitGroup
	commentServiceNew := GetCommentServiceInstance()
	wg.Add(1)
	user, err := commentServiceNew.GetUserLoginInfoById(plainComment.UserId)
	if err != nil {
		comment.User = user
	}
	comment.Id = plainComment.Id
	comment.Content = plainComment.Content
	comment.CreateDate = plainComment.CreatedAt.Format("2006-01-02 15:04:05")
	wg.Done()
	wg.Wait()
	return nil
}

func (commentService *CommentServiceImpl) GetCommentCnt(videoId int64) (int64, error) {
	return dao.GetCommentCnt(videoId)
}
