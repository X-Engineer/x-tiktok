package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"x-tiktok/service"
)

// RelationActionResp 关注和取消关注需要返回结构。
type RelationActionResp struct {
	Response
}

type UserListResponse struct {
	Response
	UserList []service.User `json:"user_list"`
}

// RelationAction no practical effect, just check if token is valid
func RelationAction(c *gin.Context) {
	//token := c.Query("token")
	//todo token的鉴权还没有做
	//if _, exist := usersLoginInfo[token]; exist {
	userId, err1 := strconv.ParseInt(c.Query("token"), 10, 64)
	toUserId, err2 := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	actionType, err3 := strconv.ParseInt(c.Query("action_type"), 10, 64)
	//fmt.Sprintf("%d %d %d", userId, toUserId, actionType)
	fmt.Println(userId)
	fmt.Println(toUserId)
	fmt.Println(actionType)
	// 传入参数格式有问题。
	if nil != err1 || nil != err2 || nil != err3 || actionType < 1 || actionType > 2 {
		fmt.Printf("fail")
		c.JSON(http.StatusOK, RelationActionResp{
			Response{
				StatusCode: -1,
				StatusMsg:  "请求参数格式错误",
			},
		})
		return
	}
	// 正常处理
	fsi := service.NewFSIInstance()
	switch {
	// 关注
	case 1 == actionType:
		go fsi.FollowAction(userId, toUserId)
	// 取关
	case 2 == actionType:
		go fsi.CancelFollowAction(userId, toUserId)
	}
	c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "操作成功"})
	//} else {
	//	c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	//}
}

// FollowList all users have same follow list
func FollowList(c *gin.Context) {
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)

	fmt.Println(userId)

	if err != nil {
		fmt.Printf("fail")
		c.JSON(http.StatusOK, UserListResponse{
			Response{
				StatusCode: -1,
				StatusMsg:  "请求参数格式错误",
			},
			nil,
		})
		return
	}

	fsi := service.NewFSIInstance()
	followings, err1 := fsi.GetFollowings(userId)
	if err1 != nil {
		fmt.Printf("fail")
		c.JSON(http.StatusOK, UserListResponse{
			Response{
				StatusCode: -1,
				StatusMsg:  "获取关注列表失败",
			},
			nil,
		})
		return
	}

	c.JSON(http.StatusOK, UserListResponse{
		Response{
			StatusCode: 0,
			StatusMsg:  "获取关注列表成功",
		},
		followings,
	})
}

// FollowerList all users have same follower list
func FollowerList(c *gin.Context) {
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)

	fmt.Println(userId)

	if err != nil {
		fmt.Printf("fail")
		c.JSON(http.StatusOK, UserListResponse{
			Response{
				StatusCode: -1,
				StatusMsg:  "请求参数格式错误",
			},
			nil,
		})
		return
	}

	fsi := service.NewFSIInstance()
	followers, err1 := fsi.GetFollowers(userId)
	if err1 != nil {
		fmt.Printf("fail")
		c.JSON(http.StatusOK, UserListResponse{
			Response{
				StatusCode: -1,
				StatusMsg:  "获取粉丝列表失败",
			},
			nil,
		})
		return
	}

	c.JSON(http.StatusOK, UserListResponse{
		Response{
			StatusCode: 0,
			StatusMsg:  "获取粉丝列表成功",
		},
		followers,
	})
}

// FriendList all users have same friend list
func FriendList(c *gin.Context) {
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: []service.User{},
	})
}
