package wechat

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/spf13/viper"
)

func Login(c *gin.Context) {

	code := c.DefaultQuery("code", "null")

	if code == "null" {
		glog.Infof("接收为空 code 微信获取code失败")
		return
	}

	//请求access_token
	urlToken := fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?"+
		"appid=%s"+
		"&secret=%s"+
		"&code=%s"+
		"&grant_type=authorizatLSion_code", viper.GetString("wechat.xcx.appid"), viper.GetString("wechat.xcx.secret"), code)

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
		return
	}

	var wxerr Err
	//判断返回参数 获取access_token
	probe := gjson.Get(string(body), "errcode")
	if probe.String() != "" {
		glog.Infof("request access fail, body is %s", body)
		if err = json.Unmarshal(body, &wxerr); err != nil {
			glog.Infof("bind wxerr is fail")
			c.JSONP(http.StatusOK, gin.H{
				"code": "0410",
			})
		}
		c.JSONP(http.StatusOK, gin.H{
			"code": "0411",
		})
		return
	}

	var acc AccessToken
	if err := json.Unmarshal(body, &acc); err != nil {
		glog.Infof("bind is fail")
		c.JSONP(http.StatusOK, gin.H{
			"code": "0403",
		})
	}

	c.JSONP(http.StatusOK, gin.H{
		"code": "200",
	})
}
