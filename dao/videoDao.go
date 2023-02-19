package dao

import (
	"log"
	"time"
	"x-tiktok/config"
)

type Video struct {
	Id        int64 `json:"id"`
	AuthorId  int64
	Title     string `json:"title"`
	PlayUrl   string `json:"play_url"`
	CoverUrl  string `json:"cover_url"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// TableName 数据库表名映射到 Video 结构体
func (Video) TableName() string {
	return "video"
}

// SaveVideo 保存视频记录到数据库中
func SaveVideo(video Video) error {
	result := Db.Save(&video)
	if result.Error != nil {
		log.Println("数据库保存视频失败！", result.Error)
		return result.Error
	}
	return nil
}

// UploadVideo 上传视频
func UploadVideo(videoName string, authorId int64, videoTitle string) error {
	var video Video
	video.AuthorId = authorId
	video.Title = videoTitle
	video.PlayUrl = config.PLAY_URL_PREFIX + videoName + ".mp4"
	video.CoverUrl = video.PlayUrl + config.COVER_URL_SUFFIX
	video.CreatedAt = time.Now()
	video.UpdatedAt = time.Now()
	return SaveVideo(video)
}

// GetVideosByUserId 根据用户 Id 获取该用户已发布的所有视频
func GetVideosByUserId(userId int64) ([]Video, error) {
	// 预定义容量，避免多次扩容
	videos := make([]Video, 0, config.VIDEO_INIT_NUM_PER_AUTHOR)
	result := Db.Where(&Video{AuthorId: userId}).Find(&videos)
	if result.Error != nil {
		log.Println("获取用户已发布视频失败！")
		return nil, result.Error
	}
	return videos, nil
}

// GetVideosByLatestTime 按投稿时间倒序的视频列表
func GetVideosByLatestTime(latestTime time.Time) ([]Video, error) {
	videos := make([]Video, config.VIDEO_NUM_PER_REFRESH)
	result := Db.Where("created_at < ?", latestTime).
		Order("created_at desc").
		Limit(config.VIDEO_NUM_PER_REFRESH).
		Find(&videos)
	if result.RowsAffected == 0 {
		log.Println("没有更多视频了！")
		return videos, nil
	}
	if result.Error != nil {
		log.Println("获取视频 Feed 失败！")
		return nil, result.Error
	}
	return videos, nil
}

// GetVideoByVideoId 根据视频 Id 获取视频信息
func GetVideoByVideoId(videoId int64) (Video, error) {
	var video Video
	result := Db.Where("id = ?", videoId).First(&video)
	if result.Error != nil {
		log.Println("根据视频 Id 获取视频失败！")
		return video, result.Error
	}
	return video, nil
}

// GetVideoListById 根据videoIdList查询视频信息
func GetVideoListById(videoIdList []int64) ([]Video, error) {
	var videoList []Video
	result := Db.Model(Video{}).
		Where("id in (?)", videoIdList).
		Find(&videoList)
	if result.Error != nil {
		return videoList, result.Error
	}
	return videoList, nil
}

func GetVideoCnt(userId int64) (int64, error) {
	var count int64
	result := Db.Model(Video{}).
		Where("author_id = ?", userId).
		Count(&count)
	if result.Error != nil {
		log.Println("根据userId获取作品数量失败！")
		return -1, nil
	}
	return count, nil
}
