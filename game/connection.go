package game

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net"
	"sync"
	"time"
	"xiuianserver/model"
)

var ConnOnlineMap sync.Map

type Connection struct {
	Conn     *websocket.Conn
	ConnID   string
	isClosed bool
	ExitChan chan bool
	msgChan  chan interface{}
	lock     sync.Mutex
}
type OnlineStatusMsg struct {
	RemoteAddr net.Addr
	NsqTopic   string
}
type RespPlayerMessage struct {
	Name string           `json:"name"`
	Data model.RespPlayer `json:"data"`
}

func ListenPlayerMap() {
	tick := time.NewTicker(3 * time.Second)
	for {
		select {
		case <-tick.C:
			//log.Println(ConnectionMap)
			//log.Println(game.GetPlayerManager().Player)
			//case ok := <-conn.ExitChan:
			//	if ok {
			//
			//	}
		}
	}
}
func NewConnection(conn *websocket.Conn, connID string) *Connection {
	return &Connection{
		Conn:     conn,
		ConnID:   connID,
		isClosed: false,
		msgChan:  make(chan interface{}, 100),
		ExitChan: make(chan bool, 1),
	}
}

func (conn *Connection) Start() {
	log.Println("start game server...")
	//人物登录
	//flag, err := game.GetPlayerManager().PlayerLogin(conn.Conn)
	//if err != nil {
	//	log.Println("player login err", err)
	//	return
	//}
	//if flag {
	//	go conn.WSRead()
	//	go conn.WSWrite()
	//	osm := &OnlineStatusMsg{
	//		conn.Conn.RemoteAddr(),
	//		"systembroadcast" + strconv.Itoa(int(conn.GetConnID())),
	//	}
	//
	//	ConnOnlineMap.Store(conn.ConnID, osm)
	//}
	go conn.WSRead()
	go conn.WSWrite()
	osm := &OnlineStatusMsg{
		conn.Conn.RemoteAddr(),
		"systembroadcast" + conn.GetConnID(),
	}
	ConnOnlineMap.Store(conn.ConnID, osm)
}
func (conn *Connection) WSRead() {
	log.Println(conn.ConnID, "start read")
	defer func() {
		if err := recover(); err != nil {
			log.Println("ReadClientMsg  recover", err)
		}
	}()
	defer conn.Stop()
	ncb := NewHandlerConnBranch(conn)
	for {
		select {

		default:
			_, data, err := conn.Conn.ReadMessage()
			if err != nil {
				log.Println("read data err", err)
				continue
			}
			fmt.Println("client message " + string(data))
			rpt := model.ReqPlayerType{}
			err = json.Unmarshal(data, &rpt)
			if err != nil {
				log.Println("Unmarshal err", err)
				return
			}
			ncb.HandlerBranch(rpt.Name, data)
		}
	}
}
func (conn *Connection) WSWrite() {
	log.Println(conn.ConnID, "start write")
	defer func() {
		if err := recover(); err != nil {
			log.Println("SendServerMsg  recover", err)
		}
	}()
	defer conn.Stop()
	defer log.Println(conn.Conn.RemoteAddr().String(), "conn writer exit")

	for {
		select {
		case data := <-conn.msgChan:
			result, _ := json.Marshal(data)
			if err := conn.Conn.WriteMessage(websocket.TextMessage, []byte(string(result))); err != nil {
				log.Println("send message fail", err)
				continue
			}
		case <-conn.ExitChan:
			return
		}
	}
}
func (conn *Connection) Stop() {
	//conn.lock.Lock()
	//defer conn.lock.Unlock()
	log.Println("conn stop...connid", conn.ConnID)
	playerMsg := GetPlayerManager().Player[conn.GetConnID()]
	GetRoomManager().PlayerLeaveRoom(playerMsg.Rid, playerMsg.Id)
	GetPlayerManager().RemovePlayer(playerMsg.Id)
	GetPlayerManager().Broadcast()
	GetRoomManager().Broadcast()
	GetRoomManager().BroadcastRoomPlayer(playerMsg.Rid)
	ConnOnlineMap.Delete(conn.GetConnID())
	if conn.isClosed == true {
		return
	}
	conn.isClosed = true
	conn.Conn.Close()
	//告知writer关闭
	conn.ExitChan <- true
	//ConnectionMap.Delete(conn.GetConnID())
	//回收资源
	close(conn.ExitChan)
	close(conn.msgChan)
}
func (conn *Connection) GetConnID() string {
	return conn.ConnID
}
