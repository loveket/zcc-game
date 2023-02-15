package game

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"sync"
	"time"
	"xiuianserver/model"
	"xiuianserver/utils"
)

var roomId uint32 = 0

var roomManager *RoomManager

type RoomManager struct {
	Room map[uint32]*Room
	lock sync.RWMutex
}

func GetRoomManager() *RoomManager {
	if roomManager == nil {
		roomManager = new(RoomManager)
		roomManager.Room = make(map[uint32]*Room)
	}
	return roomManager
}

func (rm *RoomManager) CreateRoom() *Room {
	rm.lock.Lock()
	defer rm.lock.Unlock()
	roomId++
	room := NewRoom(roomId)
	rm.Room[roomId] = room
	return room
}
func (rm *RoomManager) DeleteRoom(rid uint32) {
	roomId--
	delete(rm.Room, rid)
	rm.Broadcast()
}
func (rm *RoomManager) PlayJoinRoom(rid, uid uint32) *Room {
	rm.lock.Lock()
	defer rm.lock.Unlock()
	if roomMsg, ok := rm.Room[rid]; ok {
		roomMsg.JoinRoom(uid)
		return roomMsg
	}
	return nil
}
func (rm *RoomManager) GetRoomList() (roomList []model.RespRoomMessage) {
	roomList = make([]model.RespRoomMessage, 0)
	for roomid, room := range rm.Room {
		grv := room.GetRoomView(room.Players)
		master := room.GetRoomMaster(roomid)
		if master == "" || len(master) == 0 {
			log.Println("房间异常")
			rm.DeleteRoom(roomId)
			return
		}
		rp := model.RespRoomMessage{
			Id:         roomid,
			MasterName: master,
			Players:    grv,
		}
		roomList = append(roomList, rp)
	}
	return
}
func (rm *RoomManager) PlayerLeaveRoom(rid, uid uint32) error {
	rm.lock.Lock()
	defer rm.lock.Unlock()
	if roomMsg, ok := rm.Room[rid]; ok {
		roomPlay := roomMsg.Players
		roomSize := len(roomPlay)
		if roomSize < 1 {
			delete(rm.Room, rid)
			return errors.New("空房间")
		}
		if roomSize == 1 /*&& roomMsg.Players[0].Id == uid*/ {
			rm.DeleteRoom(rid)
			return nil
		}
		for i := 0; i < roomSize; i++ {
			if uid == roomMsg.Players[i].Id {
				playerList := append(roomMsg.Players[:i], roomMsg.Players[i+1:]...)
				rm.Room[rid].Players = playerList
				return nil
			}
		}
	}
	return nil
}
func (rm *RoomManager) Broadcast() {
	roomList := rm.GetRoomList()
	//广播给玩家
	for _, player := range GetPlayerManager().Player {
		fmt.Println("333333333333", player)
		rpl := model.RespRoomList{Name: utils.MsgRoomList, Data: roomList}
		result, _ := json.Marshal(rpl)
		fmt.Println("----", string(result))
		if err := player.Connection.WriteMessage(websocket.TextMessage, []byte(string(result))); err != nil {
			log.Println("send message fail", err)
			continue
		}
	}

}
func (rm *RoomManager) BroadcastRoomPlayer(rid uint32) {
	if roomList, ok := rm.Room[rid]; ok {
		roomPlayer := roomList.GetRoomPlayerSlice()
		rpl := model.RespRoomPlayerList{Name: utils.MsgRoom, Data: roomPlayer}
		result, _ := json.Marshal(rpl)
		for _, player := range roomList.Players {
			if err := player.Connection.WriteMessage(websocket.TextMessage, []byte(string(result))); err != nil {
				log.Println("send message fail", err)
				continue
			}
		}
	}
}
func (rm *RoomManager) StartRoom(rid uint32) error {
	if roomList, ok := rm.Room[rid]; ok {
		err := roomList.start()
		if err != nil {
			return err
		}
		//开启同步
		go rm.StartRoomListener(rid)
	} else {
		return errors.New("房间不存在")
	}
	return nil
}
func (rm *RoomManager) StartRoomListener(rid uint32) {
	tick := time.NewTicker(100 * time.Millisecond)
	tick1 := time.NewTicker(16 * time.Millisecond)
	defer tick.Stop()
	defer tick1.Stop()
	for {
		select {
		case <-tick.C:
			if room, ok := rm.Room[rid]; ok {
				err := room.SendServerMsg()
				if err != nil {
					return
				}
			}
		case <-tick1.C:
			if room, ok := rm.Room[rid]; ok {
				room.TimePast()
			}
		}
	}
}
