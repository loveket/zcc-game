package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"xiuianserver/game"
	"xiuianserver/log"
	"xiuianserver/model"
	"xiuianserver/redis"
)

var ctx = context.Background()

func WebSocketConn(c *gin.Context) {
	username := c.Param("user")
	fmt.Println("*****", username)
	if len(username) == 0 || username == "null" {
		c.JSON(201, "状态异常")
		return
	}
	upgrader := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
		//if r.Header["user"] == nil || r.Header["user"][0] == "" || len(r.Header["user"][0]) == 0 {
		//	return false
		//}
		//username = r.Header["user"][0]
		return true
	}}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.Writer.Write([]byte(err.Error()))
		return
	}
	//atomic.AddUint32(&cid, 1)
	//dealConn := connection.NewConnection(conn, user)

	user := &model.UserModel{}
	err = json.Unmarshal([]byte(redis.RedisClient.HGet(ctx, "user", username).Val()), user)
	if err != nil {
		fmt.Println("login Unmarshal err", err)
		return
	}
	player := game.NewPlayer(user.Id, user.Username, conn)
	playerG := game.GetPlayerManager()
	playerG.AddPlayer(player)
	go player.Connection.Start()
}
func RegisterUser(c *gin.Context) {
	user := &model.UserModel{}
	err := c.ShouldBind(user)
	if err != nil {
		fmt.Println("bind user failed", err)
		return
	}
	user.Id = uuid.NewV4().String()
	if ok, err := redis.RedisClient.HExists(ctx, "user", user.Username).Result(); ok && err == nil {
		c.JSON(201, "当前用户已经存在")
		return
	}
	if len(user.Username) == 0 || len(user.Password) == 0 {
		c.JSON(201, "用户名或密码不能为空")
		return
	}
	data, err := json.Marshal(user)
	if err != nil {
		fmt.Println("user Marshal", err)
		return
	}
	result1 := redis.RedisClient.HSet(ctx, "user", user.Username, data)
	if result1.Err() != nil {
		fmt.Println("redis存储用户失败", result1.Err())
		c.JSON(201, "状态异常，请重新注册")
		return
	}
	c.JSON(200, "注册成功")
}
func LoginUser(c *gin.Context) {
	user1 := &model.UserModel{}
	err := c.ShouldBind(user1)
	if err != nil {
		fmt.Println("bind user failed", err)
		return
	}
	if ok, err := redis.RedisClient.HExists(ctx, "user", user1.Username).Result(); !ok && err == nil {
		c.JSON(201, "用户不存在")
		return
	}
	user2 := &model.UserModel{}
	err = json.Unmarshal([]byte(redis.RedisClient.HGet(ctx, "user", user1.Username).Val()), user2)
	if err != nil {
		fmt.Println("login Unmarshal err", err)
		return
	}
	if user1.Username == user2.Username && user1.Password == user2.Password {
		log.LoggerSingle.Info(user2.Username + "登录成功")
		c.JSON(200, gin.H{"data": user2.Username, "id": user2.Id})
		return
	}
	c.JSON(201, "用户名或密码错误")
}
