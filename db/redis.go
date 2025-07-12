package db

import (
	"context"

	"github.com/manthan307/corebase/utils/helper"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func NewRedisClient(log *zap.Logger) *redis.Client {
	redisURL := helper.GetEnv("REDIS_URL", "")
	if redisURL == "" {
		panic("REDIS_URL is not set")
	}

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		panic(err)
	}

	redisClient := redis.NewClient(opt)

	err = redisClient.Ping(context.Background()).Err()
	if err != nil {
		panic(err)
	}

	log.Info("✅ Connected to Redis")

	return redisClient
}
