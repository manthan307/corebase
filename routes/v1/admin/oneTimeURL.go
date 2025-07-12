package admin

import (
	"context"
	"fmt"
	"time"

	"github.com/manthan307/corebase/db"
	"github.com/manthan307/corebase/utils/helper"
	"github.com/redis/go-redis/v9"
)

func OneTimeURL(ctx context.Context, client *db.Client, role string, exp time.Duration) (string, error) {
	token := helper.GenerateRandomString(64)
	key := fmt.Sprintf("setup_token:%s", token)

	if err := client.RedisClient.Set(ctx, key, role, exp).Err(); err != nil {
		return "", err
	}

	port := helper.GetEnv("PORT", "8000")
	baseURL := helper.GetEnv("HOST", fmt.Sprintf("http://localhost:%s", port))
	url := fmt.Sprintf("%s/_/setup?token=%s", baseURL, token)

	return url, nil
}

func ValidateToken(ctx context.Context, token string, client *db.Client) (string, error) {
	key := fmt.Sprintf("setup_token:%s", token)
	role, err := client.RedisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("token not found or expired")
	} else if err != nil {
		return "", err
	}

	return role, nil
}

func DeleteToken(ctx context.Context, token string, client *db.Client) error {
	key := fmt.Sprintf("setup_token:%s", token)
	if err := client.RedisClient.Del(ctx, key).Err(); err != nil {
		return err
	}
	return nil
}

func ConsumeToken(ctx context.Context, token string, client *db.Client) (string, error) {
	role, err := ValidateToken(ctx, token, client)
	if err != nil {
		return "", err
	}
	if err := DeleteToken(ctx, token, client); err != nil {
		return "", err
	}
	return role, nil
}
