package model

import (
	"context"
	"github.com/go-redis/redis/v8"
)

const (
	RedisServer = "tesseract_redis:6379"
)
var (
	ctx = context.Background()
	rdb = redis.NewClient(&redis.Options{
		Addr: RedisServer,
		Password: "",
		DB: 0,
	})
)

func SendToRedis(taskID string, text string) error {
	err := rdb.Set(ctx, taskID, text, 0).Err()
	return err
}

func GetFromRedis(taskID string) (res string, err error) {
	text, err := rdb.Get(ctx, taskID).Result()
	return text, err
}


