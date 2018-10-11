package wechatpay

import (
	"ginweb/componet/pay/wxpayC"
	"ginweb/server/payS/wxpayS"
	"ginweb/util/aesED"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/spf13/viper"
	"net/http"
	"strconv"
	"time"
)

var wechat_client *wxpayC.WechatPay

func XcxPay(c *gin.Context) {
	//由前端传递 openid
	token := c.DefaultPostForm("token", "null")
	if token == "null" {
		glog.Infof("获取openid失败，获取为null")
		c.JSONP(http.StatusOK, gin.H{
			"status": "0400",
		})
		return
	}
	openid, err := aesED.Decrypt(token)
	if err != nil {
		glog.Infof("token 解析 openid 失败")
		c.JSONP(http.StatusOK, gin.H{
			"status": "0400",
		})
		return
	}
	//获取类型
	TradeType := c.DefaultPostForm("TradeType", "JSAPI")
	if TradeType != "JSAPI" {
		glog.Infof("获取tradetype失败，获取为 null")
		c.JSONP(http.StatusOK, gin.H{
			"status": "0401",
		})
		return
	}
	//获取下单总价格
	price, err := strconv.Atoi(c.DefaultPostForm("price", "1"))
	if err != nil {
		glog.Infof("获取price失败，获取为 null")
		c.JSONP(http.StatusOK, gin.H{
			"status": "0402",
		})
		return
	}

	//订单生成
	wechat_client = wxpayS.CreateOrder(wechat_client)
	//获取ip
	ip := c.ClientIP()
	payResult := new(wxpayC.UnifyOrderResult)
	payResult, err = wxpayS.UnifiedOrder(ip, openid, "JSAPI", price, wechat_client)
	if err != nil {
		c.JSONP(http.StatusOK, gin.H{
			"status": "0403",
		})
		return
	}
	if err != nil {
		glog.Infof("%T==/n==%P==/n==%s", payResult, payResult, payResult)
		c.JSONP(http.StatusOK, gin.H{
			"status": "test",
		})
		return
	}

	if payResult.AppId != viper.GetString("wechat.xcx.appid") && payResult.MchId != viper.GetString("wechat.pay.mcid") {
		glog.Error("订单被篡改！")
		c.JSONP(http.StatusOK, gin.H{
			"status": "0404",
		})
		return
	}

	//下单回调处理
	res := make(map[string]interface{}, 0)
	res["appId"] = payResult.AppId
	res["nonceStr"] = payResult.NonceStr
	res["package"] = "prepay_id=" + payResult.PrepayId
	res["signType"] = "MD5"
	res["timeStamp"] = strconv.FormatInt(time.Now().Unix(), 10)
	resign := wxpayC.GetSign(res, viper.GetString("wechat.pay.apikey"))

	c.JSONP(http.StatusOK, gin.H{
		"appId":     payResult.AppId,
		"nonceStr":  payResult.NonceStr,
		"package":   "prepay_id=" + payResult.PrepayId,
		"signType":  "MD5",
		"timeStamp": res["timeStamp"],
		"sign":      resign,
	})
}
