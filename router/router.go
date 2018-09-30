package router

import (
	"ginweb/handler/sd"
	"ginweb/handler/vaptcha"
	"ginweb/handler/wechat"
	"ginweb/router/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Load(g *gin.Engine, mv ...gin.HandlerFunc) *gin.Engine {
	//Middlewares

	g.Use(gin.Recovery())
	g.Use(middleware.Options)
	g.Use(middleware.Secure)
	g.Use(mv...)

	//404 Handler
	g.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "The incorrect API route")
	})

	//The health check heandlers
	svcd := g.Group("/sd")
	{
		svcd.GET("/health", sd.HealthCheck)
		svcd.GET("/disk", sd.DiskCheck)
		svcd.GET("/cpu", sd.CPUCheck)
		svcd.GET("/ram", sd.RAMCheck)
	}

	//The vaptcha
	vap := g.Group("/test")
	{
		vap.POST("/validate", vaptcha.Validate)

		vap.GET("/downtime", vaptcha.Outage)

	}

	wx := g.Group("/wx")
	{
		wx.Any("/login", wechat.Login)
	}

	return g
}
