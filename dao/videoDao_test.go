package dao

import (
	"log"
	"testing"
	"time"
	"x-tiktok/config"
)

func TestUploadVideo(t *testing.T) {
	UploadVideo("VID_2023_1_29", 1, "测试视频1")
}

func TestGetVideosByUserId(t *testing.T) {
	res, err := GetVideosByUserId(1)
	if err == nil {
		for _, re := range res {
			log.Println(re)
		}
	}
}

func TestGetVideosByLatestTime(t *testing.T) {
	// 时区修正
	mockTime, _ := time.ParseInLocation(config.GO_STARTER_TIME, "2023-01-29 21:20:04", time.Local)
	log.Println(mockTime)
	res, err := GetVideosByLatestTime(mockTime)
	if err == nil {
		for _, re := range res {
			log.Println(re)
		}
	}
}

func TestGetVideoByVideoId(t *testing.T) {
	video, err := GetVideoByVideoId(1)
	if err == nil {
		log.Println(video)
	}
}
