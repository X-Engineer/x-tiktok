package config

import "time"

// 存放相关配置

const GO_STARTER_TIME = "2006-01-02 15:04:05"

// 视频模块相关配置
const (
	VIDEO_NUM_PER_REFRESH     = 6
	VIDEO_INIT_NUM_PER_AUTHOR = 10
	// 阿里 OSS 相关配置
	OSS_ACCESS_KEY_ID     = "LTAI5tGj5YiM8KGsUdqEy978"
	OSS_ACCESS_KEY_SECRET = "WRhFkA5KDnynjuAad3FqWHTnCpPDDw"
	OSS_BUCKET_NAME       = "xlab-open-source"
	OSS_ENDPOINT          = "http://oss-cn-beijing.aliyuncs.com"
	CUSTOM_DOMAIN         = "https://oss.x-lab.info/"
	OSS_VIDEO_DIR         = "zhicheng-ning/bytedance-go/x-tiktok/videos/"
	OSS_USER_AVATAR_DIR   = "zhicheng-ning/bytedance-go/x-tiktok/users/avatar.jpg"
	PLAY_URL_PREFIX       = CUSTOM_DOMAIN + OSS_VIDEO_DIR
	COVER_URL_SUFFIX      = "?x-oss-process=video/snapshot,t_2000,m_fast"
)

// jwt密钥
var SECRETE = "x-engineer"

const LIKE = 1

var LatestRequestTime = make(map[string]time.Time, 100)

var ExpireTime = time.Hour * 24
