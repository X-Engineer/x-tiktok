package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

var Ctx = context.Background()

var RdbTest *redis.Client

// InitRedis 初始化 Redis 连接，redis 默认 16 个 DB
func InitRedis() {
	RdbTest = redis.NewClient(&redis.Options{
		Addr:     "ip:port",
		Password: "redis-passwd",
		DB:       0,
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
