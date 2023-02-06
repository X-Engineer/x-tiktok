package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"x-tiktok/service"
)

type FavoriteActionResponse struct {
	Response
}

type GetFavouriteListResponse struct {
	Response
	VideoList []service.Video `json:"video_list"`
}

// FavoriteAction 点赞操作
func FavoriteAction(c *gin.Context) {
	videoId, _ := strconv.ParseInt(c.Query("video_id"), 10, 64)
	actionType, _ := strconv.ParseInt(c.Query("action_type"), 10, 32)
	Fni := service.NewLikeServImpInstance()
	err := Fni.FavoriteAction(c.GetInt64("userId"), videoId, int32(actionType))
	if err == nil {
		log.Printf("方法like.FavouriteAction(userid, videoId, int32(actiontype) 成功")
		c.JSON(http.StatusOK, FavoriteActionResponse{
			Response{
				StatusCode: 0,
				StatusMsg:  "favourite action success",
			},
		})
	} else {
		log.Printf("方法like.FavouriteAction(userid, videoId, int32(actiontype) 失败：%v", err)
		c.JSON(http.StatusOK, FavoriteActionResponse{
			Response{
				StatusCode: -1,
				StatusMsg:  "favourite action fail",
			},
		})
	}
}

// FavoriteList 获取用户的点赞视频列表
func FavoriteList(c *gin.Context) {
	strUserId := c.Query("user_id")
	//likeCnt:=dao.VideoLikedCount()
	userId, _ := strconv.ParseInt(strUserId, 10, 64)
	Fni := service.NewLikeServImpInstance()

	_, err := Fni.GetLikesList(userId)
	//video类型数组的假数据
	//var video_list []service.Video
	videoList := make([]service.Video, 1, 2)
	videoList[0].Id = 1
	videoList[0].IsFavorite = true
	videoList[0].AuthorId = 1
	videoList[0].CoverUrl = ""
	videoList[0].PlayUrl = ""
	videoList[0].Title = ""
	videoList[0].CommentCount = 10
	videoList[0].FavoriteCount = 100
	videoList[0].IsFavorite = true
	if err == nil {
		log.Printf("方法like.GetFavouriteList(userid) 成功")
		c.JSON(http.StatusOK, GetFavouriteListResponse{
			Response: Response{StatusCode: 0, StatusMsg: "get favouriteList success"}, VideoList: videoList,
		})
	} else {
		log.Printf("方法like.GetFavouriteList(userid) 失败：%v", err)
		c.JSON(http.StatusOK, Response{StatusCode: -1, StatusMsg: "get favouriteList fail "})
	}
}
