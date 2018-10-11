package orderS

import (
	"github.com/golang/glog"
	"strconv"
	"time"
)

//本地订单号生成，依赖时间关系
func LocalOrder() string {
	time := time.Now().UnixNano()
	timeUnixNano := strconv.FormatInt(time, 10)
	/*
			业务处理
	*/
	glog.Infof("本地订单号：%s", timeUnixNano)
	return string(timeUnixNano)
}