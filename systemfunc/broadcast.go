package systemfunc

import (
	"encoding/json"
	"github.com/nsqio/go-nsq"
	"log"
	"sync"
	"xiuianserver/connection"
)

type SystemMessage struct {
	data string
}

func NewSystemMessage(data string) *SystemMessage {
	return &SystemMessage{data: data}
}

func (sm *SystemMessage) Broadcast() {
	config := nsq.NewConfig()
	product, err := nsq.NewProducer("192.168.44.129:4150", config)
	if err != nil {
		log.Println("[nsq product err]", err)
		return
	}
	msg, err := json.Marshal(sm.data)
	if err != nil {
		log.Println("[marshal SystemMessage err]", err)
		return
	}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		connection.ConnOnlineMap.Range(func(key, value any) bool {
			msgList := value.(*connection.OnlineStatusMsg)
			product.Publish(msgList.NsqTopic, msg)
			return true
		})
	}()
	wg.Wait()
}
