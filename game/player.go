package game

import (
	"github.com/gorilla/websocket"
	"xiuianserver/model"
)

type Player struct {
	Id           string
	Nickname     string
	Rid          uint32
	IsRoomMaster bool
	NsqTopics    map[string]string
	Connection   *Connection
	KaPool       *KaPool
	Friend       *Friend
	OnlineStatus bool
}

func NewPlayer(id string, nickname string, conn *websocket.Conn) *Player {
	//nsq主题名为玩家uuid
	topics := map[string]string{"SystemFunc": id}
	return &Player{
		Id:         id,
		Nickname:   nickname,
		NsqTopics:  topics,
		Connection: NewConnection(conn, id),
		KaPool:     &KaPool{},
	}
}
func (p *Player) GetNickname() string {

	return p.Nickname
}
func (p *Player) GetPlayerView(player *Player) model.RespPlayer {
	return model.RespPlayer{
		Id:       player.Id,
		Nickname: player.Nickname,
		Rid:      player.Rid,
	}
}
