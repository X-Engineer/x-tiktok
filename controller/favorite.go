package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"x-tiktok/service"
	"x-tiktok/util"
)

type FavoriteActionResponse struct {
	Response
}

type GetFavouriteListResponse struct {
	Response
	VideoList []service.Video `json:"video_list"`
}

// 赞操作
func FavoriteAction(c *gin.Context) {
	video_id, _ := strconv.ParseInt(c.Query("video_id"), 10, 64)
	action_type, _ := strconv.ParseInt(c.Query("action_type"), 10, 32)
	token := c.Query("token")
	//获取用户信息
	claims, _ := util.ParseToken(token)
	c.Set("userId", claims.ID)

	Fni := service.NewLikeServImpInstance()
	err := Fni.FavoriteAction(c.GetInt64("userId"), video_id, int32(action_type))
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

// 获取点赞列表
func FavoriteList(c *gin.Context) {
	strUserId := c.Query("user_id")
	//likeCnt:=dao.VideoLikedCount()
	userId, _ := strconv.ParseInt(strUserId, 10, 64)
	Fni := service.NewLikeServImpInstance()

	_, err := Fni.GetLikesList(userId)
	//video类型数组的假数据
	var video_list []service.Video
	video_list[0].Id = 1
	video_list[0].IsFavorite = true
	video_list[0].AuthorId = 1
	video_list[0].CoverUrl = ""
	video_list[0].PlayUrl = ""
	video_list[0].Title = ""
	video_list[0].CommentCount = 10
	video_list[0].FavoriteCount = 100
	video_list[0].IsFavorite = true
	if err == nil {
		log.Printf("方法like.GetFavouriteList(userid) 成功")
		c.JSON(http.StatusOK, GetFavouriteListResponse{
			Response: Response{StatusCode: 0, StatusMsg: "get favouriteList success"}, VideoList: video_list,
		})
	} else {
		log.Printf("方法like.GetFavouriteList(userid) 失败：%v", err)
		c.JSON(http.StatusOK, Response{StatusCode: -1, StatusMsg: "get favouriteList fail "})
	}
}
