package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	_ "net/http/pprof"
	"sync/atomic"
	"xiuianserver/connection"
)

var cid uint32 = 0

func wsHandle(writer http.ResponseWriter, request *http.Request) {
	if request.RequestURI != "/" {
		return
	}
	upgrader := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
		return true
	}}
	conn, err := upgrader.Upgrade(writer, request, nil)
	//defer conn.Close()
	if err != nil {
		writer.Write([]byte(err.Error()))
	}
	atomic.AddUint32(&cid, 1)
	dealConn := connection.NewConnection(conn, atomic.LoadUint32(&cid))

	go dealConn.Start()
}
func main() {
	log.Println("start server...")
	go connection.ListenPlayerMap()
	http.HandleFunc("/", wsHandle)
	http.ListenAndServe(":7000", nil)

}
