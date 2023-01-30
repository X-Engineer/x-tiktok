package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"x-tiktok/service"
)

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	token := c.PostForm("token")
	log.Println("token:", token)
	//if _, exist := usersLoginInfo[token]; !exist {
	//	c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	//	return
	//}

	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	title := c.PostForm("title")
	log.Printf("视频 title: %v\n", title)
	videoService := service.GetVideoServiceInstance()
	// 从 token 中获取 userId
	err = videoService.Publish(data, title, 1)
	if err != nil {
		log.Println("上传文件失败")
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  fmt.Sprintf("《%s》视频上传成功", title),
	})
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: DemoVideos,
	})
}
