package config

import "time"

// 存放相关配置

const GO_STARTER_TIME = "2006-01-02 15:04:05"

// 视频模块相关配置
const (
	VIDEO_NUM_PER_REFRESH     = 6
	VIDEO_INIT_NUM_PER_AUTHOR = 10
	// 阿里 OSS 相关配置
	OSS_ACCESS_KEY_ID     = "OSS_ACCESS_KEY_ID"
	OSS_ACCESS_KEY_SECRET = "OSS_ACCESS_KEY_SECRET"
	OSS_BUCKET_NAME       = "OSS_BUCKET_NAME"
	OSS_ENDPOINT          = "OSS_ENDPOINT"
	CUSTOM_DOMAIN         = "CUSTOM_DOMAIN"
	OSS_VIDEO_DIR         = "OSS_VIDEO_DIR"
	PLAY_URL_PREFIX       = CUSTOM_DOMAIN + OSS_VIDEO_DIR
	COVER_URL_SUFFIX      = "?x-oss-process=video/snapshot,t_2000,m_fast"
)

const LIKE = 1

var LatestRequestTime = make(map[string]time.Time, 100)
