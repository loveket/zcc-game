package nsq

import (
	"github.com/nsqio/go-nsq"
	"log"
	log_diy "xiuianserver/log"
)

type NsqClient struct {
	RemoteAddr string
	Producer   *nsq.Producer
	Log        *log_diy.Log
}

func NewNsqClient(addr string) *NsqClient {
	config := nsq.NewConfig()
	pb, err := nsq.NewProducer(addr, config)
	if err != nil {
		log.Println("NewProducer err", err)
		return nil
	}
	return &NsqClient{RemoteAddr: addr, Producer: pb, Log: log_diy.NewLog(addr)}
}
func (nc *NsqClient) Pub(topic string, message []byte) {
	err := nc.Producer.Publish(topic, message)
	if err != nil {
		nc.Log.ErrorMSG(err.Error())
	}
}
