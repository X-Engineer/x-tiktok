package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

var Ctx = context.Background()
var NilError = redis.Nil
var RdbTest *redis.Client

// UserFollowings 根据用户id找到他关注的人
var UserFollowings *redis.Client

// UserFollowers 根据用户id找到他的粉丝
var UserFollowers *redis.Client

// UserFriends 根据用户id找到他的好友
var UserFriends *redis.Client

// RdbVCid 存储video与comment的关系
var RdbVCid *redis.Client

// RdbCVid 根据commentId找videoId
var RdbCVid *redis.Client

// RdbCIdComment 根据commentId 找comment
var RdbCIdComment *redis.Client

// RdbUVid 根据userId找到他点赞过的videoId
var RdbUVid *redis.Client

// RdbVUid 根据videoId找到点赞过它的userId
var RdbVUid *redis.Client

const (
	ProdRedisAddr = "ip:port"
	ProRedisPwd   = "redis-passwd"
)

// InitRedis 初始化 Redis 连接，redis 默认 16 个 DB
func InitRedis() {
	RdbTest = redis.NewClient(&redis.Options{
		Addr:     ProdRedisAddr,
		Password: ProRedisPwd,
		DB:       0,
	})
	RdbVCid = redis.NewClient(&redis.Options{
		Addr:     ProdRedisAddr,
		Password: ProRedisPwd,
		DB:       1,
	})
	RdbCVid = redis.NewClient(&redis.Options{
		Addr:     ProdRedisAddr,
		Password: ProRedisPwd,
		DB:       2,
	})
	RdbCIdComment = redis.NewClient(&redis.Options{
		Addr:     ProdRedisAddr,
		Password: ProRedisPwd,
		DB:       3,
	})
	RdbUVid = redis.NewClient(&redis.Options{
		Addr:     ProdRedisAddr,
		Password: ProRedisPwd,
		DB:       4,
	})
	RdbVUid = redis.NewClient(&redis.Options{
		Addr:     ProdRedisAddr,
		Password: ProRedisPwd,
		DB:       5,
	})
	UserFollowings = redis.NewClient(&redis.Options{
		Addr:     ProdRedisAddr,
		Password: ProRedisPwd,
		DB:       11,
	})
	UserFollowers = redis.NewClient(&redis.Options{
		Addr:     ProdRedisAddr,
		Password: ProRedisPwd,
		DB:       12,
	})
	UserFriends = redis.NewClient(&redis.Options{
		Addr:     ProdRedisAddr,
		Password: ProRedisPwd,
		DB:       13,
	})
}

// 测试连接 Redis
func connRedis() {
	InitRedis()
	_, err := RdbTest.Ping(Ctx).Result()
	if err != nil {
		log.Panicf("连接 redis 错误，错误信息: %v", err)
	} else {
		log.Println("Redis 连接成功！")
	}
}

// Go 操作 Redis
// 更多命令参考：https://www.cnblogs.com/itbsl/p/14198111.html
func setValue(key string, value interface{}) {
	InitRedis()
	// 设置 2 min 过期，如果 expiration 为 0 表示永不过期
	RdbTest.Set(Ctx, key, value, 2*time.Minute)
}

// 测试获取值
func getValue(key string) {
	InitRedis()
	val, err := RdbTest.Get(Ctx, key).Result()
	switch {
	case err == redis.Nil:
		log.Println("key does not exist")
	case err != nil:
		log.Println("Get failed", err)
	case val == "":
		log.Println("value is empty")
	case val != "":
		log.Println("value is", val)
	}
}
