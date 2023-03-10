package redis

import (
	"context"
	"log"
	"testing"
)

func Test_Redis(t *testing.T) {
	NewRedis()
	ctx := context.Background()
	result, err := RedisClient.Set(ctx, "player1", 100, 0).Result()
	if len(result) > 0 && err == nil {
		log.Println("test redis success")
	}
}
