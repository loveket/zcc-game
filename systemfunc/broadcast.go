package systemfunc

import (
	"encoding/json"
	"log"
	"sync"
	"xiuianserver/game"
	nsq_pb "xiuianserver/nsq"
)

type SystemMessage struct {
	data string
}

func NewSystemMessage(data string) *SystemMessage {
	return &SystemMessage{data: data}
}

func (sm *SystemMessage) Broadcast() {
	msg, err := json.Marshal(sm.data)
	if err != nil {
		log.Println("[marshal SystemMessage err]", err)
		return
	}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		game.ConnOnlineMap.Range(func(key, value any) bool {
			msgList := value.(*game.OnlineStatusMsg)
			nsq_pb.NsqPub.Publish(msgList.NsqTopic, msg)
			return true
		})
	}()
	wg.Wait()
}
