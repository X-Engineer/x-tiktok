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
	VideoList []service.Video `json:"video_list,omitempty"`
}

// 赞操作
func FavoriteAction(c *gin.Context) {
	user_id, _ := strconv.ParseInt(c.Query("userId"), 10, 64)
	video_id, _ := strconv.ParseInt(c.Query("video_id"), 10, 64)
	action_type, _ := strconv.ParseInt(c.Query("action_type"), 10, 32)

	//用户信息校验
	token := c.Query("token")
	res, _ := util.ParseToken(token)
	if res == nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Token鉴权失败"},
		})
	} else {
		usi := service.UserServiceImpl{}
		//id, _ := strconv.ParseInt(userId, 10, 64)
		if user, err := usi.GetUserLoginInfoById(user_id); err != nil {
			c.JSON(http.StatusOK, UserResponse{
				Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
			})
		} else {
			c.JSON(http.StatusOK, UserResponse{
				Response: Response{StatusCode: 0, StatusMsg: "Query Success"},
				User:     User(user),
			})
		}
	}

	Fni := service.NewLikeServImpInstance()
	err := Fni.FavoriteAction(user_id, video_id, int32(action_type))

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
				StatusCode: 1,
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
	//curId, _ := strconv.ParseInt(strCurId, 10, 64)
	//like := GetVideo()
	//videos, err := like.GetFavouriteList(userId, curId)
	_, err := Fni.GetLikesList(userId)
	//先定义一个假的空数组
	var videolis []service.Video
	if err == nil {
		log.Printf("方法like.GetFavouriteList(userid) 成功")
		c.JSON(http.StatusOK, GetFavouriteListResponse{
			Response: Response{StatusCode: 0, StatusMsg: "get favouriteList success"},
			//调用video接口获取视频具体信息
			VideoList: videolis,
		})
	} else {
		log.Printf("方法like.GetFavouriteList(userid) 失败：%v", err)
		c.JSON(http.StatusOK, GetFavouriteListResponse{
			Response: Response{StatusCode: 1, StatusMsg: "get favouriteList fail "},
		})
	}
}
