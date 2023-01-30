package service

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/google/uuid"
	"log"
	"mime/multipart"
	"sync"
	"x-tiktok/config"
	"x-tiktok/dao"
)

type VideoServiceImpl struct {
}

var (
	videoServiceImp  *VideoServiceImpl
	videoServiceOnce sync.Once
)

// GetVideoServiceInstance Go 单例模式：https://www.liwenzhou.com/posts/Go/singleton/
func GetVideoServiceInstance() *VideoServiceImpl {
	videoServiceOnce.Do(func() {
		videoServiceImp = &VideoServiceImpl{}
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
