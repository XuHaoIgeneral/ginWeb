package wechat

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
)

type XcxSessionKey struct {
	Openid      string `form:"openid" json:"openid" binding:"required"`
	Session_key string `form:"session_key" json:"session_key" binding:"required"`
	Unionid     string `form:"unionid" json:"unionid" binding:"required"`
}

func XcxLogin(c *gin.Context) {

	code := c.DefaultPostForm("code", "null")
	if code == "null" {
		log.Infof("接收为空 code 微信获取code失败")
		c.JSONP(http.StatusOK, gin.H{
			"status": "0400",
		})
	}

	//请求session_key
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?"+
		"appid=%s"+
		"&secret=%s"+
		"&js_code=%s"+
		"&grant_type=%s",
		viper.GetString("wechat.xcx.appid"), viper.GetString("wechat.xcx.secret"), code, "authorization_code")

	fmt.Println(url)
	resp, err := http.Get(url)
	defer resp.Body.Close()


	if err != nil {
		log.Infof("请求session_key错误")
		c.JSONP(http.StatusOK, gin.H{
			"code": "0401",
		})
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Infof("请求session_key body错误")
		c.JSONP(http.StatusOK, gin.H{
			"code": "0402",
		})
		return
	}

	//判断返回参数 获取access_token
	probe := gjson.Get(string(body), "unionid")
	if probe.String() == "" {
		log.Infof("获取unionid 错误 返回为:%s",string(body))
		c.JSONP(http.StatusOK, gin.H{
			"code": "0403",
		})
		return
	}

	var xcx XcxSessionKey
	if err := json.Unmarshal(body, &xcx); err != nil {
		log.Infof("bind is fail")
		c.JSONP(http.StatusOK, gin.H{
			"code": "0403",
		})
	}

	c.JSONP(http.StatusOK, gin.H{
		"code": "200",
	})
}
