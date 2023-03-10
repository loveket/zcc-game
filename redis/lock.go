package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

var lockflag string = "zcc"
var timeout time.Duration = 5 * time.Second
var stopChan = make(chan struct{})
var unlockScript = `
		if redis.call('get',KEYS[1])==ARGV[1]
		then
			return redis.call('del',KEYS[1])
		else
			return 0
		end
	`
var resetScript = `
	if redis.call('get',KEYS[1])==ARGV[1]
		then
		return redis.call('expire',KEYS[1],ARGV[2])
	else
		return 0
	end
`

func Lock(ctx context.Context, uuid string, client *redis.Client) {
	for {
		locksucess, err := client.SetNX(ctx, lockflag, uuid, timeout).Result()
		if err == nil && locksucess {
			//监听过期时间
			go listenResetTime(ctx, uuid)
			log.Println("lock success")
			return
		}
	}
}
func UnLock(ctx context.Context, uuid string, client *redis.Client) {
	for {
		script := redis.NewScript(unlockScript)
		result, err := script.Run(ctx, client, []string{lockflag}, uuid).Result()
		if err != nil || result.(int64) == 0 {
			log.Println("unlock failed")
			continue
		}
		stopChan <- struct{}{}
		return
	}

}
func listenResetTime(ctx context.Context, uuid string) {
	tick := time.NewTicker(timeout / 2)
	defer tick.Stop()
	script := redis.NewScript(resetScript)
	for {
		select {
		case <-stopChan:
			return
		case <-tick.C:
			result, err := script.Run(ctx, RedisClient, []string{lockflag}, uuid, timeout).Result()
			if err != nil || result == int64(0) {
				log.Println("reset time failed")
			}
		}
	}
}
