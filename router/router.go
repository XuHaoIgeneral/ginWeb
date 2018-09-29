package router

import (
	"ginWeb/handler/wechat"
	"ginWeb/router/middleware"
	"ginWeb/handler/vaptcha"
	"ginweb/handler/sd"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Load(g *gin.Engine,mv ...gin.HandlerFunc) *gin.Engine  {
	//Middlewares

	g.Use(gin.Recovery())
	g.Use(middleware.Options)
	g.Use(middleware.Secure)
	g.Use(mv...)

	//404 Handler
	g.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound,"The incorrect API route")
	})

	//The health check heandlers
	svcd:=g.Group("/sd")
	{
		svcd.GET("/health",sd.HealthCheck)
		svcd.GET("/disk",sd.DiskCheck)
		svcd.GET("/cpu",sd.CPUCheck)
		svcd.GET("/ram",sd.RAMCheck)
	}

	//The vaptcha
	vap:=g.Group("/test")
	{
		vap.POST("/validate",vaptcha.Validate)

		vap.GET("/downtime",vaptcha.Outage)

	}

	wx:=g.Group("/")
	{
		wx.Any("/login",wechat.Login)
	}

	return g
}
