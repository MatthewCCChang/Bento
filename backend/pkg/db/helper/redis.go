package helper

import (
	"fmt"
	"context"

	"github.com/redis/go-redis/v9"
)

func UpdateRedis(ctx context.Context, rdb *redis.Client, key string, value interface{}) (error){
	rdb.Set(ctx, key, value, redis.KeepTTL)
	fmt.Println("Successful update for Redis")
	return nil
}