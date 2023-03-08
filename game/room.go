package game

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"time"
	"xiuianserver/model"
	"xiuianserver/utils"
)

type Room struct {
	Id                uint32
	Players           []*Player
	PendingInput      []interface{} //*model.PeopleMessage
	LastTime          float64
	LastPlayerFrameId map[uint32]uint64
}

func NewRoom(rid uint32) *Room {
	return &Room{
		Id:                rid,
		Players:           make([]*Player, 0),
		PendingInput:      make([]interface{}, 0),
		LastTime:          float64(time.Now().Unix()) / 1e11,
		LastPlayerFrameId: make(map[uint32]uint64),
	}
}
func (r *Room) JoinRoom(uid uint32) {
	if playerMsg, ok := playerManager.Player[uid]; ok {
		playerMsg.Rid = r.Id
		playerMsg.IsRoomMaster = true
		r.Players = append(r.Players, playerMsg)
	}
}
func (r *Room) GetRoomPlayerSlice() (roomPlayerList []model.RoomPlayerMessage) {
	roomPlayerList = make([]model.RoomPlayerMessage, 0)
	for i := 0; i < len(r.Players); i++ {
		rpm := model.RoomPlayerMessage{
			Id:       r.Players[i].Id,
			Nickname: r.Players[i].Nickname,
			Rid:      r.Players[i].Rid,
		}
		roomPlayerList = append(roomPlayerList, rpm)
	}
	return
}
func (r *Room) GetRoomView(player []*Player) (roomList []model.RoomPlayerMessage) {
	roomList = make([]model.RoomPlayerMessage, 0)
	for _, v := range player {
		gpv := v.GetPlayerView(v)
		rpm := model.RoomPlayerMessage{
			Id:       gpv.Id,
			Nickname: gpv.Nickname,
			Rid:      gpv.Rid,
		}
		roomList = append(roomList, rpm)
	}
	return
}
func (r *Room) start() error {
	playerMsgList := make([]*model.PlayerMessage, 0)

	for i := 0; i < len(r.Players); i++ {
		randNum := rand.Intn(301)
		pm := &model.PlayerMessage{
			Id:         r.Players[i].Id,
			Nickname:   r.Players[i].Nickname,
			Type:       utils.Actor1,
			WeaponType: utils.Weapon1,
			BulletType: utils.Bullet2,
			Hp:         100,
			Position: &model.Position{
				X: float64(-150 + randNum),
				Y: float64(-150 + randNum),
			},
			Direction: &model.Direction{
				X: 1,
				Y: 0,
			},
		}
		playerMsgList = append(playerMsgList, pm)
	}
	sl := &model.StateList{
		Actors:       playerMsgList,
		Bullets:      make([]interface{}, 0),
		NextBulletId: 1,
	}
	state := &model.State{
		State: sl,
	}
	rs := &model.ResponState{
		Name: utils.MsgGameStart,
		Data: state,
	}
	result, err := json.Marshal(rs)
	if err != nil {
		//log.Println(err)
		return errors.New("start game json.Marshal err")
	}
	for _, v := range r.Players {
		log.Println("zcc", string(result))
		if err := v.Connection.WriteMessage(websocket.TextMessage, []byte(string(result))); err != nil {
			log.Println("create game send message fail", err)
			return err
		}

	}
	return nil
}
func (r *Room) SendServerMsg() error {
	defer func() {
		if err := recover(); err != nil {
			log.Println("SendServerMsg  recover", err)
			return
		}
	}()
	if r == nil {
		return errors.New("房间不存在")
	}
	inputs := r.PendingInput
	r.PendingInput = nil
	for _, player := range r.Players {
		rm := model.RespMessage{
			Name: utils.MsgServerSync,
			Data: model.RespMessageData{
				LastFrameId: r.LastPlayerFrameId[player.Id],
				Input:       inputs,
			},
		}
		result, err := json.Marshal(rm)
		if err != nil {
			log.Println("json.Marshal fail", err)
			continue
		}
		if err := player.Connection.WriteMessage(websocket.TextMessage, []byte(string(result))); err != nil {
			log.Println("send message fail", err)
			return err
		}
	}
	return nil
}
func (r *Room) TimePast() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("TimePast  recover", err)
			return
		}
	}()
	curTime := float64(time.Now().Unix()) / 1e11
	//if r.LastTime == 0 {
	//	r.LastTime = curTime
	//}
	//dt := curTime - r.LastTime
	//fmt.Println(dt)
	tp := &model.TimePast{
		Type: utils.TimePast,
		Dt:   curTime,
	}
	r.PendingInput = append(r.PendingInput, tp)
	r.LastTime = curTime
}
func (r *Room) GetRoomMaster(rid uint32) string {
	masterName := ""
	for _, player := range r.Players {
		if player.IsRoomMaster {
			masterName = player.Nickname
			return masterName
		}
	}
	return ""
}
