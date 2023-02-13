package util

import (
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

// GetQueryParams 获取 GET 请求中的 Query 参数
func GetQueryParams(c *gin.Context) map[string]any {
	query := c.Request.URL.Query()
	var queryMap = make(map[string]any, len(query))
	for k := range query {
		queryMap[k] = c.Query(k)
	}
	return queryMap
}

// GetPostFormParams 获取 POST 请求中的参数
func GetPostFormParams(c *gin.Context) (map[string]any, error) {
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		if !errors.Is(err, http.ErrNotMultipart) {
			return nil, err
		}
	}
	var postMap = make(map[string]any, len(c.Request.PostForm))
	for k, v := range c.Request.PostForm {
		if len(v) > 1 {
			postMap[k] = v
		} else if len(v) == 1 {
			postMap[k] = v[0]
		}
	}

	return postMap, nil
}

// GetBody 获取请求体的数据
func GetBody(c *gin.Context) []byte {
	// 读取body数据
	body, err := c.GetRawData()
	if err != nil {
		return nil
	}
	//把读过的字节流重新放到body
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	return body
}
