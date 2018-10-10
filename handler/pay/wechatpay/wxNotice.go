package wechatpay

import (
	"encoding/xml"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

//微信扫码回调地址(gin框架)
//func PayNotifyUrl(c *gin.Context)  {
//	wechat_client.PaynotifyUrl(c)
//}

func PayNotifyUrl(c *gin.Context) {

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		glog.Error(err, "read notify body error")
	}
	var wx_notify_req PayNotifyResult
	err = xml.Unmarshal(body, &wx_notify_req)
	if err != nil {
		glog.Error(err, "read http body xml failed! err :"+err.Error())
	}
	var reqMap map[string]interface{}
	reqMap = make(map[string]interface{}, 0)
	
	reqMap["return_code"] = wx_notify_req.ReturnCode
	reqMap["return_msg"] = wx_notify_req.ReturnMsg
	reqMap["appid"] = wx_notify_req.AppId
	reqMap["mch_id"] = wx_notify_req.MchId
	reqMap["nonce_str"] = wx_notify_req.NonceStr
	reqMap["result_code"] = wx_notify_req.ResultCode
	reqMap["openid"] = wx_notify_req.OpenId
	reqMap["is_subscribe"] = wx_notify_req.IsSubscribe
	reqMap["trade_type"] = wx_notify_req.TradeType
	reqMap["bank_type"] = wx_notify_req.BankType
	reqMap["total_fee"] = wx_notify_req.TotalFee
	reqMap["fee_type"] = wx_notify_req.FeeType
	reqMap["cash_fee"] = wx_notify_req.CashFee
	reqMap["cash_fee_type"] = wx_notify_req.CashFeeType
	reqMap["transaction_id"] = wx_notify_req.TransactionId
	reqMap["out_trade_no"] = wx_notify_req.OutTradeNo
	reqMap["attach"] = wx_notify_req.Attach
	reqMap["time_end"] = wx_notify_req.TimeEnd

	wechat_client = CreateOrder(wechat_client)
	//进行签名校验
	if wechat_client.VerifySign(reqMap, wx_notify_req.Sign) {
		record, err := json.Marshal(wx_notify_req)
		if err != nil {
			glog.Error(err, "wechat pay marshal err :"+err.Error())
		}
		//TODO 加入你的代码，处理返回值
		fmt.Println(string(record))
		// err = wechat_pay_recoed_producer.Publish("wechat_pay", record)
		if err != nil {
			glog.Error(err, "wechat publish record err:"+err.Error())
		}
		c.XML(http.StatusOK, gin.H{
			"return_code": "SUCCESS",
			"return_msg":  "OK",
		})
	} else {
		c.XML(http.StatusOK, gin.H{
			"return_code": "FAIL",
			"return_msg":  "failed to verify sign, please retry!",
		})
	}
	return
}
