package cache

import (
	"github.com/go-redis/redis/v7"
	"learning/conf"
)

var RedisClient *redis.Client

func SetUp() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     conf.AppConfig.Redis.Addr,
		Password: conf.AppConfig.Redis.Password,
		DB:       0, // use default DB
	})
}

const (
	CaptchaPrefix = "captcha"
)
