package systemfunc

import (
	"fmt"
	"xiuianserver/game"
	nsq_pb "xiuianserver/nsq"
)

type SystemFunc struct {
}
type SystemMessageRequest struct {
	Data string
}
type SystemMessageResponse struct {
	Status string
}

//func NewSystemMessage(data string) *SystemMessage {
//	return &SystemMessage{data: data}
//}

//	func (sm *SystemFunc) Broadcast() {
//		msg, err := json.Marshal(sm.data)
//		if err != nil {
//			log.Println("[marshal SystemMessage err]", err)
//			return
//		}
//		wg := &sync.WaitGroup{}
//		wg.Add(1)
//		go func() {
//			defer wg.Done()
//			game.ConnOnlineMap.Range(func(key, value any) bool {
//				msgList := value.(*game.OnlineStatusMsg)
//				nsq_pb.NsqPub.Publish(msgList.NsqTopic, msg)
//				return true
//			})
//		}()
//		wg.Wait()
//	}
func (sm *SystemFunc) Broadcast(req SystemMessageRequest, resp *SystemMessageResponse) error {
	players := game.GetPlayerManager()
	for _, player := range players.Player {
		fmt.Println("[pub]" + player.NsqTopics["SystemFunc"] + "---[Msg]" + req.Data)
		err := nsq_pb.NsqPub.Publish(player.NsqTopics["SystemFunc"], []byte(req.Data))
		if err != nil {
			resp.Status += player.Nickname + "failed;"
			continue
		}
		resp.Status += player.Nickname + "success;"
	}
	return nil
}
