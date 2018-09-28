package models

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"github.com/golang/glog"
	_ "github.com/lib/pq"
)

var x *xorm.Engine

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "123"
	dbName   = "pg_test"
)

func Init() {
	// 创建 ORM 引擎与数据库
	var err error
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)
	x, err = xorm.NewEngine("postgres", psqlInfo)
	if err != nil {
		glog.Fatalf("Fail to create engine: %v\n", err)
	}

	// 同步结构体与数据表  ****若本已经存在，会报错

	/*
		加入逻辑判断 判断下列的表是否存在
	 */
	if err = x.Sync2(new(Wx)); err != nil {
		glog.Fatal("Fail to sync database wx %s", err)
	}

	if err = x.Sync2(new(User)); err != nil {
		glog.Fatalf("Fail to sync database: %v\n", err)
	}

	if err = x.Sync2(new(TransactionDetail)); err != nil {
		glog.Fatalf("Fail to sync database: %s", err)
	}

	if err = x.Sync2(new(Question)); err != nil {
		glog.Fatalf("Fail to sync database %s", err)
	}
}
