package utils

import (
	"github.com/gin-gonic/gin"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method //请求方法
		if method == "OPTIONS" {
			c.Status(200)
		}
		c.Header("Access-Control-Allow-Origin", "*") // 这是允许访问所有域

		c.Header("Access-Control-Allow-Methods", "GET,DELETE,POST,PUT,OPTIONS") //服务器支持的所有跨域请求的方法,为了避免浏览次请求的多次'预检'请求
		//  header的类型
		c.Header("Access-Control-Allow-Headers", "*")
		//              允许跨域设置                                                                                                      可以返回其他子段
		//c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar") // 跨域关键设置 让浏览器可以解析
		//c.Header("Access-Control-Max-Age", "172800")                                                                                                                                                           // 缓存请求信息 单位为秒
		//c.Header("Access-Control-Allow-Credentials", "true") //  跨域请求是否需要带cookie信息 默认设置为true
		c.Set("content-type", "application/json") // 设置返回格式是json

		// 处理请求
		c.Next() //  处理请求
	}
}
