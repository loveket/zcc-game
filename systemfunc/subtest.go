package systemfunc

import (
	"errors"
	"fmt"
	"github.com/nsqio/go-nsq"
	"log"
	"sync"
	"time"
	"xiuianserver/game"
)

// nsq 订阅测试
type Sub struct {
}

func (s *Sub) HandleMessage(msg *nsq.Message) error {
	tick := time.NewTicker(1 * time.Second)
	defer tick.Stop()
	for {
		select {
		case <-time.After(10 * time.Second):
			return errors.New("timeout")
		case <-tick.C:
			fmt.Println("[receive msg]=", string(msg.Body))
		}
	}
	return nil
}
func BroadcastTest() {
	log.Println("start consume test....")
	wg := &sync.WaitGroup{}
	config := nsq.NewConfig()
	players := game.GetPlayerManager()
	for _, player := range players.Player {
		wg.Add(1)
		go func(player *game.Player) {
			defer wg.Done()
			consume, _ := nsq.NewConsumer(player.NsqTopics["SystemFunc"], "hpc", config)
			consume.AddHandler(&Sub{})
			err := consume.ConnectToNSQD("192.168.44.131:4150")
			if err != nil {
				log.Println("ConnectToNSQD err", err)
				return
			}

		}(player)
	}
	wg.Wait()
}
