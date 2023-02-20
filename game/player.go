package game

import (
	"github.com/gorilla/websocket"
	"xiuianserver/model"
)

type Player struct {
	Id           uint32
	Nickname     string
	Connection   *websocket.Conn
	Rid          uint32
	IsRoomMaster bool
	KaPool       *KaPool
}

func NewPlayer(id uint32, nickname string, conn *websocket.Conn) *Player {
	return &Player{
		Id:         id,
		Nickname:   nickname,
		Connection: conn,
		KaPool:     &KaPool{},
	}
}
func (p *Player) GetNickname() string {

	return p.Nickname
}
func (p *Player) GetPlayerView(player *Player) model.RespPlayer {
	return model.RespPlayer{
		Id:       playerId,
		Nickname: player.Nickname,
		Rid:      player.Rid,
	}
}

//func (p *Player) ReadMsgFromConn() (inputs *model.PeopleMessage, err error) {
//	fmt.Println("8888888888888888888888888888")
//	for {
//		_, data, err := p.Connection.ReadMessage()
//		if err != nil {
//			log.Println("read data err", err)
//			return nil, err
//		}
//		fmt.Println("client message from ws" + string(data))
//		fmt.Println("8888888888899999999999999")
//		rpt := model.ReqPlayerType{}
//		err = json.Unmarshal(data, &rpt)
//		if err != nil {
//			log.Println("Unmarshal err", err)
//			return nil, err
//		}
//		if rpt.Name == utils.MsgClientSync {
//			reqmsg := model.ReqMessageData{}
//			err = json.Unmarshal(data, &reqmsg)
//			if err != nil {
//				log.Println("Unmarshal err", err)
//				return nil, err
//			}
//			return &reqmsg.Input, nil
//		}
//	}
//}
