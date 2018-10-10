package wechatpay

import (
	"errors"
	"ginweb/server/aesED"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/spf13/viper"
	"net/http"
	"strconv"
	"time"
)

var wechat_client *WechatPay

func XcxPay(c *gin.Context) {
	wechat_client.Xcxpay(c)
}

func (this *WechatPay) Xcxpay(c *gin.Context) {
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
	glog.Infof(openid)
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
	this = CreateOrder(this)
	//获取ip
	ip := c.ClientIP()
	payResult := new(UnifyOrderResult)
	glog.Infof("this===%s", this.AppId)
	glog.Infof("openid===%s", openid)
	payResult, err = UnifiedOrder(ip, openid, "JSAPI", price, this)
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
	res["appid"] = payResult.AppId
	res["nonceStr"] = payResult.NonceStr
	res["prepayId"] = payResult.PrepayId
	res["signType"] = "MD5"
	res["timeStamp"] = strconv.FormatInt(time.Now().Unix(), 10)
	resign := GetSign(res, viper.GetString("wechat.wcx.secret"))

	c.JSONP(http.StatusOK, gin.H{
		"appid":     payResult.AppId,
		"nonceStr":  payResult.NonceStr,
		"package":   "prepay_id=" + payResult.PrepayId,
		"signType":  "MD5",
		"timeStamp": res["timeStamp"],
		"sign":      resign,
	})
}

//本地订单号生成，依赖时间关系
func localhostOrder() string {
	time := time.Now().UnixNano()
	timeUnixNano := strconv.FormatInt(time, 10)
	/*
			业务处理
	*/
	glog.Infof("本地订单号：%s", timeUnixNano)
	return string(timeUnixNano)
}

//订单生成
func CreateOrder(wechat_client *WechatPay) *WechatPay {
	wechat_cert := []byte("1234567890qazxswa")
	wechat_key := []byte("1234567890wsxzaq")
	wechat_client = New(
		viper.GetString("wechat.xcx.appid"),
		viper.GetString("wechat.pay.mcid"),
		viper.GetString("wechat.pay.apikey"),
		wechat_key,
		wechat_cert, )
	glog.Infof("%s", wechat_client)
	return wechat_client
}

//统一下单
func UnifiedOrder(ip, openid, TradeType string, price int, wechat_client *WechatPay) (*UnifyOrderResult, error) {

	var pay_data UnitOrder
	pay_data.NotifyUrl = viper.GetString("wechat.WECHAT_NOTIFY_URL")

	switch TradeType {
	case "NATIVE":
		pay_data.TradeType = "NATIVE"
	case "JSAPI":
		pay_data.TradeType = "JSAPI"
		glog.Infof("openid 为：%s", openid)
		pay_data.Openid = openid
	case "MWEB":
		pay_data.TradeType = "MWEB"
	}
	localhostOrder := localhostOrder()
	glog.Infof(localhostOrder)
	pay_data.Body = "测试-支付"
	pay_data.SpbillCreateIp = ip
	pay_data.TotalFee = price
	pay_data.OutTradeNo = localhostOrder //本地订单号

	//订单下达
	result, err := wechat_client.Pay(pay_data)

	if err != nil {
		glog.Infof("微信服务订单发送生成失败")
		return nil, errors.New("send fail about wechat order")
	}

	if result.ReturnCode == "FAIL" {
		glog.Infof("微信订单服务失败 %s", result.ReturnCode)
		return nil, errors.New("fail about wechat order server")
	}
	if result.ReturnMsg == "" {
		glog.Infof("签名失败 ，错误原因为 %s", result.ReturnMsg)
		return nil, errors.New("sign is fail!")
	}
	if result.ResultCode == "FAIL" {
		glog.Infof("业务结果失败")
		return nil, errors.New("wechat pay order active fail")
	}
	glog.Infof("result type: %T,%P", result, result)
	return result, nil
}
