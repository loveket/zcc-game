package nsq

import (
	"github.com/nsqio/go-nsq"
	"log"
	gc "xiuianserver/config"
)

var NsqPub *nsq.Producer

func NewNsqClient() {
	config := nsq.NewConfig()
	pb, err := nsq.NewProducer(gc.GlobalConfig.NsqConfig.RemoteAddr, config)
	if err != nil {
		log.Println("NewProducer err", err)
		return
	}
	NsqPub = pb
}
