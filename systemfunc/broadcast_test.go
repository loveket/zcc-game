package systemfunc

import (
	"github.com/nsqio/go-nsq"
	"log"
	"strconv"
	"sync"
	"testing"
	"time"
	"xiuianserver/connection"
)

func init() {
	for i := 0; i < 50; i++ {
		osm := &connection.OnlineStatusMsg{
			nil,
			"systembroadcast" + strconv.Itoa(i),
		}
		connection.ConnOnlineMap.Store(i, osm)
	}
}
func NsqConsumeTest() {
	log.Println("start consume test....")
	wg := &sync.WaitGroup{}
	config := nsq.NewConfig()
	connection.ConnOnlineMap.Range(func(key, value any) bool {
		wg.Add(1)
		data := value.(*connection.OnlineStatusMsg)
		consume, _ := nsq.NewConsumer(data.NsqTopic, "hpc", config)
		consume.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
			log.Println("[sub]===>", string(message.Body))
			wg.Done()
			return nil
		}))
		err := consume.ConnectToNSQD("192.168.44.129:4150")
		if err != nil {
			log.Println("ConnectToNSQD err", err)
			return false
		}
		return true
	})
	wg.Wait()
}
func TestSystemMessage_Broadcast(t *testing.T) {
	go NsqConsumeTest()
	time.Sleep(3 * time.Second)
	data := "系统将在5分钟后维护。。。"
	sm := NewSystemMessage(data)
	sm.Broadcast()
}
