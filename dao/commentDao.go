package dao

import (
	"errors"
	"log"
	"time"
)

type Comment struct {
	Id         int64     //评论id
	UserId     int64     //评论用户id
	VideoId    int64     //视频id
	Content    string    //评论内容
	ActionType int64     //发布评论为1，取消评论为2
	CreatedAt  time.Time //评论发布的日期mm-dd
	UpdatedAt  time.Time
}

// TableName 修改表名映射
func (Comment) TableName() string {
	return "comment"
}

func InsertComment(comment Comment) (Comment, error) {
	if err := Db.Model(Comment{}).Create(&comment).Error; err != nil {
		log.Println(err.Error())
		return Comment{}, err
	}
	return comment, nil
}

func DeleteComment(commentId int64) error {
	var comment Comment
	// 先查询是否有此评论
	result := Db.Where("id = ?", commentId).
		First(&comment)
	if result.Error != nil {
		return errors.New("del comment is not exist")
	}
	// 删除评论，将action_type置为2
	result = Db.Model(Comment{}).
		Where("id=?", commentId).
		Update("action_type", 2)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetCommentList(videoId int64) ([]Comment, error) {
	var commentList []Comment
	result := Db.Model(Comment{}).Where(map[string]interface{}{"video_id": videoId, "action_type": 1}).
		Order("created_at desc").
		Find(&commentList)
	if result.Error != nil {
		log.Println(result.Error)
		return commentList, errors.New("get comment list failed")
	}
	return commentList, nil
}

func GetCommentCnt(videoId int64) (int64, error) {
	var count int64
	result := Db.Model(Comment{}).Where(map[string]interface{}{"video_id": videoId, "action_type": 1}).
		Count(&count)
	if result.Error != nil {
		return 0, errors.New("find comments count failed")
	}
	return count, nil
}
