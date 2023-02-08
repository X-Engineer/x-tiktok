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
	// 耗时0.09s
	for _, videoId := range []int64{15, 16, 17, 18, 19} {
		video, _ := GetVideoByVideoId(videoId)
		log.Println(video)
	}
}

// 耗时0.02s
func TestGetVideoListById(t *testing.T) {
	videoList, err := GetVideoListById([]int64{15, 16, 17, 18, 19})
	if err == nil {
		log.Println(len(videoList))
		//log.Println(videoList)
	}
	for _, video := range videoList {
		log.Println(video)
	}
}

func TestGetVideoCnt(t *testing.T) {
	count, err := GetVideoCnt(5)
	if err == nil {
		log.Println(count)
	}
}
