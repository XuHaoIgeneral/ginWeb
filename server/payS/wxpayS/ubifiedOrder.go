package wxpayS

import (
	"errors"
	"ginweb/componet/pay/wxpayC"
	"ginweb/server/orderS"
	"github.com/golang/glog"
	"github.com/spf13/viper"
)

//统一下单
func UnifiedOrder(ip, openid, TradeType string, price int, wechat_client *wxpayC.WechatPay) (*wxpayC.UnifyOrderResult, error) {

	var pay_data wxpayC.UnitOrder
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
	localhostOrder := orderS.LocalOrder()
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
	return result, nil
}


//订单生成 （初始化订单）
func CreateOrder(wechat_client *wxpayC.WechatPay) *wxpayC.WechatPay {
	wechat_cert := []byte("1234567890qazxswa")
	wechat_key := []byte("1234567890wsxzaq")
	wechat_client = wxpayC.New(
		viper.GetString("wechat.xcx.appid"),
		viper.GetString("wechat.pay.mcid"),
		viper.GetString("wechat.pay.apikey"),
		wechat_key,
		wechat_cert, )
	return wechat_client
}