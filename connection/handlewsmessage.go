package connection

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"xiuianserver/game"
	"xiuianserver/model"
	"xiuianserver/utils"
)

type IBranchMethod interface {
	apiPlayerList()
	apiRoomCreate()
	apiRoomList()
	apiRoomJoin(data []byte)
	apiRoomLeave()
	apiGameStart()
	msgClientSync(data []byte)
	msgHp(data []byte)
	apiGameKaPool(data []byte)
}
type handlerConnBranch struct {
	conn *Connection
}

func NewHandlerConnBranch(conn *Connection) *handlerConnBranch {
	return &handlerConnBranch{conn: conn}
}
func (h *handlerConnBranch) HandlerBranch(rpt string, data []byte) {
	switch rpt {
	case utils.ApiPlayerList:
		h.apiPlayerList()
	case utils.ApiRoomCreate:
		h.apiRoomCreate()
	case utils.ApiRoomList:
		h.apiRoomList()
	case utils.ApiRoomJoin:
		h.apiRoomJoin(data)
	case utils.ApiRoomLeave:
		h.apiRoomLeave()
	case utils.ApiChatSend:
		h.apiChatSend(data)
	case utils.ApiGameStart:
		h.apiGameStart()
	case utils.MsgClientSync:
		h.msgClientSync(data)
	case utils.MsgHp:
		h.msgHp(data)
	case utils.ApiGameKaPool:
		h.apiGameKaPool(data)
	}
}
func (h *handlerConnBranch) apiPlayerList() {
	playerList := game.GetPlayerManager().GetPlayerList()
	resp := model.Response{Name: utils.ApiPlayerList, Data: model.ResponseBody{Success: true, Error: "", Res: playerList}}
	game.GetPlayerManager().Broadcast()
	game.GetRoomManager().Broadcast()
	h.conn.msgChan <- resp
}
func (h *handlerConnBranch) apiRoomCreate() {
	if !h.conn.isClosed {
		cr := game.GetRoomManager().CreateRoom()
		room := game.GetRoomManager().PlayJoinRoom(cr.Id, h.conn.GetConnID())
		if room == nil {
			log.Println("创建房间失败")
			return
		}
		gps := room.GetRoomPlayerSlice()
		resp := model.Response{Name: utils.ApiRoomCreate, Data: model.ResponseBody{Success: true, Error: "", Res: model.RespRoomMessage{
			Id:      room.Id,
			Players: gps,
		}}}
		h.conn.msgChan <- resp
		game.GetRoomManager().Broadcast()
	} else {
		log.Println("未登录，不能创建房间")
		return
	}
}
func (h *handlerConnBranch) apiRoomList() {
	roomList := game.GetRoomManager().GetRoomList()
	//rpl := model.RespRoomList{Name: utils.ApiRoomList, Data: roomList}
	resp := model.Response{Name: utils.ApiPlayerList, Data: model.ResponseBody{Success: true, Error: "", Res: roomList}}
	h.conn.msgChan <- resp
}
func (h *handlerConnBranch) apiRoomJoin(data []byte) {
	reqrid := model.ReqRoomId{}
	err := json.Unmarshal(data, &reqrid)
	if err != nil {
		log.Println("Unmarshal err", err)
		return
	}
	game.GetRoomManager().PlayJoinRoom(reqrid.Rid, h.conn.GetConnID())
	roomPlayerList := game.GetRoomManager().Room[reqrid.Rid].GetRoomPlayerSlice()
	resp := model.Response{Name: utils.ApiRoomJoin, Data: model.ResponseBody{Success: true, Error: "", Res: model.RespRoomMessage{
		Id:      reqrid.Rid,
		Players: roomPlayerList,
	}}}
	if roomList, ok := game.GetRoomManager().Room[reqrid.Rid]; ok {
		for _, player := range roomList.Players {
			result, _ := json.Marshal(resp)
			fmt.Println("----", string(result))
			if err := player.Connection.WriteMessage(websocket.TextMessage, []byte(string(result))); err != nil {
				log.Println("send message fail", err)
				continue
			}
		}
	}
	game.GetRoomManager().BroadcastRoomPlayer(reqrid.Rid)
}
func (h *handlerConnBranch) apiChatSend(data []byte) {
	req := model.ReqChatHall{}
	err := json.Unmarshal(data, &req)
	if err != nil {
		log.Println("Unmarshal err", err)
		return
	}
	//Todo   redis后续引入
	player := game.GetPlayerManager().Player[h.conn.GetConnID()]
	resp := model.RespChatHall{
		NickName: player.GetNickname(),
		Time:     req.Time,
		Message:  req.Message,
	}
	rpl := model.ResponState{Name: utils.MsgPlayerHallChat, Data: resp}
	result, _ := json.Marshal(rpl)
	go func() {
		for _, player := range game.GetPlayerManager().Player {
			if err := player.Connection.WriteMessage(websocket.TextMessage, []byte(string(result))); err != nil {
				log.Println("send message fail", err)
				continue
			}
		}
		return
	}()
}
func (h *handlerConnBranch) apiRoomLeave() {
	if player, ok := game.GetPlayerManager().Player[h.conn.GetConnID()]; ok {
		err := game.GetRoomManager().PlayerLeaveRoom(player.Rid, player.Id)
		if err != nil {
			fmt.Println("玩家移除储物", err)
			return
		}
		game.GetPlayerManager().Broadcast()
		game.GetRoomManager().Broadcast()
		game.GetRoomManager().BroadcastRoomPlayer(player.Rid)
		resp := model.Response{Name: utils.ApiRoomLeave, Data: model.ResponseBody{Success: true, Error: "", Res: ""}}
		h.conn.msgChan <- resp
		//fmt.Println("after", game.GetRoomManager().Room[player.Rid].Players)
	} else {
		fmt.Println("玩家不存在", h.conn.GetConnID())
		return
	}
}
func (h *handlerConnBranch) apiGameStart() {
	playerMsg := game.GetPlayerManager().Player[h.conn.ConnID]
	err := game.GetRoomManager().StartRoom(playerMsg.Rid)
	var resp model.Response
	if err != nil {
		log.Println(err)
		resp = model.Response{Name: utils.ApiGameStart, Data: model.ResponseBody{Success: false, Error: err.Error(), Res: ""}}
	}
	resp = model.Response{Name: utils.ApiGameStart, Data: model.ResponseBody{Success: true, Error: "", Res: ""}}
	h.conn.msgChan <- resp
}
func (h *handlerConnBranch) apiGameKaPool(data []byte) {
	t := &model.ReqKaPool{}
	var err error
	err = json.Unmarshal(data, &t)
	if err != nil {
		log.Println("Unmarshal err", err)
		return
	}
	var result []string
	playerMsg := game.GetPlayerManager().Player[h.conn.ConnID]
	if t.Times == 1 {
		result = append(result, playerMsg.KaPool.OneLuckyDraw()...)
	} else if t.Times == 10 {
		result = append(result, playerMsg.KaPool.TenLuckyDraw()...)
	} else {
		err = errors.New("请求抽卡出错")
		return
	}
	fmt.Println("******", result)
	var resp model.Response
	if err != nil {
		log.Println(err)
		resp = model.Response{Name: utils.ApiGameKaPool, Data: model.ResponseBody{Success: false, Error: err.Error(), Res: ""}}
		return
	}
	resp = model.Response{Name: utils.ApiGameKaPool, Data: model.ResponseBody{Success: true, Error: "", Res: result}}
	h.conn.msgChan <- resp
}
func (h *handlerConnBranch) msgClientSync(data []byte) {
	reqmsg := model.ReqMessageData{}
	err := json.Unmarshal(data, &reqmsg)
	if err != nil {
		log.Println("Unmarshal err", err)
		return
	}
	playerMsg := game.GetPlayerManager().Player[h.conn.ConnID]
	roomMsg := game.GetRoomManager().Room[playerMsg.Rid]
	if reqmsg.Input.Type == utils.WeaponShoot {
		rwm := &model.RespWeaponShootMessage{
			Owner:     reqmsg.Input.Id,
			Position:  reqmsg.Input.Position,
			Direction: reqmsg.Input.Direction,
			Type:      reqmsg.Input.Type,
		}
		roomMsg.PendingInput = append(roomMsg.PendingInput, rwm)
	} else {
		roomMsg.PendingInput = append(roomMsg.PendingInput, &reqmsg.Input)
	}
	roomMsg.LastPlayerFrameId[playerMsg.Id] = reqmsg.FrameId
	fmt.Println("MsgClientSync2", reqmsg)
}
func (h *handlerConnBranch) msgHp(data []byte) {
	actors := model.ReqActorsMsg{}
	err := json.Unmarshal(data, &actors)
	if err != nil {
		log.Println("Unmarshal err", err)
		return
	}
	//log.Println("actors+++++++++++++", actors.Data)
	for _, actor := range actors.Data {
		player := game.GetPlayerManager().Player[actor.Id]
		if game.GetRoomManager().Room[player.Rid] == nil {
			continue
		}
		if len(game.GetRoomManager().Room[player.Rid].Players) == 1 {
			resp := model.Response{Name: utils.MsgGameEnd, Data: model.ResponseBody{Success: true, Error: "", Res: "you win"}}
			result, _ := json.Marshal(resp)
			if err := player.Connection.WriteMessage(websocket.TextMessage, []byte(string(result))); err != nil {
				log.Println("send message fail", err)
			}
			err := game.GetRoomManager().PlayerLeaveRoom(player.Rid, player.Id)
			if err != nil {
				log.Println("leave room err", err)
			}
		}
		if actor.Hp <= 0 {
			resp := model.Response{Name: utils.MsgGameEnd, Data: model.ResponseBody{Success: false, Error: "", Res: "you loser"}}
			result, _ := json.Marshal(resp)
			if err := player.Connection.WriteMessage(websocket.TextMessage, []byte(string(result))); err != nil {
				log.Println("send message fail", err)
			}
			err := game.GetRoomManager().PlayerLeaveRoom(player.Rid, player.Id)
			if err != nil {
				log.Println("leave room err", err)
			}
		}
	}
}
