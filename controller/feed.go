package controller

import (
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"strconv"
	"time"
	"x-tiktok/service"
)

type FeedResponse struct {
	Response
	VideoList []service.Video `json:"video_list"`
	NextTime  int64           `json:"next_time"`
}

// Feed 不限制登录状态，返回按投稿时间倒序的视频列表，视频数由服务端控制，单次最多30个
func Feed(c *gin.Context) {
	latestTime := c.Query("latest_time")
	//log.Println("返回视频的最新投稿时间戳:", latestTime)
	var convTime time.Time
	if latestTime != "0" {
		t, _ := strconv.ParseInt(latestTime, 10, 64)
		if t > math.MaxInt32 {
			convTime = time.Now()
		} else {
			convTime = time.Unix(t, 0)
		}
	} else {
		convTime = time.Now()
	}
	//log.Println("返回视频的最新投稿时间:", convTime)
	// 获取登录用户的 id（等待用户模块存入用户id到context）
	//userId, _ := strconv.ParseInt(c.GetString("userId"), 10, 64)
	userId := c.GetInt64("userId")
	videoService := service.GetVideoServiceInstance()
	videos, nextTime, err := videoService.Feed(convTime, userId)
	if err != nil {
		c.JSON(http.StatusOK, FeedResponse{
			Response:  Response{StatusCode: 1, StatusMsg: "刷新视频流失败"},
			VideoList: nil,
			NextTime:  nextTime.Unix(),
		})
		return
	}
	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0, StatusMsg: "刷新视频流成功!"},
		VideoList: videos,
		NextTime:  nextTime.Unix(),
	})
}
