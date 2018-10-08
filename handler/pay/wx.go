package pay

import (
	"ginweb/handler/wechat"
	"ginweb/pkg/wechatpay"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/spf13/viper"
	"net/http"
	"time"
)

var wechat_client *wechatpay.WechatPay

//接收前端传递的openid
func XcxPay(c *gin.Context) {

	openid := c.DefaultPostForm("openid", "null")
	if openid == "null" {
		glog.Infof("获取openid失败，获取为null")
		c.JSONP(http.StatusOK, gin.H{
			"status": "0400",
		})
		return
	}

	//订单生成
	CreateOrder()

	ip := c.ClientIP()
	openid := xcx.Openid
}

//本地订单号生成，依赖时间关系
func localhostOrder() string {
	timeUnixNano := string(time.Now().UnixNano())

	/*
			业务处理
	*/

	return timeUnixNano
}

//订单生成
func CreateOrder() {
	wechat_cert := "1234567890qazxsw"
	wechat_key := "1234567890wsxzaq"
	wechat_client = wechatpay.New(
		viper.GetString("wechat.xcx.appid"),
		viper.GetString("wechat.pay.mcid"),
		viper.GetString("wechat.wcx.secret"),
		[]byte(wechat_key),
		[]byte(wechat_cert))
}

func UnifiedOrder(ip, openid, TradeType string) error {
	var pay_data wechatpay.UnitOrder

	pay_data.NotifyUrl = viper.GetString("https://sd.wlinno.com")
	if TradeType == "NATIVE" { //二维码支付
		pay_data.TradeType = "NATIVE"
	} else if TradeType == "JSAPI" { //公众号支付，小程序支付
		/*
		获取openid

		 */
		pay_data.TradeType = "JSAPI"
		pay_data.Openid = openid
	} else if TradeType == "MWEB" { //H5支付
		pay_data.TradeType = "JSAPI"
		pay_data.Openid = openid
	}

	pay_data.Body = "测试-支付"
	pay_data.SpbillCreateIp = ip
	pay_data.TotalFee = 1
	pay_data.OutTradeNo = "123456789" //本地订单号
	result, err := wechat_client.Pay(pay_data)
	if err != nil {
		glog.Infof("微信服务订单发送生成失败")
		return err
	}

	if result.ReturnCode == "FAIL" {
		glog.Infof("微信订单服务失败")
		return err
	}

}

////创建支付
//func CreatePay()  {
//	wechat_cert, err := ioutil.ReadFile("conf/wechat/apiclient_cert.pem")
//	if err != nil {
//		panic(err)
//	}
//	wechat_key, err := ioutil.ReadFile("conf/wechat/apiclient_key.pem")
//	wechat_client = wechatpay.New(
//		viper.GetString("wechat.xcx.appid"),
//		viper.GetString("wechat.pay.mcid"),
//		viper.GetString("wechat.wcx.secret"),
//		wechat_key,
//		wechat_cert)
//
//	if err != nil {
//		panic(err)
//	}
//
//}
