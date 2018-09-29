package wechat

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/spf13/viper"
)

func Login(c *gin.Context) {
	code := c.DefaultQuery("code", "null")
	state := c.Query("state")
	glog.Info(state)
	if code == "null" {
		glog.Infof("接收为空 code 微信获取code失败")
		return
	}
	urlToken := fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?"+
		"appid=%s"+
		"&secret=%s"+
		"&code=%s"+
		"&grant_type=authorization_code", viper.GetString("appid"), viper.GetString("secret"), code)

	//请求access_token
	resp, err := http.Get(urlToken)
	if err != nil {
		glog.Infof("请求access_token错误")
		c.JSONP(http.StatusOK, gin.H{
			"code": "0401",
		})
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Infof("请求access_token body错误")
		c.JSONP(http.StatusOK, gin.H{
			"code": "0402",
		})
	}

	var acc AccessToken
	if err := json.Unmarshal(body, &acc); err != nil {
		glog.Infof("bind is fail")
		c.JSONP(http.StatusOK, gin.H{
			"code": "0403",
		})
	}
	if acc.Errcode == 40029 {
		glog.Infof("error about resp body")
		c.JSONP(http.StatusOK, gin.H{
			"code": "0404",
		})
	}
	c.JSONP(http.StatusOK, gin.H{
		"code": "200",
	})
}
