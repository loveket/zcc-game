package model

type RespPlayerList struct {
	Name string       `json:"name"`
	Data []RespPlayer `json:"data"`
}
type RespRoomPlayerList struct {
	Name string              `json:"name"`
	Data []RoomPlayerMessage `json:"data"`
}
type RespPlayer struct {
	Id       string `json:"id"`
	Nickname string `json:"nickname"`
	Rid      uint32 `json:"rid"`
}
type RespRoom struct {
	Name string          `json:"name"`
	Data RespRoomMessage `json:"data"`
}
type RespRoomList struct {
	Name string            `json:"name"`
	Data []RespRoomMessage `json:"data"`
}
type RespRoomMessage struct {
	Id         uint32              `json:"id"`
	MasterName string              `json:"masterName"`
	Players    []RoomPlayerMessage `json:"players"`
}
type RoomPlayerMessage struct {
	Id       string `json:"id"`
	Nickname string `json:"nickname"`
	Rid      uint32 `json:"rid"`
}
type ReqPlayerType struct {
	Name string `json:"name"`
}
type ReqPlayerMessage struct {
	Nickname string `json:"nickname"`
}

// **********
type Response struct {
	Name string       `json:"name"`
	Data ResponseBody `json:"data"`
}
type ResponseBody struct {
	Success bool        `json:"success"`
	Error   string      `json:"error"`
	Res     interface{} `json:"res"`
}

type ResponState struct {
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}
type State struct {
	State *StateList `json:"state"`
}
type StateList struct {
	Actors       []*PlayerMessage `json:"actors"`
	Bullets      []interface{}    `json:"bullets"`
	NextBulletId uint32           `json:"nextBulletId"`
}
type PlayerMessage struct {
	Id         string     `json:"id"`
	Nickname   string     `json:"nickname"`
	Type       string     `json:"type"`
	WeaponType string     `json:"weaponType"`
	BulletType string     `json:"bulletType"`
	Hp         int64      `json:"hp"`
	Position   *Position  `json:"position"`
	Direction  *Direction `json:"direction"`
}
type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}
type Direction struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type RespMessage struct {
	Name string          `json:"name"`
	Data RespMessageData `json:"data"`
}
type ReqMessageData struct {
	FrameId uint64        `json:"frameId"`
	Input   PeopleMessage `json:"input"`
}
type ReqRoomId struct {
	Rid uint32 `json:"rid"`
}
type ReqChatHall struct {
	Time    string `json:"time"`
	Message string `json:"message"`
}
type RespChatHall struct {
	NickName string `json:"nickname"`
	Time     string `json:"time"`
	Message  string `json:"message"`
}
type RespMessageData struct {
	LastFrameId uint64        `json:"lastFrameId"`
	Input       []interface{} `json:"inputs"`
}
type PeopleMessage struct {
	Id        uint64  `json:"id"`
	Type      string  `json:"type"`
	Position  Pos     `json:"position"`
	Direction Pos     `json:"direction"`
	Dt        float64 `json:"dt"`
}
type RespWeaponShootMessage struct {
	Owner     uint64 `json:"owner"`
	Type      string `json:"type"`
	Position  Pos    `json:"position"`
	Direction Pos    `json:"direction"`
}
type Pos struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}
type TimePast struct {
	Type string  `json:"type"`
	Dt   float64 `json:"dt"`
}
type ReqActorsMsg struct {
	Data []*PlayerMessage `json:"data"`
}
type ReqKaPool struct {
	Times int `json:"times"`
}
