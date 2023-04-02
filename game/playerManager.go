package game

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"sync"
	"xiuianserver/model"
	"xiuianserver/utils"
)

var playerManager *PlayerManager

type PlayerManager struct {
	Player map[string]*Player
	lock   sync.RWMutex
}

func GetPlayerManager() *PlayerManager {
	if playerManager == nil {
		playerManager = new(PlayerManager)
		playerManager.Player = make(map[string]*Player)
	}
	return playerManager
}

//	func (pm *PlayerManager) PlayerLogin(conn *websocket.Conn) (bool, error) {
//		var player *Player
//		flag := false
//		if player == nil {
//			_, data, err := conn.ReadMessage()
//			if err != nil {
//				log.Println("read data err", err)
//				return flag, err
//			}
//			rpt := model.ReqPlayerType{}
//			err = json.Unmarshal(data, &rpt)
//			if err != nil {
//				log.Println("Unmarshal err", err)
//				return flag, err
//			}
//			if rpt.Name == utils.ApiPlayerJoin {
//				//fmt.Println("client message nickname" + string(data))
//				rpm := model.ReqPlayerMessage{}
//				err = json.Unmarshal(data, &rpm)
//				if err != nil || len(rpm.Nickname) == 0 {
//					log.Println("Unmarshal err", err)
//					return flag, err
//				}
//				player := GetPlayerManager().CreatePlayer(rpm.Nickname, conn)
//				result := player.GetPlayerView(player)
//				z := model.Response{Name: utils.ApiPlayerJoin, Data: model.ResponseBody{Success: true, Error: "", Res: result}}
//				result1, err := json.Marshal(z)
//				if err != nil {
//					log.Println("Marshal err", err)
//					return flag, err
//				}
//				conn.WriteMessage(websocket.TextMessage, []byte(string(result1)))
//				flag = true
//			}
//		}
//		return flag, nil
//	}
func (pm *PlayerManager) AddPlayer(player *Player) {
	pm.lock.Lock()
	defer pm.lock.Unlock()
	pm.Player[player.Id] = player
	//playSlice := player.GetPlayerView(player)
	//playerList = append(playerList, playSlice)
	//广播玩家
	//pm.Broadcast()
}
func (pm *PlayerManager) RemovePlayer(id string) {
	pm.lock.Lock()
	defer pm.lock.Unlock()
	delete(pm.Player, id)
}
func (pm *PlayerManager) GetPlayerList() (playList []model.RespPlayer) {
	pm.lock.Lock()
	defer pm.lock.Unlock()
	playList = make([]model.RespPlayer, 0)
	for _, player := range pm.Player {
		rp := model.RespPlayer{
			Id:       player.Id,
			Nickname: player.Nickname,
			Rid:      player.Rid,
		}
		playList = append(playList, rp)
	}
	return
}
func (pm *PlayerManager) Broadcast() {
	playList := make([]model.RespPlayer, 0)
	for _, player := range pm.Player {
		rp := model.RespPlayer{
			Id:       player.Id,
			Nickname: player.Nickname,
			Rid:      player.Rid,
		}
		playList = append(playList, rp)
	}
	for _, player := range pm.Player {
		rpl := model.RespPlayerList{Name: utils.MsgPlayerList, Data: playList}
		result, _ := json.Marshal(rpl)
		fmt.Println("----", string(result))
		if err := player.Connection.Conn.WriteMessage(websocket.TextMessage, []byte(string(result))); err != nil {
			log.Println("send message fail", err)
			continue
		}
	}
}
