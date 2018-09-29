package models

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"github.com/golang/glog"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

var x *xorm.Engine

func Init() {
	// 创建 ORM 引擎与数据库
	var err error
	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		viper.GetString("pgsql.host"),
		viper.GetString("pgsql.port"),
		viper.GetString("pgsql.user"),
		viper.GetString("pgsql.password"),
		viper.GetString("pgsql.dbName"))
	fmt.Println(psqlInfo)
	x, err = xorm.NewEngine(viper.GetString("pgsql.name"), psqlInfo)
	if err != nil {
		glog.Fatalf("Fail to create engine: %v\n", err)
	}

	// 同步结构体与数据表
	var dataStruct []interface{}

	dataStruct = append(dataStruct, new(Wx))
	dataStruct = append(dataStruct, new(User))
	dataStruct = append(dataStruct, new(TransactionDetail))
	dataStruct = append(dataStruct, new(Question))

	for i := range dataStruct {
		if err = x.Sync2(dataStruct[i]); err != nil {
			glog.Fatal("Fail to sync database wx %s", err)
		}
	}

}
