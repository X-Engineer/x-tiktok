package test

import (
	"net/http"
	"testing"
)

func TestLike(t *testing.T) {
	e := newExpect(t)

	userId, token := getTestUserToken(testUserA, e)

	likeActionResp := e.POST("/douyin/favorite/action/").
		WithQuery("token", token).WithQuery("video_id", 1).WithQuery("action_type", 1).
		WithFormField("token", token).WithFormField("video_id", 1).WithFormField("action_type", 1).
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	likeActionResp.Value("status_code").Number().Equal(0)

	likeListResp := e.GET("/douyin/favorite/list/").
		WithQuery("token", token).WithQuery("user_id", userId).
		WithFormField("token", token).WithFormField("user_id", userId).
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	likeListResp.Value("status_code").Number().Equal(0)
	likeListResp.Value("video_list").Array().Length().Gt(0)
}
