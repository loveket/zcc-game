package http

import (
	"github.com/gin-gonic/gin"
	"xiuianserver/service"
)

func HttpApi(r *gin.Engine) {

	u := r.Group("user")
	{
		u.POST("/register", service.RegisterUser)
		u.POST("/login", service.LoginUser)
	}
	p := r.Group("play")
	{
		p.GET("/ws/:user", service.WebSocketConn)
	}

}
