package service

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/google/uuid"
	"log"
	"mime/multipart"
	"sync"
	"time"
	"x-tiktok/config"
	"x-tiktok/dao"
)

type VideoServiceImpl struct {
	UserService
}

var (
	videoServiceImp  *VideoServiceImpl
	videoServiceOnce sync.Once
)

// GetVideoServiceInstance Go 单例模式：https://www.liwenzhou.com/posts/Go/singleton/
func GetVideoServiceInstance() *VideoServiceImpl {
	videoServiceOnce.Do(func() {
		videoServiceImp = &VideoServiceImpl{
			UserService: &UserServiceImpl{},
		}
	})
	return videoServiceImp
}

func (videoService *VideoServiceImpl) Publish(data *multipart.FileHeader, title string, userId int64) error {
	// 保证唯一的 videoName
	videoName := uuid.New().String()
	err := UploadVideoToOSS(data, videoName)
	if err != nil {
		return err
	}
	err = dao.UploadVideo(videoName, userId, title)
	if err != nil {
		log.Println("视频存入数据库失败！")
		return err
	}
	return nil
}

func (videoService *VideoServiceImpl) Feed(latestTime time.Time, userId int64) ([]Video, time.Time, error) {
	videos := make([]Video, 0, config.VIDEO_NUM_PER_REFRESH)
	plainVideos, err := dao.GetVideosByLatestTime(latestTime)
	if err != nil {
		log.Println("GetVideosByLatestTime:", err)
		return nil, time.Time{}, err
	}
	err = videoService.getRespVideos(&videos, &plainVideos, userId)
	if err != nil {
		log.Println("getRespVideos:", err)
		return nil, time.Time{}, err
	}
	return videos, plainVideos[len(plainVideos)-1].CreatedAt, nil
}

func (videoService *VideoServiceImpl) PublishList(userId int64) ([]Video, error) {
	videos := make([]Video, 0, config.VIDEO_INIT_NUM_PER_AUTHOR)
	plainVideos, err := dao.GetVideosByUserId(userId)
	if err != nil {
		log.Println("GetVideosByUserId:", err)
		return nil, err
	}
	err = videoService.getRespVideos(&videos, &plainVideos, userId)
	if err != nil {
		log.Println("getRespVideos:", err)
		return nil, err
	}
	return videos, nil
}

func UploadVideoToOSS(file *multipart.FileHeader, videoName string) error {
	// 创建OSSClient实例。
	// yourEndpoint填写Bucket对应的Endpoint，以华东1（杭州）为例，填写为https://oss-cn-hangzhou.aliyuncs.com。其它Region请按实际情况填写。
	// 阿里云 RAM 用户信息
	client, err := oss.New(config.OSS_ENDPOINT, config.OSS_ACCESS_KEY_ID, config.OSS_ACCESS_KEY_SECRET)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	// 填写存储空间名称，例如examplebucket。
	bucket, err := client.Bucket(config.OSS_BUCKET_NAME)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	fd, err := file.Open()
	if err != nil {
		log.Println("file open failed!")
		return err
	}
	defer fd.Close()

	err = bucket.PutObject(config.OSS_VIDEO_DIR+videoName+".mp4", fd)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	return nil
}

func (videoService *VideoServiceImpl) getRespVideos(videos *[]Video, plainVideos *[]dao.Video, userId int64) error {
	for _, tmpVideo := range *plainVideos {
		var video Video
		videoService.combineVideo(&video, &tmpVideo, userId)
		*videos = append(*videos, video)
	}
	return nil
}

// 组装 controller 层所需的 Video 结构
func (videoService *VideoServiceImpl) combineVideo(video *Video, plainVideo *dao.Video, userId int64) error {
	// 建立协程组，确保所有协程的任务完成后才退出本方法
	var wg sync.WaitGroup
	wg.Add(4)
	video.Video = *plainVideo
	// 视频作者信息
	go func(v *Video) {
		user, err := videoService.GetUserLoginInfoById(video.AuthorId)
		if err != nil {
			return
		}
		v.Author = user
		wg.Done()
	}(video)

	// 视频点赞数量
	go func(v *Video) {
		// 等待点赞服务，获取视频点赞量
		v.FavoriteCount = 10
		wg.Done()
	}(video)

	// 视频评论数量
	go func(v *Video) {
		// 等待评论服务，获取视频评论量
		v.CommentCount = 10
		wg.Done()
	}(video)

	// 当前登录用户/游客是否对该视频点过赞
	go func(v *Video) {
		// 等待点赞服务，获取是否点赞
		v.IsFavorite = false
		wg.Done()
	}(video)

	wg.Wait()
	return nil
}
