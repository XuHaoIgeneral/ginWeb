package wechatLogin

import (
	"ginweb/server/LoginS/wxLoginS"
	"ginweb/util/aesED"
	"github.com/gin-gonic/gin"
	"net/http"
)

//接收小程序登陆请求并返回
func XcxLogin(c *gin.Context) {
	
	code := c.DefaultPostForm("code", "null")
	xcx, err := wxLoginS.GetSessionKey(code)
	token, err := aesED.Encrypt(xcx.Openid)
	if err != nil {
		c.JSONP(http.StatusOK, gin.H{
			"code": "0",
		})
		return
	}
	c.JSONP(http.StatusOK, gin.H{
		"token": token,
	})
}
