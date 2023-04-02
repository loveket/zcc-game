package game

import (
	"context"
	"xiuianserver/redis"
)

type Friend struct {
	ctx         context.Context
	ApplyPlayer []string
	Friends     []string
}

func NewFriend() *Friend {
	return &Friend{
		ctx:     context.Background(),
		Friends: make([]string, 100),
	}
}

func (f *Friend) SearchFriend(name string) bool {
	if ok, err := redis.RedisClient.HExists(f.ctx, "user", name).Result(); ok && err == nil {
		return true
	}
	return false
}
