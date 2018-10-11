package wxLoginS

import (
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
)

const (
	xcxUrl = "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=%s"
)

type XcxSessionKey struct {
	Openid      string `form:"openid" json:"openid" binding:"required"`
	Session_key string `form:"session_key" json:"session_key" binding:"required"`
	Unionid     string `form:"unionid" json:"unionid" binding:"required"`
}

//请求session_key openid unionid 并且绑定
func GetSessionKey(code string) (*XcxSessionKey, error) {
	var xcx XcxSessionKey
	url := fmt.Sprintf(xcxUrl, viper.GetString("wechat.xcx.appid"), viper.GetString("wechat.xcx.secret"), code, "authorization_code")

	resp, err := http.Get(url)

	defer resp.Body.Close()
	if err != nil {
		glog.Infof("请求session_key错误")
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Infof("请求session_key body错误")
		return nil, err
	}

	//判断返回参数是否带有unionid，获取openid session_key
	probe := gjson.Get(string(body), "unionid")
	if probe.String() == "" {
		glog.Infof("获取unionid为空 返回为:%s", string(body))
	}

	if err := json.Unmarshal(body, &xcx); err != nil {
		glog.Infof("bind is fail")
		return nil, err
	}
	return &xcx, nil
}
