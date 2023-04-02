package main

import (
	"github.com/gin-gonic/gin"
	"log"
	_ "net/http/pprof"
	api "xiuianserver/api/http"
	gc "xiuianserver/config"
	"xiuianserver/game"
	"xiuianserver/redis"
	"xiuianserver/utils"
)

var router *gin.Engine

func init() {
	router = gin.Default()
	router.Use(utils.Cors())
}
func main() {
	log.Println("[zcc-game] serve already start...")
	go game.ListenPlayerMap()
	redis.NewRedis()
	//http.HandleFunc("/", wsHandle)
	//http.ListenAndServe(":7000", nil)
	api.HttpApi(router)
	router.Run(gc.GlobalConfig.HttpConfig.IpAddr)
}
