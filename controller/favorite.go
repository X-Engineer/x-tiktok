package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"x-tiktok/dao"
	"x-tiktok/service"
)

type FavoriteActionResponse struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type GetFavouriteListResponse struct {
	StatusCode int32           `json:"status_code"`
	StatusMsg  string          `json:"status_msg,omitempty"`
	VideoList  []service.Video `json:"video_list,omitempty"`
}

// FavoriteAction no practical effect, just check if token is valid
func FavoriteAction(c *gin.Context) {
	//token := c.Query("token")
	user_id, _ := strconv.ParseInt(c.Query("userId"), 10, 64)
	video_id, _ := strconv.ParseInt(c.Query("video_id"), 10, 64)
	action_type, _ := strconv.ParseInt(c.Query("action_type"), 10, 8)
	like := new(service.LikeServiceImpl)

	//用户信息校验
	token := c.Query("token")
	res, _ := util.ParseToken(token)
	if res == nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Token鉴权失败"},
		})
	} else {
		usi := service.UserServiceImpl{}
		id, _ := strconv.ParseInt(userId, 10, 64)
		if user, err := usi.GetUserLoginInfoById(id); err != nil {
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
	//if _, exist := usersLoginInfo[token]; exist {
	//	c.JSON(http.StatusOK, Response{StatusCode: 0})
	//} else {
	//	c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	//}

	err := like.FavoriteAction(user_id, video_id, int8(action_type))
	if err == nil {
		log.Printf("方法like.FavouriteAction(userid, videoId, int32(actiontype) 成功")
		c.JSON(http.StatusOK, FavoriteActionResponse{
			StatusCode: 0,
			StatusMsg:  "favourite action success",
		})
	} else {
		log.Printf("方法like.FavouriteAction(userid, videoId, int32(actiontype) 失败：%v", err)
		c.JSON(http.StatusOK, FavoriteActionResponse{
			StatusCode: 1,
			StatusMsg:  "favourite action fail",
		})
	}
}

func FavoriteList(c *gin.Context) {
	strUserId := c.Query("user_id")
	//likeCnt:=dao.VideoLikedCount()
	userId, _ := strconv.ParseInt(strUserId, 10, 64)
	curId, _ := strconv.ParseInt(strCurId, 10, 64)
	//like := GetVideo()
	//videos, err := like.GetFavouriteList(userId, curId)
	videos, err := dao.GetLikeListByUserId(userId)
	if err == nil {
		log.Printf("方法like.GetFavouriteList(userid) 成功")
		c.JSON(http.StatusOK, GetFavouriteListResponse{
			StatusCode: 0,
<<<<<<< Updated upstream
		},
		VideoList: nil,
	})
=======
			StatusMsg:  "get favouriteList success",
			VideoList:  videos,
		})
	} else {
		log.Printf("方法like.GetFavouriteList(userid) 失败：%v", err)
		c.JSON(http.StatusOK, GetFavouriteListResponse{
			StatusCode: 1,
			StatusMsg:  "get favouriteList fail ",
		})
	}
>>>>>>> Stashed changes
}
